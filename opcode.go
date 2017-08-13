package main

// HAXE Intermediate Language Operation
type hxilOp int

const (
	hxilMov hxilOp = iota
	hxilInt
	hxilFloat
	hxilBool
	hxilBytes
	hxilString
	hxilNull

	hxilAdd
	hxilSub
	hxilMul
	hxilSDiv
	hxilUDiv
	hxilSMod
	hxilUMod
	hxilShl
	hxilSShr
	hxilUShr
	hxilAnd
	hxilhxilr
	hxilXor

	hxilNeg
	hxilNot
	hxilIncr
	hxilDecr

	hxilCall0
	hxilCall1
	hxilCall2
	hxilCall3
	hxilCall4
	hxilCallN
	hxilCallMethod
	hxilCallThis
	hxilCallClosure

	hxilStaticClosure
	hxilInstanceClosure
	hxilVirtualClosure

	hxilGetGlobal
	hxilSetGlobal
	hxilField
	hxilSetField
	hxilGetThis
	hxilSetThis
	hxilDynGet
	hxilDynSet

	hxilJTrue
	hxilJFalse
	hxilJNull
	hxilJNotNull
	hxilJSLt
	hxilJSGte
	hxilJSGt
	hxilJSLte
	hxilJULt
	hxilJUGte
	hxilJNotLt
	hxilJNotGte
	hxilJEq
	hxilJNotEq
	hxilJAlways

	hxilToDyn
	hxilToSFloat
	hxilToUFloat
	hxilToInt
	hxilSafeCast
	hxilUnsafeCast
	hxilToVirtual

	hxilLabel
	hxilRet
	hxilThrow
	hxilRethrow
	hxilSwitch
	hxilNullCheck
	hxilTrap
	hxilEndTrap

	hxilGetI8
	hxilGetI16
	hxilGetMem
	hxilGetArray
	hxilSetI8
	hxilSetI16
	hxilSetMem
	hxilSetArray

	hxilNew
	hxilArraySize
	hxilType
	hxilGetType
	hxilGetTID

	hxilRef
	hxilUnref
	hxilSetref

	hxilMakeEnum
	hxilEnumAlloc
	hxilEnumIndex
	hxilEnumField
	hxilSetEnumField

	hxilAssert
	hxilRefData
	hxilRefhxilffset
	hxilNop
)

type hxilOpData struct {
	name string
	args int
}

var (
	hxilOpCodes = []hxilOpData{
		hxilMov:    {"mov", 2},
		hxilInt:    {"int", 2},
		hxilFloat:  {"float", 2},
		hxilBool:   {"bool", 2},
		hxilBytes:  {"bytes", 2},
		hxilString: {"string", 2},
		hxilNull:   {"null", 1},

		hxilAdd:   {"add", 3},
		hxilSub:   {"sub", 3},
		hxilMul:   {"mul", 3},
		hxilSDiv:  {"sdiv", 3},
		hxilUDiv:  {"udiv", 3},
		hxilSMod:  {"smod", 3},
		hxilUMod:  {"umod", 3},
		hxilShl:   {"shl", 3},
		hxilSShr:  {"sshr", 3},
		hxilUShr:  {"ushr", 3},
		hxilAnd:   {"and", 3},
		hxilhxilr: {"or", 3},
		hxilXor:   {"xoe", 3},

		hxilNeg:  {"neg", 2},
		hxilNot:  {"not", 2},
		hxilIncr: {"incr", 1},
		hxilDecr: {"decr", 1},

		hxilCall0:       {"call", 2},
		hxilCall1:       {"call", 3},
		hxilCall2:       {"call", 4},
		hxilCall3:       {"call", 5},
		hxilCall4:       {"call", 6},
		hxilCallN:       {"call", -1},
		hxilCallMethod:  {"callmethod", -1},
		hxilCallThis:    {"callthis", -1},
		hxilCallClosure: {"callclosure", -1},

		hxilStaticClosure:   {"staticclosure", 2},
		hxilInstanceClosure: {"instanceclosure", 3},
		hxilVirtualClosure:  {"virtualclosure", 3},

		hxilGetGlobal: {"getglobal", 2},
		hxilSetGlobal: {"setglobal", 2},
		hxilField:     {"field", 3},
		hxilSetField:  {"setfield", 3},
		hxilGetThis:   {"getthis", 2},
		hxilSetThis:   {"setthis", 2},
		hxilDynGet:    {"dynget", 3},
		hxilDynSet:    {"dynset", 3},

		hxilJTrue:    {"jtrue", 2},
		hxilJFalse:   {"jfalse", 2},
		hxilJNull:    {"jnull", 2},
		hxilJNotNull: {"jnotnull", 2},
		hxilJSLt:     {"jslt", 3},
		hxilJSGte:    {"jsgte", 3},
		hxilJSGt:     {"jsgt", 3},
		hxilJSLte:    {"jslte", 3},
		hxilJULt:     {"jult", 3},
		hxilJUGte:    {"jugte", 3},
		hxilJNotLt:   {"jnotlt", 3},
		hxilJNotGte:  {"jnotgte", 3},
		hxilJEq:      {"jeq", 3},
		hxilJNotEq:   {"jnoteq", 3},
		hxilJAlways:  {"jalways", 1},

		hxilToDyn:      {"todyn", 2},
		hxilToSFloat:   {"tosfloat", 2},
		hxilToUFloat:   {"toufloat", 2},
		hxilToInt:      {"toint", 2},
		hxilSafeCast:   {"safecast", 2},
		hxilUnsafeCast: {"unsafecast", 2},
		hxilToVirtual:  {"tovirtual", 2},

		hxilLabel:     {"label", 0},
		hxilRet:       {"ret", 1},
		hxilThrow:     {"throw", 1},
		hxilRethrow:   {"rethrow", 1},
		hxilSwitch:    {"switch", -1},
		hxilNullCheck: {"nullcheck", 1},
		hxilTrap:      {"trap", 2},
		hxilEndTrap:   {"endtrap", 1},

		hxilGetI8:    {"geti8", 3},
		hxilGetI16:   {"geti16", 3},
		hxilGetMem:   {"getmem", 3},
		hxilGetArray: {"getarray", 3},
		hxilSetI8:    {"seti8", 3},
		hxilSetI16:   {"seti16", 3},
		hxilSetMem:   {"setmem", 3},
		hxilSetArray: {"setarray", 3},

		hxilNew:       {"new", 1},
		hxilArraySize: {"arraysize", 2},
		hxilType:      {"type", 2},
		hxilGetType:   {"gettype", 2},
		hxilGetTID:    {"gettid", 2},

		hxilRef:    {"ref", 2},
		hxilUnref:  {"unref", 2},
		hxilSetref: {"setref", 2},

		hxilMakeEnum:     {"makeenum", -1},
		hxilEnumAlloc:    {"enumalloc", 2},
		hxilEnumIndex:    {"enumindex", 2},
		hxilEnumField:    {"enumfield", 4},
		hxilSetEnumField: {"setenumfield", 3},

		hxilAssert:       {"assert", 0},
		hxilRefData:      {"refdata", 2},
		hxilRefhxilffset: {"refoffset", 3},
		hxilNop:          {"nop", 0},
	}
)
