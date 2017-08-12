package main

type hlTypeKind int

// HL Type
//go:generate stringer -type=hlTypeKind
const (
	HVOID hlTypeKind = iota
	HUI8
	HUI16
	HI32
	HI64
	HF32
	HF64
	HBOOL
	HBYTES
	HDYN
	HFUN
	HOBJ
	HARRAY
	HTYPE
	HREF
	HVIRTUAL
	HDYNOBJ
	HABSTRACT
	HENUM
	HNULL
	HLAST

	_H_FORCE_INT hlTypeKind = 0x7FFFFFFF
)

type hlType interface {
	//Kind() hlTypeKind
	//String() string
	/*
		hl_type_kind kind;
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

type hlTypeGeneric struct {
	kind hlTypeKind
}

type hlTypeVoid struct {
}

type hlTypeUI8 struct {
}

type hlTypeUI16 struct {
}

type hlTypeI32 struct {
}

type hlTypeI64 struct {
}

type hlTypeF32 struct {
}

type hlTypeF64 struct {
}

type hlTypeBool struct {
}

type hlTypeBytes struct {
}

type hlTypeDyn struct {
}

type hlTypeFun struct {
	args []int
	ret  int
	/*
		hl_type **args;
		hl_type *ret;
		int nargs;
		// storage for closure
		hl_type *parent;
		struct {
			hl_type_kind kind;
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

type hlTypeObj struct {
	super       int
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

type hlTypeType struct {
}

type hlTypeRef struct {
	tparam int
}

type hlTypeVirtual struct {
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

type hlTypeDynObj struct {
}

type hlTypeAbstract struct {
	name int
}

type hlTypeEnum struct {
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

type hlTypeNull struct {
}

type hlField struct {
	name   int
	hash   string
	hlType int
}
