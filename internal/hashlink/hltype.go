package hashlink

// Hashlink Data Type
//go:generate stringer -type=Id
type HdtId int

func (id HdtId) New() hlType {
	var t hlType

	switch id {
	case VoidType:
		t = new(hxtVoid)
	case UI8Type:
		t = new(hxtUI8)
	case UI16Type:
		t = new(hxtUI16)
	case I32Type:
		t = new(hxtI32)
	case I64Type:
		t = new(hxtI64)
	case F32Type:
		t = new(hxtF32)
	case F64Type:
		t = new(hxtF64)
	case BoolType:
		t = new(hxtBool)
	case BytesType:
		t = new(hxtBytes)
	case DynType:
		t = new(hxtDyn)
	case FunType:
		t = new(hxtFun)
	case ObjType:
		t = new(hxtObj)
	case ArrayType:
		t = new(hxtArray)
	case TypeType:
		t = new(hxtType)
	case RefType:
		t = new(hxtRef)
	case VirtualType:
		t = new(hxtVirtual)
	case DynObjType:
		t = new(hxtDynObj)
	case AbstractType:
		t = new(hxtAbstract)
	case EnumType:
		t = new(hxtEnum)
	case NullType:
		t = new(hxtNull)
	}

	return t
}

const (
	VoidType HdtId = iota
	UI8Type
	UI16Type
	I32Type
	I64Type
	F32Type
	F64Type
	BoolType
	BytesType
	DynType
	FunType
	ObjType
	ArrayType
	TypeType
	RefType
	VirtualType
	DynObjType
	AbstractType
	EnumType
	NullType
)

type hlType interface{}
type hlString []byte

type hxtVoid int
type hxtUI8 byte
type hxtUI16 uint16
type hxtI32 int32
type hxtI64 int64
type hxtF32 float32
type hxtF64 float64
type hxtBool bool
type hxtBytes []byte
type hxtDyn struct {
}

type hxtFun struct {
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

func (t *hxtFun) Unmarshal(ctx *Data, b *hlbStream) {
	nArg := int(b.byte())
	t.arg = make([]hlType, nArg)
	for i := 0; i < nArg; i++ {
		t.arg[i] = ctx.LookupType(b.index())
	}
	t.ret = ctx.LookupType(b.index())
}

type hxtObj struct {
	name     hlString
	super    int
	global   int
	lField   []hlField
	lProto   []hlProto
	lBinding []hlBinding
}

func (t *hxtObj) Unmarshal(ctx *Data, b *hlbStream) {
	t.name = ctx.LookupString(b.index())
	t.super = b.index()
	t.global = b.index()
	nField := b.index()
	nProto := b.index()
	nBinding := b.index()

	t.lField = make([]hlField, nField)
	for i := 0; i < nField; i++ {
		t.lField[i].name = ctx.LookupString(b.index())
		//t.field[i].hash = // TODO Hash name
		t.lField[i].typePtr = ctx.LookupType(b.index())
	}
	t.lProto = make([]hlProto, nProto)
	for i := 0; i < nProto; i++ {
		t.lProto[i].name = ctx.LookupString(b.index())
		t.lProto[i].fIndex = b.index()
		t.lProto[i].pIndex = b.index()
	}
	t.lBinding = make([]hlBinding, nBinding)
	for i := 0; i < nBinding; i++ {
		t.lBinding[i].index = b.index()
		t.lBinding[i].fIndex = b.index()
	}
}

type hxtArray struct {
}

type hxtType struct {
}

type hxtRef struct {
	tparam hlType
}

func (t *hxtRef) Unmarshal(ctx *Data, b *hlbStream) {
	t.tparam = ctx.LookupType(b.index())
}

type hxtVirtual struct {
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

func (t *hxtVirtual) Unmarshal(ctx *Data, b *hlbStream) {
	nField := b.index()
	t.field = make([]hlField, nField)
	for i := 0; i < nField; i++ {
		t.field[i].name = ctx.LookupString(b.index()) // read_ustring
		//t.field[i].hash // TODO Generate hash of name
		t.field[i].typePtr = ctx.LookupType(b.index())
	}
}

type hxtDynObj struct {
}

type hxtAbstract struct {
	name hlString
}

func (t *hxtAbstract) Unmarshal(ctx *Data, b *hlbStream) {
	t.name = ctx.LookupString(b.index())
}

type hxtEnum struct {
	name        hlString
	lConstruct  []hlEnumConstruct
	globalValue int
}

func (t *hxtEnum) Unmarshal(ctx *Data, b *hlbStream) {
	t.name = ctx.LookupString(b.index())
	t.globalValue = b.index()
	nConstruct := int(b.byte())
	t.lConstruct = make([]hlEnumConstruct, nConstruct)
	for i := 0; i < nConstruct; i++ {
		t.lConstruct[i].name = ctx.LookupString(b.index())
		nParam := b.index()
		t.lConstruct[i].arg = make([]hlType, nParam)
		for j := 0; j < nParam; j++ {
			t.lConstruct[i].arg[j] = ctx.LookupType(b.index())
		}
	}
}

type hxtNull struct {
	tparam hlType
}

func (t *hxtNull) Unmarshal(ctx *Data, b *hlbStream) {
	t.tparam = ctx.LookupType(b.index())
}

type hlField struct {
	name    hlString
	hash    uint32
	typePtr hlType
}

type hlProto struct {
	name   hlString
	hash   uint32
	fIndex int
	pIndex int
}

type hlBinding struct {
	index  int
	fIndex int
}

type hlEnumConstruct struct {
	name hlString
	arg  []hlType
}

type hlNative struct {
	lib    hlString
	name   hlString
	t      hlType
	findex int
}

type hlFunction struct {
	typePtr hlType
	fIndex  int
	lReg    []hlType
	lInst   []hxilInst
}

type Flags int

func (f Flags) HasDebug() bool { return f&1 == 1 }
