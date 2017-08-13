package main

type hxitId int

// Haxe Intermediate Type
//go:generate stringer -type=hxitId
const (
	hxitVoid hxitId = iota
	hxitUi8
	hxitUi16
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
	hxitDynobj
	hxitAbstract
	hxitEnum
	hxitNull
	hxitLast

	_H_FORCE_INT hxitId = 0x7FFFFFFF
)

type hxType interface {
	Kind() string
	/*
		hl_type_typeId typeId;
		union {
			const uchar *abs_name;
			hl_type_fun *fun;
			hl_type_obj *obj;
			hl_type_enum *tenum;
			hl_type_virtual *virt;
			hl_type	*tparam;
		};
		void **vobj_proto;
		unsigned int *mark_bits;
	*/
}

type hxTypeBase struct {
	typeId hxitId
}

func (t hxTypeBase) Kind() string {
	return t.typeId.String()
}

type hxTypeGeneric struct {
	hxTypeBase
}

type hxTypeFun struct {
	hxTypeBase
	args []int
	ret  int
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

type hxTypeObj struct {
	hxTypeBase
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

type hxTypeRef struct {
	hxTypeBase
	tparam int
}

type hxTypeVirtual struct {
	hxTypeBase
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

type hxTypeAbstract struct {
	hxTypeBase
	name string
}

type hxTypeEnum struct {
	hxTypeBase
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

type hxTypeNull struct {
	hxTypeBase
	tparam int
}

type hlField struct {
	name    string
	hash    string
	typeIdx int
}
