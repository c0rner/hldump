package hashlink

import (
	"fmt"
)

// HAXE Intermediate Language Operation Code
type HilOp int

const (
	OpMov HilOp = iota
	OpInt
	OpFloat
	OpBool
	OpBytes
	OpString
	OpNull

	OpAdd
	OpSub
	OpMul
	OpSDiv
	OpUDiv
	OpSMod
	OpUMod
	OpShl
	OpSShr
	OpUShr
	OpAnd
	Opr
	OpXor

	OpNeg
	OpNot
	OpIncr
	OpDecr

	OpCall0
	OpCall1
	OpCall2
	OpCall3
	OpCall4
	OpCallN
	OpCallMethod
	OpCallThis
	OpCallClosure

	OpStaticClosure
	OpInstanceClosure
	OpVirtualClosure

	OpGetGlobal
	OpSetGlobal
	OpField
	OpSetField
	OpGetThis
	OpSetThis
	OpDynGet
	OpDynSet

	OpJTrue
	OpJFalse
	OpJNull
	OpJNotNull
	OpJSLt
	OpJSGte
	OpJSGt
	OpJSLte
	OpJULt
	OpJUGte
	OpJNotLt
	OpJNotGte
	OpJEq
	OpJNotEq
	OpJAlways

	OpToDyn
	OpToSFloat
	OpToUFloat
	OpToInt
	OpSafeCast
	OpUnsafeCast
	OpToVirtual

	OpLabel
	OpRet
	OpThrow
	OpRethrow
	OpSwitch
	OpNullCheck
	OpTrap
	OpEndTrap

	OpGetI8
	OpGetI16
	OpGetMem
	OpGetArray
	OpSetI8
	OpSetI16
	OpSetMem
	OpSetArray

	OpNew
	OpArraySize
	OpType
	OpGetType
	OpGetTID

	OpRef
	OpUnref
	OpSetref

	OpMakeEnum
	OpEnumAlloc
	OpEnumIndex
	OpEnumField
	OpSetEnumField

	OpAssert
	OpRefData
	OpRefOpffset
	OpNop
)

type OpData struct {
	name string
	args int
}

var (
	OpCodes = []OpData{
		OpMov:    {"mov", 2},
		OpInt:    {"int", 2},
		OpFloat:  {"float", 2},
		OpBool:   {"bool", 2},
		OpBytes:  {"bytes", 2},
		OpString: {"string", 2},
		OpNull:   {"null", 1},

		OpAdd:  {"add", 3},
		OpSub:  {"sub", 3},
		OpMul:  {"mul", 3},
		OpSDiv: {"sdiv", 3},
		OpUDiv: {"udiv", 3},
		OpSMod: {"smod", 3},
		OpUMod: {"umod", 3},
		OpShl:  {"shl", 3},
		OpSShr: {"sshr", 3},
		OpUShr: {"ushr", 3},
		OpAnd:  {"and", 3},
		Opr:    {"or", 3},
		OpXor:  {"xoe", 3},

		OpNeg:  {"neg", 2},
		OpNot:  {"not", 2},
		OpIncr: {"incr", 1},
		OpDecr: {"decr", 1},

		OpCall0:       {"call", 2},
		OpCall1:       {"call", 3},
		OpCall2:       {"call", 4},
		OpCall3:       {"call", 5},
		OpCall4:       {"call", 6},
		OpCallN:       {"call", -1},
		OpCallMethod:  {"callmethod", -1},
		OpCallThis:    {"callthis", -1},
		OpCallClosure: {"callclosure", -1},

		OpStaticClosure:   {"staticclosure", 2},
		OpInstanceClosure: {"instanceclosure", 3},
		OpVirtualClosure:  {"virtualclosure", 3},

		OpGetGlobal: {"getglobal", 2},
		OpSetGlobal: {"setglobal", 2},
		OpField:     {"field", 3},
		OpSetField:  {"setfield", 3},
		OpGetThis:   {"getthis", 2},
		OpSetThis:   {"setthis", 2},
		OpDynGet:    {"dynget", 3},
		OpDynSet:    {"dynset", 3},

		OpJTrue:    {"jtrue", 2},
		OpJFalse:   {"jfalse", 2},
		OpJNull:    {"jnull", 2},
		OpJNotNull: {"jnotnull", 2},
		OpJSLt:     {"jslt", 3},
		OpJSGte:    {"jsgte", 3},
		OpJSGt:     {"jsgt", 3},
		OpJSLte:    {"jslte", 3},
		OpJULt:     {"jult", 3},
		OpJUGte:    {"jugte", 3},
		OpJNotLt:   {"jnotlt", 3},
		OpJNotGte:  {"jnotgte", 3},
		OpJEq:      {"jeq", 3},
		OpJNotEq:   {"jnoteq", 3},
		OpJAlways:  {"jalways", 1},

		OpToDyn:      {"todyn", 2},
		OpToSFloat:   {"tosfloat", 2},
		OpToUFloat:   {"toufloat", 2},
		OpToInt:      {"toint", 2},
		OpSafeCast:   {"safecast", 2},
		OpUnsafeCast: {"unsafecast", 2},
		OpToVirtual:  {"tovirtual", 2},

		OpLabel:     {"label", 0},
		OpRet:       {"ret", 1},
		OpThrow:     {"throw", 1},
		OpRethrow:   {"rethrow", 1},
		OpSwitch:    {"switch", -1},
		OpNullCheck: {"nullcheck", 1},
		OpTrap:      {"trap", 2},
		OpEndTrap:   {"endtrap", 1},

		OpGetI8:    {"geti8", 3},
		OpGetI16:   {"geti16", 3},
		OpGetMem:   {"getmem", 3},
		OpGetArray: {"getarray", 3},
		OpSetI8:    {"seti8", 3},
		OpSetI16:   {"seti16", 3},
		OpSetMem:   {"setmem", 3},
		OpSetArray: {"setarray", 3},

		OpNew:       {"new", 1},
		OpArraySize: {"arraysize", 2},
		OpType:      {"type", 2},
		OpGetType:   {"gettype", 2},
		OpGetTID:    {"gettid", 2},

		OpRef:    {"ref", 2},
		OpUnref:  {"unref", 2},
		OpSetref: {"setref", 2},

		OpMakeEnum:     {"makeenum", -1},
		OpEnumAlloc:    {"enumalloc", 2},
		OpEnumIndex:    {"enumindex", 2},
		OpEnumField:    {"enumfield", 4},
		OpSetEnumField: {"setenumfield", 3},

		OpAssert:     {"assert", 0},
		OpRefData:    {"refdata", 2},
		OpRefOpffset: {"refoffset", 3},
		OpNop:        {"nop", 0},
	}
)

// Haxe Intermediate Language Instruction
type HilInst struct {
	op    HilOp
	arg   []int
	extra []int
}

func (o *HilInst) Print(ctx *Data) {
	switch o.op {
	case OpInt:
		fmt.Printf("%s %d, %d ; ", OpCodes[o.op].name, o.arg[0], o.arg[1])
		fmt.Printf("%d, %d\n", ctx.ints[o.arg[0]], ctx.ints[o.arg[1]])
	case OpField:
		fmt.Printf("%s %d, %d[%d] ; \n", OpCodes[o.op].name, o.arg[0], o.arg[1], o.arg[2])
	case OpJTrue, OpJFalse, OpJNull, OpJNotNull:
		fmt.Printf("%s %d, %d\n", OpCodes[o.op].name, o.arg[0], o.arg[1])
	case OpJEq, OpJNotEq, OpSwitch:
		fmt.Printf("%s %d, %d, %d\n", OpCodes[o.op].name, o.arg[0], o.arg[1], o.arg[2])
	case OpJAlways:
		fmt.Printf("%s %d\n", OpCodes[o.op].name, o.arg[0])
	case OpCall0:
		tgt := ctx.LookupFunction(o.arg[1])
		fmt.Printf("%s ", OpCodes[o.op].name)
		fmt.Printf("%d,f@%d ; ", o.arg[0], o.arg[1])
		switch tgt := tgt.(type) {
		case *hlNative:
			fmt.Printf(".%s.%s", tgt.libPtr, tgt.namePtr)
		case *hlFunction:
			fmt.Printf(".%s()", tgt.field)
		}
		fmt.Println()
	case OpCall1:
		tgt := ctx.LookupFunction(o.arg[1])
		fmt.Printf("%s ", OpCodes[o.op].name)
		fmt.Printf("%d,f@%d(%d) ; ", o.arg[0], o.arg[1], o.arg[2])
		switch tgt := tgt.(type) {
		case *hlNative:
			fmt.Printf(".%s.%s", tgt.libPtr, tgt.namePtr)
		case *hlFunction:
			fmt.Printf(".%s()", tgt.field)
		}
		fmt.Println()
	case OpCall2:
		tgt := ctx.LookupFunction(o.arg[1])
		fmt.Printf("%s ", OpCodes[o.op].name)
		fmt.Printf("%d,f@%d(%d, %d) ; ", o.arg[0], o.arg[1], o.arg[2], o.arg[3])
		switch tgt := tgt.(type) {
		case *hlNative:
			fmt.Printf(".%s.%s", tgt.libPtr, tgt.namePtr)
		case *hlFunction:
			fmt.Printf(".%s()", tgt.field)
		}
		fmt.Println()
	case OpCall3:
		tgt := ctx.LookupFunction(o.arg[1])
		fmt.Printf("%s ", OpCodes[o.op].name)
		fmt.Printf("%d,f@%d(%d, %d, %d) ; ", o.arg[0], o.arg[1], o.arg[2], o.arg[3], o.arg[4])
		switch tgt := tgt.(type) {
		case *hlNative:
			fmt.Printf(".%s.%s", tgt.libPtr, tgt.namePtr)
		case *hlFunction:
			if tgt.obj != nil {
				ooo := tgt.obj.(*ObjType)
				fmt.Printf("%s.%s()", ooo.namePtr, tgt.field)
			}
		}
		fmt.Println()
	case OpString:
		fmt.Printf("%s ", OpCodes[o.op].name)
		fmt.Printf("%d,@%d ; \"%.40s\"", o.arg[0], o.arg[1], ctx.strings.String(o.arg[1]))
		fmt.Println()
		/*
			case OpSetField:
				fmt.Printf("%s ", OpCodes[o.op].name)
				tgt := ctx.getType(o.arg[0])
				switch tgt := tgt.(type) {
				case *hxtObj:
					fmt.Printf("%s.%d, %d ; (%d)", tgt.name, o.arg[1], o.arg[2], len(tgt.lField))
		*/
	default:
		fmt.Printf("%s", OpCodes[o.op].name)
		fmt.Println()
	}
}
