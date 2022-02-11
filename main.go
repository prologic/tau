//go:build !windows
// +build !windows

package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/NicoNex/tau/compiler"
	"github.com/NicoNex/tau/obj"
	"github.com/NicoNex/tau/parser"
	"github.com/NicoNex/tau/vm"
	"golang.org/x/term"
)

func encode(bcode *compiler.Bytecode) ([]byte, error) {
	var (
		buf bytes.Buffer
		enc = gob.NewEncoder(&buf)
	)

	for _, c := range bcode.Constants {
		gob.Register(c)
	}

	if err := enc.Encode(bcode); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func decode(r io.Reader) (*compiler.Bytecode, error) {
	var (
		b   *compiler.Bytecode
		dec = gob.NewDecoder(r)
	)

	gob.Register(obj.NewInteger(0))
	gob.Register(obj.NewBoolean(false))
	gob.Register(obj.NewNull())
	gob.Register(obj.NewClass())
	gob.Register(obj.NewReturn(nil))
	gob.Register(obj.NewFloat(0))
	gob.Register(obj.NewList())
	gob.Register(obj.NewMap())
	gob.Register(obj.NewString(""))
	gob.Register(obj.NewError(""))
	gob.Register(obj.NewClosure(nil, []obj.Object{}))
	gob.Register(obj.Builtin(func(arg ...obj.Object) obj.Object { return nil }))
	gob.Register(obj.NewFunction([]string{}, obj.NewEnv(), interface{}(nil)))

	return b, dec.Decode(&b)
}

func readFile(fname string) []byte {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return b
}

func writeFile(fname string, cont []byte) {
	if err := ioutil.WriteFile(fname, cont, 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func evalREPL() {
	var env = obj.NewEnv()

	initState, err := term.MakeRaw(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore(0, initState)

	t := term.NewTerminal(os.Stdin, ">>> ")
	obj.Stdout = t

	for {
		input, err := t.ReadLine()
		if err != nil {
			// Quit without error on Ctrl^D.
			if err != io.EOF {
				fmt.Println(err)
			}
			return
		}

		res, errs := parser.Parse(input)
		if len(errs) != 0 {
			for _, e := range errs {
				fmt.Fprintln(t, e)
			}
			continue
		}

		if val := res.Eval(env); val != nil && val != obj.NullObj {
			fmt.Fprintln(t, val)
		}
	}
}

func vmREPL() {
	var (
		consts      []obj.Object
		globals     []obj.Object
		symbolTable *compiler.SymbolTable
	)

	globals = make([]obj.Object, vm.GlobalSize)
	symbolTable = compiler.NewSymbolTable()

	for i, b := range obj.Builtins {
		symbolTable.DefineBuiltin(i, b.Name)
	}

	initState, err := term.MakeRaw(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore(0, initState)

	t := term.NewTerminal(os.Stdin, ">>> ")
	obj.Stdout = t

	for {
		input, err := t.ReadLine()
		if err != nil {
			// Quit without error on Ctrl^D.
			if err != io.EOF {
				fmt.Println(err)
			}
			return
		}

		res, errs := parser.Parse(input)
		if len(errs) != 0 {
			for _, e := range errs {
				fmt.Fprintln(t, e)
			}
			continue
		}

		c := compiler.NewWithState(symbolTable, consts)
		c.Compile(res)
		tvm := vm.NewWithGlobalStore(c.Bytecode(), globals)

		if err := tvm.Run(); err != nil {
			fmt.Fprintf(t, "runtime error: %v\n", err)
			continue
		}

		fmt.Fprintln(t, tvm.LastPoppedStackElem())
	}
}

func execFilesVM(files []string) {
	for _, f := range files {
		var bytecode *compiler.Bytecode

		if strings.HasSuffix(f, ".tauc") {
			file, err := os.Open(f)
			if err != nil {
				fmt.Println(err)
				continue
			}

			bytecode, err = decode(bufio.NewReader(file))
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			b := readFile(f)
			res, errs := parser.Parse(string(b))
			if len(errs) != 0 {
				for _, e := range errs {
					fmt.Println(e)
				}
				return
			}

			c := compiler.New()
			if err := c.Compile(res); err != nil {
				fmt.Println(err)
				return
			}
			bytecode = c.Bytecode()
		}

		tvm := vm.New(bytecode)
		if err := tvm.Run(); err != nil {
			fmt.Printf("runtime error: %v\n", err)
			return
		}

		fmt.Println(tvm.LastPoppedStackElem())
	}
}

func execFilesEval(files []string) {
	for _, f := range files {
		var env = obj.NewEnv()

		b := readFile(f)
		res, errs := parser.Parse(string(b))
		if len(errs) != 0 {
			for _, e := range errs {
				fmt.Println(e)
			}
			continue
		}

		val := res.Eval(env)
		if val != obj.NullObj && val != nil {
			fmt.Println(val)
		}
	}
}

func compileFiles(files []string) {
	for _, f := range files {
		b := readFile(f)

		res, errs := parser.Parse(string(b))
		if len(errs) != 0 {
			for _, e := range errs {
				fmt.Println(e)
			}
			return
		}

		c := compiler.New()
		if err := c.Compile(res); err != nil {
			fmt.Println(err)
			continue
		}

		cnt, err := encode(c.Bytecode())
		if err != nil {
			fmt.Println(err)
			continue
		}

		ext := filepath.Ext(f)
		writeFile(f[:len(f)-len(ext)]+".tauc", cnt)
	}
}

func main() {
	var (
		useVM   bool
		compile bool
	)

	flag.BoolVar(&useVM, "vm", false, "Use the Tau VM instead of eval method.")
	flag.BoolVar(&compile, "c", false, "Compile the tau file.")
	flag.Parse()

	if compile {
		compileFiles(flag.Args())
		return
	}

	if flag.NArg() > 0 {
		if useVM {
			execFilesVM(flag.Args())
		} else {
			execFilesEval(flag.Args())
		}
		return
	}

	if useVM {
		vmREPL()
	} else {
		evalREPL()
	}
}
