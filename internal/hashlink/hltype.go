package hashlink

// Hashlink Data Type
//go:generate stringer -type=Id
type HdtId int

const (
	VoidTypeId HdtId = iota
	UI8TypeId
	UI16TypeId
	I32TypeId
	I64TypeId
	F32TypeId
	F64TypeId
	BoolTypeId
	BytesTypeId
	DynTypeId
	FunTypeId
	ObjTypeId
	ArrayTypeId
	TypeTypeId
	RefTypeId
	VirtualTypeId
	DynObjTypeId
	AbstractTypeId
	EnumTypeId
	NullTypeId
)

func (id HdtId) NewType() hlType {
	var t hlType

	switch id {
	case VoidTypeId:
		t = new(VoidType)
	case UI8TypeId:
		t = new(UI8Type)
	case UI16TypeId:
		t = new(UI16Type)
	case I32TypeId:
		t = new(I32Type)
	case I64TypeId:
		t = new(I64Type)
	case F32TypeId:
		t = new(F32Type)
	case F64TypeId:
		t = new(F64Type)
	case BoolTypeId:
		t = new(BoolType)
	case BytesTypeId:
		t = new(BytesType)
	case DynTypeId:
		t = new(DynType)
	case FunTypeId:
		t = new(FunType)
	case ObjTypeId:
		t = new(ObjType)
	case ArrayTypeId:
		t = new(ArrayType)
	case TypeTypeId:
		t = new(TypeType)
	case RefTypeId:
		t = new(RefType)
	case VirtualTypeId:
		t = new(VirtualType)
	case DynObjTypeId:
		t = new(DynObjType)
	case AbstractTypeId:
		t = new(AbstractType)
	case EnumTypeId:
		t = new(EnumType)
	case NullTypeId:
		t = new(NullType)
	}

	return t
}

type hlType interface{}

type VoidType int
type UI8Type byte
type UI16Type uint16
type I32Type int32
type I64Type int64
type F32Type float32
type F64Type float64
type BoolType bool
type BytesType []byte
type DynType struct{}

type FunType struct {
	arg []hlType
	ret hlType
	/*
		hl_type **args;
		hl_type *ret;
		int nargs;
		// storage for closure
		hl_type *parent;
		struct {
			hl_type_typeId typeId;
			void *p;
		} closure_type;
		struct {
			hl_type **args;
			hl_type *ret;
			int nargs;
			hl_type *parent;
		} closure;
	*/
}

func (t *FunType) Unmarshal(ctx *Data, b *hlbStream) {
	nArg := int(b.byte())
	t.arg = make([]hlType, nArg)
	for i := 0; i < nArg; i++ {
		t.arg[i] = ctx.LookupType(b.index())
	}
	t.ret = ctx.LookupType(b.index())
}

type ObjType struct {
	nameIdx  int
	super    int
	global   int
	lField   []hlField
	lProto   []hlProto
	lBinding []hlBinding
}

func (t *ObjType) Unmarshal(ctx *Data, b *hlbStream) {
	t.nameIdx = b.index()
	t.super = b.index()
	t.global = b.index()
	nField := b.index()
	nProto := b.index()
	nBinding := b.index()

	t.lField = make([]hlField, nField)
	for i := 0; i < nField; i++ {
		t.lField[i].nameIdx = b.index()
		//t.field[i].hash = // TODO Hash name
		t.lField[i].typePtr = ctx.LookupType(b.index())
	}
	t.lProto = make([]hlProto, nProto)
	for i := 0; i < nProto; i++ {
		t.lProto[i].nameIdx = b.index()
		t.lProto[i].funIdx = b.index()
		t.lProto[i].protoIdx = b.index()
	}
	t.lBinding = make([]hlBinding, nBinding)
	for i := 0; i < nBinding; i++ {
		t.lBinding[i].index = b.index()
		t.lBinding[i].funIdx = b.index()
	}
}

type ArrayType struct {
}

type TypeType struct {
}

type RefType struct {
	tparam hlType
}

func (t *RefType) Unmarshal(ctx *Data, b *hlbStream) {
	t.tparam = ctx.LookupType(b.index())
}

type VirtualType struct {
	field []hlField
	/*
		hl_obj_field *fields;
		int nfields;
		// runtime
		int dataSize;
		int *indexes;
		hl_field_lookup *lookup;
	*/
}

func (t *VirtualType) Unmarshal(ctx *Data, b *hlbStream) {
	nField := b.index()
	t.field = make([]hlField, nField)
	for i := 0; i < nField; i++ {
		t.field[i].nameIdx = b.index()
		//t.field[i].hash // TODO Generate hash of name
		t.field[i].typePtr = ctx.LookupType(b.index())
	}
}

type DynObjType struct {
}

type AbstractType struct {
	nameIdx int
}

func (t *AbstractType) Unmarshal(ctx *Data, b *hlbStream) {
	t.nameIdx = b.index()
}

type EnumType struct {
	nameIdx     int
	lConstruct  []hlEnumConstruct
	globalValue int
}

func (t *EnumType) Unmarshal(ctx *Data, b *hlbStream) {
	t.nameIdx = b.index()
	t.globalValue = b.index()
	nConstruct := int(b.byte())
	t.lConstruct = make([]hlEnumConstruct, nConstruct)
	for i := 0; i < nConstruct; i++ {
		t.lConstruct[i].nameIdx = b.index()
		nParam := b.index()
		t.lConstruct[i].arg = make([]hlType, nParam)
		for j := 0; j < nParam; j++ {
			t.lConstruct[i].arg[j] = ctx.LookupType(b.index())
		}
	}
}

type NullType struct {
	tparam hlType
}

func (t *NullType) Unmarshal(ctx *Data, b *hlbStream) {
	t.tparam = ctx.LookupType(b.index())
}

type hlField struct {
	nameIdx int
	hash    uint32
	typePtr hlType
}

type hlProto struct {
	nameIdx  int
	hash     uint32
	funIdx   int
	protoIdx int
}

type hlBinding struct {
	index  int
	funIdx int
}

type hlEnumConstruct struct {
	nameIdx int
	arg     []hlType
}

type hlNative struct {
	libIdx  int
	nameIdx int
	t       hlType
	funIdx  int
}

type hlFunction struct {
	typePtr hlType
	funIdx  int
	lReg    []hlType
	lInst   []hxilInst
}

type Flags int

func (f Flags) HasDebug() bool { return f&1 == 1 }
