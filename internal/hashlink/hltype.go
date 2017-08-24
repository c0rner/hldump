package hashlink

func (id HdtId) NewType() hlType {
	var t hlType

	switch id {
	case VoidT:
		t = new(VoidType)
	case UI8T:
		t = new(UI8Type)
	case UI16T:
		t = new(UI16Type)
	case I32T:
		t = new(I32Type)
	case I64T:
		t = new(I64Type)
	case F32T:
		t = new(F32Type)
	case F64T:
		t = new(F64Type)
	case BoolT:
		t = new(BoolType)
	case BytesT:
		t = new(BytesType)
	case DynT:
		t = new(DynType)
	case FunT:
		t = new(FunType)
	case ObjT:
		t = new(ObjType)
	case ArrayT:
		t = new(ArrayType)
	case TypeT:
		t = new(TypeType)
	case RefT:
		t = new(RefType)
	case VirtualT:
		t = new(VirtualType)
	case DynObjT:
		t = new(DynObjType)
	case AbstractT:
		t = new(AbstractType)
	case EnumT:
		t = new(EnumType)
	case NullT:
		t = new(NullType)
	}

	return t
}

type hlType interface {
	Id() HdtId
}

type VoidType int

func (t *VoidType) Id() HdtId {
	return VoidT
}

type UI8Type byte

func (t *UI8Type) Id() HdtId {
	return UI8T
}

type UI16Type uint16

func (t *UI16Type) Id() HdtId {
	return UI16T
}

type I32Type int32

func (t *I32Type) Id() HdtId {
	return I32T
}

type I64Type int64

func (t *I64Type) Id() HdtId {
	return I64T
}

type F32Type float32

func (t *F32Type) Id() HdtId {
	return F32T
}

type F64Type float64

func (t *F64Type) Id() HdtId {
	return F64T
}

type BoolType bool

func (t *BoolType) Id() HdtId {
	return BoolT
}

type BytesType []byte

func (t *BytesType) Id() HdtId {
	return BytesT
}

type DynType struct{}

func (t *DynType) Id() HdtId {
	return DynT
}

type FunType struct {
	argIdx []int
	retIdx int
	argPtr []hlType
	retPtr hlType
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

func (t *FunType) Id() HdtId {
	return FunT
}

func (t *FunType) Unmarshal(ctx *Data, b *hlbStream) {
	nArg := int(b.byte())
	t.argIdx = make([]int, nArg)
	for i := 0; i < nArg; i++ {
		t.argIdx[i] = b.index()
	}
	t.retIdx = b.index()
}

type ObjType struct {
	nameIdx  int
	namePtr  []byte
	superIdx int
	superPtr *ObjType
	global   int
	offset   int
	lField   []hlField
	lProto   []hlProto
	lBinding []hlBinding
}

func (t *ObjType) Id() HdtId {
	return ObjT
}

func (t *ObjType) Unmarshal(ctx *Data, b *hlbStream) {
	t.nameIdx = b.index()
	t.superIdx = b.index()
	t.global = b.index()
	nField := b.index()
	nProto := b.index()
	nBinding := b.index()

	if t.superIdx > 0 {
		super := ctx.LookupType(t.superIdx).(*ObjType)
		t.offset = super.offset + len(super.lField)
		t.offset += len(t.lField)
	}

	t.lField = make([]hlField, nField)
	for i := 0; i < nField; i++ {
		t.lField[i].nameIdx = b.index()
		//t.field[i].hash = // TODO Hash name
		t.lField[i].typeIdx = b.index()
	}
	t.lProto = make([]hlProto, nProto)
	for i := 0; i < nProto; i++ {
		t.lProto[i].nameIdx = b.index()
		t.lProto[i].funcIdx = b.index()
		t.lProto[i].override = b.index()
	}
	t.lBinding = make([]hlBinding, nBinding)
	for i := 0; i < nBinding; i++ {
		t.lBinding[i].fldIdx = b.index()
		t.lBinding[i].funcIdx = b.index()
	}
}

type ArrayType struct {
}

func (t *ArrayType) Id() HdtId {
	return ArrayT
}

type TypeType struct {
}

func (t *TypeType) Id() HdtId {
	return TypeT
}

type RefType struct {
	paramIdx int
}

func (t *RefType) Id() HdtId {
	return RefT
}

func (t *RefType) Unmarshal(ctx *Data, b *hlbStream) {
	t.paramIdx = b.index()
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

func (t *VirtualType) Id() HdtId {
	return VirtualT
}

func (t *VirtualType) Unmarshal(ctx *Data, b *hlbStream) {
	nField := b.index()
	t.field = make([]hlField, nField)
	for i := 0; i < nField; i++ {
		t.field[i].nameIdx = b.index()
		//t.field[i].hash // TODO Generate hash of name
		t.field[i].typeIdx = b.index()
	}
}

type DynObjType struct {
}

func (t *DynObjType) Id() HdtId {
	return DynObjT
}

type AbstractType struct {
	nameIdx int
	namePtr []byte
}

func (t *AbstractType) Id() HdtId {
	return AbstractT
}

func (t *AbstractType) Unmarshal(ctx *Data, b *hlbStream) {
	t.nameIdx = b.index()
}

type EnumType struct {
	nameIdx     int
	namePtr     []byte
	lConstruct  []hlEnumConstruct
	globalValue int
}

func (t *EnumType) Id() HdtId {
	return EnumT
}

func (t *EnumType) Unmarshal(ctx *Data, b *hlbStream) {
	t.nameIdx = b.index()
	t.globalValue = b.index()
	nConstruct := int(b.byte())
	t.lConstruct = make([]hlEnumConstruct, nConstruct)
	for i := 0; i < nConstruct; i++ {
		t.lConstruct[i].nameIdx = b.index()
		nParam := b.index()
		t.lConstruct[i].argIdx = make([]int, nParam)
		for j := 0; j < nParam; j++ {
			t.lConstruct[i].argIdx[j] = b.index()
		}
	}
}

type NullType struct {
	paramIdx int
}

func (t *NullType) Id() HdtId {
	return NullT
}

func (t *NullType) Unmarshal(ctx *Data, b *hlbStream) {
	t.paramIdx = b.index()
}

type hlField struct {
	nameIdx int
	namePtr int
	hash    uint32
	typeIdx int
	typePtr int
}

type hlProto struct {
	nameIdx  int
	namePtr  int
	hash     uint32
	funcIdx  int
	funcPtr  int
	override int
}

type hlBinding struct {
	fldIdx  int
	fldPtr  int
	funcIdx int
	funcPtr int
}

type hlEnumConstruct struct {
	nameIdx int
	namePtr int
	argIdx  []int
}

type Function interface{}

type hlNative struct {
	libIdx  int
	libPtr  string
	nameIdx int
	namePtr string
	typeIdx int
	typePtr int
	funcIdx int
	funcPtr int
}

type hlFunction struct {
	typeIdx int
	typePtr int
	funcIdx int
	funcPtr int
	regIdx  []int
	inst    []HilInst
	obj     hlType
	field   []byte
}

type Flags int

func (f Flags) HasDebug() bool { return f&1 == 1 }
