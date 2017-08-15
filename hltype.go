package main

// Haxe Intermediate t
//go:generate stringer -type=hxitId
type hxitId int

const (
	hxitVoid hxitId = iota
	hxitUI8
	hxitUI16
	hxitI32
	hxitI64
	hxitF32
	hxitF64
	hxitBool
	hxitBytes
	hxitDyn
	hxitFun
	hxitObj
	hxitArray
	hxitType
	hxitRef
	hxitVirtual
	hxitDynObj
	hxitAbstract
	hxitEnum
	hxitNull
)

type hxType interface{}

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
	arg []hxType
	ret hxType
	fun *hlFunction
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

func (t *hxtFun) Unmarshal(c hldumper, b *hlBuf) {
	nArg := int(b.byte())
	t.arg = make([]hxType, nArg)
	for i := 0; i < nArg; i++ {
		t.arg[i] = c.GetType(b.index())
	}
	t.ret = c.GetType(b.index())
}

type hxtObj struct {
	name     string
	super    *hxtObj
	global   hxType
	lField   []hlField
	lProto   []hlProto
	lBinding []hlBinding
	/*
		int nfields;
		int nproto;
		int nbindings;
		const uchar *name;
		hl_type *super;
		hl_obj_field *fields;
		hl_obj_proto *proto;
		int *bindings;
		void **global_value;
		hl_module_context *m;
		hl_runtime_obj *rt;
	*/
}

func (t *hxtObj) Unmarshal(c hldumper, b *hlBuf) {
	t.name = c.GetString(b.index())
	super := b.index()
	t.global = b.index()
	nField := b.index()
	nProto := b.index()
	nBinding := b.index()

	if super > 0 {
		t.super = c.GetType(super).(*hxtObj)
	}

	t.lField = make([]hlField, nField)
	for i := 0; i < nField; i++ {
		t.lField[i].name = c.GetString(b.index())
		//t.field[i].hash = // TODO Hash name
		t.lField[i].typePtr = c.GetType(b.index())
	}
	t.lProto = make([]hlProto, nProto)
	for i := 0; i < nProto; i++ {
		t.lProto[i].name = c.GetString(b.index())
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
	tparam hxType
}

func (t *hxtRef) Unmarshal(c hldumper, b *hlBuf) {
	t.tparam = c.GetType(b.index())
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

func (t *hxtVirtual) Unmarshal(c hldumper, b *hlBuf) {
	nField := b.index()
	t.field = make([]hlField, nField)
	for i := 0; i < nField; i++ {
		t.field[i].name = c.GetString(b.index()) // read_ustring
		//t.field[i].hash // TODO Generate hash of name
		t.field[i].typePtr = c.GetType(b.index())
	}
}

type hxtDynObj struct {
}

type hxtAbstract struct {
	name string
}

func (t *hxtAbstract) Unmarshal(c hldumper, b *hlBuf) {
	t.name = c.GetString(b.index())
}

type hxtEnum struct {
	name        string
	lConstruct  []hlEnumConstruct
	globalValue int
}

func (t *hxtEnum) Unmarshal(c hldumper, b *hlBuf) {
	t.name = c.GetString(b.index())
	t.globalValue = b.index()
	nConstruct := int(b.byte())
	t.lConstruct = make([]hlEnumConstruct, nConstruct)
	for i := 0; i < nConstruct; i++ {
		t.lConstruct[i].name = c.GetString(b.index())
		nParam := b.index()
		t.lConstruct[i].arg = make([]hxType, nParam)
		for j := 0; j < nParam; j++ {
			t.lConstruct[i].arg[j] = c.GetType(b.index())
		}
	}
}

type hxtNull struct {
	tparam hxType
}

func (t *hxtNull) Unmarshal(c hldumper, b *hlBuf) {
	t.tparam = c.GetType(b.index())
}

func NewHXType(id hxitId) hxType {
	var t hxType

	switch id {
	case hxitVoid:
		t = new(hxtVoid)
	case hxitUI8:
		t = new(hxtUI8)
	case hxitUI16:
		t = new(hxtUI16)
	case hxitI32:
		t = new(hxtI32)
	case hxitI64:
		t = new(hxtI64)
	case hxitF32:
		t = new(hxtF32)
	case hxitF64:
		t = new(hxtF64)
	case hxitBool:
		t = new(hxtBool)
	case hxitBytes:
		t = new(hxtBytes)
	case hxitDyn:
		t = new(hxtDyn)
	case hxitFun:
		t = new(hxtFun)
	case hxitObj:
		t = new(hxtObj)
	case hxitArray:
		t = new(hxtArray)
	case hxitType:
		t = new(hxtType)
	case hxitRef:
		t = new(hxtRef)
	case hxitVirtual:
		t = new(hxtVirtual)
	case hxitDynObj:
		t = new(hxtDynObj)
	case hxitAbstract:
		t = new(hxtAbstract)
	case hxitEnum:
		t = new(hxtEnum)
	case hxitNull:
		t = new(hxtNull)
	}

	return t
}

type hlField struct {
	name    string
	hash    uint32
	typePtr hxType
}

type hlProto struct {
	name   string
	hash   uint32
	fIndex int
	pIndex int
}

type hlBinding struct {
	index  int
	fIndex int
}

type hlEnumConstruct struct {
	name string
	arg  []hxType
}
