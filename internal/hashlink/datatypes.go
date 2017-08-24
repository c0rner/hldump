package hashlink

// Hashlink Data Type
//go:generate stringer -type=HdtId
type HdtId int

const (
	VoidT HdtId = iota
	UI8T
	UI16T
	I32T
	I64T
	F32T
	F64T
	BoolT
	BytesT
	DynT
	FunT
	ObjT
	ArrayT
	TypeT
	RefT
	VirtualT
	DynObjT
	AbstractT
	EnumT
	NullT
)
