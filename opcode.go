package main

type hlOp int

const (
	OMov hlOp = iota
	OInt
	OFloat
	OBool
	OBytes
	OString
	ONull

	OAdd
	OSub
	OMul
	OSDiv
	OUDiv
	OSMod
	OUMod
	OShl
	OSShr
	OUShr
	OAnd
	OOr
	OXor

	ONeg
	ONot
	OIncr
	ODecr

	OCall0
	OCall1
	OCall2
	OCall3
	OCall4
	OCallN
	OCallMethod
	OCallThis
	OCallClosure

	OStaticClosure
	OInstanceClosure
	OVirtualClosure

	OGetGlobal
	OSetGlobal
	OField
	OSetField
	OGetThis
	OSetThis
	ODynGet
	ODynSet

	OJTrue
	OJFalse
	OJNull
	OJNotNull
	OJSLt
	OJSGte
	OJSGt
	OJSLte
	OJULt
	OJUGte
	OJNotLt
	OJNotGte
	OJEq
	OJNotEq
	OJAlways

	OToDyn
	OToSFloat
	OToUFloat
	OToInt
	OSafeCast
	OUnsafeCast
	OToVirtual

	OLabel
	ORet
	OThrow
	ORethrow
	OSwitch
	ONullCheck
	OTrap
	OEndTrap

	OGetI8
	OGetI16
	OGetMem
	OGetArray
	OSetI8
	OSetI16
	OSetMem
	OSetArray

	ONew
	OArraySize
	OType
	OGetType
	OGetTID

	ORef
	OUnref
	OSetref

	OMakeEnum
	OEnumAlloc
	OEnumIndex
	OEnumField
	OSetEnumField

	OAssert
	ORefData
	ORefOffset
	ONop

	OLast
)

type hlOpData struct {
	name string
	args int
}

var (
	hlOpCodes = []hlOpData{
		OMov:    {"mov", 2},
		OInt:    {"int", 2},
		OFloat:  {"float", 2},
		OBool:   {"bool", 2},
		OBytes:  {"bytes", 2},
		OString: {"string", 2},
		ONull:   {"null", 1},

		OAdd:  {"add", 3},
		OSub:  {"sub", 3},
		OMul:  {"mul", 3},
		OSDiv: {"sdiv", 3},
		OUDiv: {"udiv", 3},
		OSMod: {"smod", 3},
		OUMod: {"umod", 3},
		OShl:  {"shl", 3},
		OSShr: {"sshr", 3},
		OUShr: {"ushr", 3},
		OAnd:  {"and", 3},
		OOr:   {"or", 3},
		OXor:  {"xoe", 3},

		ONeg:  {"neg", 2},
		ONot:  {"not", 2},
		OIncr: {"incr", 1},
		ODecr: {"decr", 1},

		OCall0:       {"call", 2},
		OCall1:       {"call", 3},
		OCall2:       {"call", 4},
		OCall3:       {"call", 5},
		OCall4:       {"call", 6},
		OCallN:       {"call", -1},
		OCallMethod:  {"callmethod", -1},
		OCallThis:    {"callthis", -1},
		OCallClosure: {"callclosure", -1},

		OStaticClosure:   {"OStaticClosure", 2},
		OInstanceClosure: {"OInstanceClosure", 3},
		OVirtualClosure:  {"OVirtualClosure", 3},

		OGetGlobal: {"OGetGlobal", 2},
		OSetGlobal: {"OSetGlobal", 2},
		OField:     {"OField", 3},
		OSetField:  {"OSetField", 3},
		OGetThis:   {"OGetThis", 2},
		OSetThis:   {"OSetThis", 2},
		ODynGet:    {"ODynGet", 3},
		ODynSet:    {"ODynSet", 3},

		{"OJTrue", 2},
		{"OJFalse", 2},
		{"OJNull", 2},
		{"OJNotNull", 2},
		{"OJSLt", 3},
		{"OJSGte", 3},
		{"OJSGt", 3},
		{"OJSLte", 3},
		{"OJULt", 3},
		{"OJUGte", 3},
		{"OJNotLt", 3},
		{"OJNotGte", 3},
		{"OJEq", 3},
		{"OJNotEq", 3},
		{"OJAlways", 1},

		{"OToDyn", 2},
		{"OToSFloat", 2},
		{"OToUFloat", 2},
		{"OToInt", 2},
		{"OSafeCast", 2},
		{"OUnsafeCast", 2},
		{"OToVirtual", 2},

		{"OLabel", 0},
		{"ORet", 1},
		{"OThrow", 1},
		{"ORethrow", 1},
		{"OSwitch", -1},
		{"ONullCheck", 1},
		{"OTrap", 2},
		{"OEndTrap", 1},

		{"OGetI8", 3},
		{"OGetI16", 3},
		{"OGetMem", 3},
		{"OGetArray", 3},
		{"OSetI8", 3},
		{"OSetI16", 3},
		{"OSetMem", 3},
		{"OSetArray", 3},

		{"ONew", 1},
		{"OArraySize", 2},
		{"OType", 2},
		{"OGetType", 2},
		{"OGetTID", 2},

		{"ORef", 2},
		{"OUnref", 2},
		{"OSetref", 2},

		{"OMakeEnum", -1},
		{"OEnumAlloc", 2},
		{"OEnumIndex", 2},
		{"OEnumField", 4},
		{"OSetEnumField", 3},

		{"OAssert", 0},
		{"ORefData", 2},
		{"ORefOffset", 3},
		{"ONop", 0},

		{"OLast", 0},
	}
)
