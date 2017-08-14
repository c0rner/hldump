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
	args []hxType
	ret  hxType
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
	nArgs := int(b.byte())
	t.args = make([]hxType, nArgs)
	for i := 0; i < nArgs; i++ {
		t.args[i] = c.GetType(b.index())
	}
	t.ret = c.GetType(b.index())
}

type hxtObj struct {
	name        string
	super       int
	field       []hlField
	globalValue int
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
	t.super = b.index() // FIXME - NULL super
	t.globalValue = b.index()
	nField := b.index()
	nProto := b.index()
	nBinding := b.index()

	t.field = make([]hlField, nField)
	for i := 0; i < nField; i++ {
		t.field[i].name = c.GetString(b.index())
		//t.field[i].hash = // TODO Hash name
		t.field[i].typePtr = c.GetType(b.index())
	}
	for i := 0; i < nProto; i++ {
		b.index() // read_ustring
		b.index()
		b.index()
	}
	for i := 0; i < nBinding; i++ {
		b.index()
		b.index()
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
	nConstructs int
	globalValue int
	/*
		const uchar *name;
		int nconstructs;
		hl_enum_construct *constructs;
		void **global_value;
	*/
}

func (t *hxtEnum) Unmarshal(c hldumper, b *hlBuf) {
	t.name = c.GetString(b.index())
	t.globalValue = b.index()
	t.nConstructs = int(b.byte())
	for i := 0; i < t.nConstructs; i++ { // TODO Fix constructs
		c.GetString(b.index())
		nParam := b.index()
		for j := 0; j < nParam; j++ {
			c.GetType(b.index())
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
	hash    string
	typePtr hxType
}
