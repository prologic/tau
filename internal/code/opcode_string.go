// Code generated by "stringer -type=Opcode"; DO NOT EDIT.

package code

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpConstant-0]
	_ = x[OpTrue-1]
	_ = x[OpFalse-2]
	_ = x[OpNull-3]
	_ = x[OpList-4]
	_ = x[OpMap-5]
	_ = x[OpClosure-6]
	_ = x[OpCurrentClosure-7]
	_ = x[OpAdd-8]
	_ = x[OpSub-9]
	_ = x[OpMul-10]
	_ = x[OpDiv-11]
	_ = x[OpMod-12]
	_ = x[OpBwAnd-13]
	_ = x[OpBwOr-14]
	_ = x[OpBwXor-15]
	_ = x[OpBwNot-16]
	_ = x[OpBwLShift-17]
	_ = x[OpBwRShift-18]
	_ = x[OpAnd-19]
	_ = x[OpOr-20]
	_ = x[OpEqual-21]
	_ = x[OpNotEqual-22]
	_ = x[OpGreaterThan-23]
	_ = x[OpGreaterThanEqual-24]
	_ = x[OpIn-25]
	_ = x[OpMinus-26]
	_ = x[OpBang-27]
	_ = x[OpIndex-28]
	_ = x[OpCall-29]
	_ = x[OpConcurrentCall-30]
	_ = x[OpReturn-31]
	_ = x[OpReturnValue-32]
	_ = x[OpJump-33]
	_ = x[OpJumpNotTruthy-34]
	_ = x[OpDot-35]
	_ = x[OpDefine-36]
	_ = x[OpGetGlobal-37]
	_ = x[OpSetGlobal-38]
	_ = x[OpGetLocal-39]
	_ = x[OpSetLocal-40]
	_ = x[OpGetBuiltin-41]
	_ = x[OpGetFree-42]
	_ = x[OpLoadModule-43]
	_ = x[OpInterpolate-44]
	_ = x[OpPop-45]
}

const _Opcode_name = "OpConstantOpTrueOpFalseOpNullOpListOpMapOpClosureOpCurrentClosureOpAddOpSubOpMulOpDivOpModOpBwAndOpBwOrOpBwXorOpBwNotOpBwLShiftOpBwRShiftOpAndOpOrOpEqualOpNotEqualOpGreaterThanOpGreaterThanEqualOpInOpMinusOpBangOpIndexOpCallOpConcurrentCallOpReturnOpReturnValueOpJumpOpJumpNotTruthyOpDotOpDefineOpGetGlobalOpSetGlobalOpGetLocalOpSetLocalOpGetBuiltinOpGetFreeOpLoadModuleOpInterpolateOpPop"

var _Opcode_index = [...]uint16{0, 10, 16, 23, 29, 35, 40, 49, 65, 70, 75, 80, 85, 90, 97, 103, 110, 117, 127, 137, 142, 146, 153, 163, 176, 194, 198, 205, 211, 218, 224, 240, 248, 261, 267, 282, 287, 295, 306, 317, 327, 337, 349, 358, 370, 383, 388}

func (i Opcode) String() string {
	if i >= Opcode(len(_Opcode_index)-1) {
		return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}
