package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type hlBuf []byte

func (b *hlBuf) skip(i int) {
	if i > len(*b) {
		i = len(*b)
	}
	*b = (*b)[i:]
}

func (b *hlBuf) byte() byte {
	res := (*b)[0]
	*b = (*b)[1:]
	return res
}

func (b *hlBuf) int32() int32 {
	res := int32(binary.LittleEndian.Uint32(*b))
	*b = (*b)[4:]
	return res
}

// This is probably incorrect but it is inconclusive
// how doubles are handled by Hashlib
// For reference see hl_read_double()
// https://github.com/HaxeFoundation/hashlink/blob/master/src/code.c
func (b *hlBuf) float64() float64 {
	var res float64
	*b = (*b)[8:]
	return res
}

// The index type is encoded to store integers
// in 1,2 or 4 bytes.  The encoding requires
// the data to be read in big endian notation.
func (b *hlBuf) index() int {
	var i int
	c := (*b)[0]

	if c&0x80 == 0 {
		return int(b.byte())
	}

	if (c & 0x40) == 0 {
		i = int(binary.BigEndian.Uint16(*b) & 0x1fff)
		*b = (*b)[2:]
		if c&0x20 == 0 {
			return i
		} else {
			return -i
		}
	}

	i = int(binary.BigEndian.Uint32(*b) & 0x1fffffff)
	*b = (*b)[4:]

	if (c & 0x20) == 0 {
		return i
	} else {
		return -i
	}
}

type hlNative struct {
	lib    string
	name   string
	t      int
	findex int
}

type hlFunction struct {
	hlType int
	fIndex int
	lReg   []int
	lOp    []byte

	/*
		int findex;
		int nregs;
		int nops;
		hl_type *type;
		hl_type **regs;
		hl_opcode *ops;
		int *debug;

		hl_type_obj *obj;
		const uchar *field;
	*/
}

type hldump struct {
	version     int
	flags       int
	lInt        []int
	lFloat      []float64
	lString     []string
	lDebugFiles []string
	lType       []hlType
	lGlobal     []int
	lNative     []hlNative
	lFunction   []*hlFunction
	entryPoint  int
}

func (h *hldump) init(b hlBuf) error {
	// Verify existence of the magic HLB identifier
	if hlMagic != string(b[0:3]) {
		return ErrNotValidHLB
	}
	b.skip(len(hlMagic))

	// Bail on fast on unsupported HLB version
	h.version = int(b.byte())
	if h.version != 2 {
		return ErrUnsupported
	}

	h.flags = int(b.index())
	nInt := int(b.index())
	nFloat := int(b.index())
	nString := int(b.index())
	nType := int(b.index())
	nGlobal := int(b.index())
	nNative := int(b.index())
	nFunction := int(b.index())
	h.entryPoint = int(b.index())

	fmt.Printf("Version: %d\nFlags: %x\n", h.version, h.flags)
	fmt.Printf("Ints: %d\nFloats: %d\nStrings: %d\nTypes: %d\nGlobals: %d\nNatives: %d\nFunctions: %d\n", nInt, nFloat, nString, nType, nGlobal, nNative, nFunction)

	h.lInt = make([]int, nInt)
	for i := 0; i < nInt; i++ {
		h.lInt[i] = int(b.int32())
		//fmt.Printf("Int%d: %d\n", i, h.lInt[i])
	}

	h.lFloat = make([]float64, nFloat)
	for i := 0; i < nFloat; i++ {
		h.lFloat[i] = b.float64()
		//fmt.Printf("Float%d: %f\n", i, h.lFloat[i])
	}

	h.lString = make([]string, nString)
	skip := int(b.int32())
	tmpBuf := b[:skip]
	b.skip(skip)
	for i := 0; i < nString; i++ {
		sz := b.index()
		h.lString[i] = string(tmpBuf[:sz])
		tmpBuf = tmpBuf[sz+1:]
		//fmt.Printf("[%.3d] %.50s\n", i, h.lString[i])
	}

	if h.HasDebug() {
		nDebugFiles := int(b.index())
		h.lDebugFiles = make([]string, nDebugFiles)
		skip := int(b.int32())
		tmpBuf := b[:skip]
		b.skip(skip)
		for i := 0; i < nDebugFiles; i++ {
			sz := b.index()
			h.lDebugFiles[i] = string(tmpBuf[:sz])
			tmpBuf = tmpBuf[sz+1:]
		}
	}

	h.lType = make([]hlType, nType)
	for i := 0; i < nType; i++ {
		//fmt.Printf("%.5d ", i)
		h.lType[i] = h.readType(&b)
	}

	h.lGlobal = make([]int, nGlobal)
	for i := 0; i < nGlobal; i++ {
		h.lGlobal[i] = int(b.index()) // get_type
	}

	h.lNative = make([]hlNative, nNative)
	for i := 0; i < nNative; i++ {
		h.lNative[i].lib = string(b.byte())
		h.lNative[i].name = string(b.byte())
		h.lNative[i].t = int(b.index()) // get_type
		h.lNative[i].findex = int(b.index())
	}

	h.lFunction = make([]*hlFunction, nFunction)
	for i := 0; i < nFunction; i++ {
		h.lFunction[i] = h.readFunction(&b)
		if h.HasDebug() {
		}
	}

	return nil
}

func (h *hldump) HasDebug() bool {
	return h.flags&1 == 1
}

func (h *hldump) readFunction(b *hlBuf) *hlFunction {
	f := new(hlFunction)
	f.hlType = b.index() // get_type
	f.fIndex = b.index()
	fmt.Printf("New func [%d] %s\n", f.fIndex, h.lString[f.hlType])
	nReg := b.index()
	nOp := b.index()
	f.lReg = make([]int, nReg)
	for i := 0; i < nReg; i++ {
		f.lReg[i] = b.index() // get_type
	}
	for i := 0; i < nOp; i++ {
	}
	return nil
	/*
		int i;
		f->type = hl_get_type(r);
		f->findex = UINDEX();
		f->nregs = UINDEX();
		f->nops = UINDEX();
		f->regs = (hl_type**)hl_malloc(&r->code->falloc, f->nregs * sizeof(hl_type*));
		for(i=0;i<f->nregs;i++)
		f->regs[i] = hl_get_type(r);
		CHK_ERROR();
		f->ops = (hl_opcode*)hl_malloc(&r->code->falloc, f->nops * sizeof(hl_opcode));
		for(i=0;i<f->nops;i++)
		hl_read_opcode(r, f, f->ops+i);
	*/
}

func (h *hldump) readType(b *hlBuf) hlType {
	kind := hlTypeKind(b.byte())
	fmt.Printf("> %s\n", kind)
	switch kind {
	/*
		case HVOID:
		case HUI8:
		case HUI16:
		case HI32:
		case HI64:
		case HF32:
		case HF64:
		case HBOOL:
		case HBYTES:
		case HDYN:
	*/
	case HFUN:
		t := new(hlTypeFun)
		nArgs := int(b.byte())
		t.args = make([]int, nArgs)
		for i := 0; i < nArgs; i++ {
			t.args[i] = int(b.index()) // get_type
			fmt.Printf("Arg%d: %d\n", i, t.args[i])
		}
		t.ret = int(b.index()) // get_type
		fmt.Printf("Ret: %d\n", t.ret)
		return t
		/*
			int nargs = READ();
			t->fun = (hl_type_fun*)hl_malloc(&r->code->alloc,sizeof(hl_type_fun));
			t->fun->nargs = nargs;
			t->fun->args = (hl_type**)hl_malloc(&r->code->alloc,sizeof(hl_type*)*nargs);
			for(i=0;i<nargs;i++)
			t->fun->args[i] = hl_get_type(r);
			t->fun->ret = hl_get_type(r);
		*/
	case HOBJ:
		t := new(hlTypeObj)
		i := int(b.index()) // read_ustring
		fmt.Printf("Name: %s\n", h.lString[i])
		// TODO missing name
		t.super = b.index()
		t.globalValue = b.index()
		nField := b.index()
		nProto := b.index()
		nBinding := b.index()
		for i := 0; i < nField; i++ {
			name := b.index() // read_ustring
			b.index()         // get_type
			fmt.Printf("  Field: %s\n", h.lString[name])
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
		return t
		/*
			{
				int i;
				const uchar *name = hl_read_ustring(r);
				int super = INDEX();
				int global = UINDEX();
				int nfields = UINDEX();
				int nproto = UINDEX();
				int nbindings = UINDEX();
				t->obj = (hl_type_obj*)hl_malloc(&r->code->alloc,sizeof(hl_type_obj));
				t->obj->name = name;
				t->obj->super = super < 0 ? NULL : r->code->types + super;
				t->obj->global_value = (void**)(int_val)global;
				t->obj->nfields = nfields;
				t->obj->nproto = nproto;
				t->obj->nbindings = nbindings;
				t->obj->fields = (hl_obj_field*)hl_malloc(&r->code->alloc,sizeof(hl_obj_field)*nfields);
				t->obj->proto = (hl_obj_proto*)hl_malloc(&r->code->alloc,sizeof(hl_obj_proto)*nproto);
				t->obj->bindings = (int*)hl_malloc(&r->code->alloc,sizeof(int)*nbindings*2);
				t->obj->rt = NULL;
				for(i=0;i<nfields;i++) {
					hl_obj_field *f = t->obj->fields + i;
					f->name = hl_read_ustring(r);
					f->hashed_name = hl_hash_gen(f->name,true);
					f->t = hl_get_type(r);
				}
				for(i=0;i<nproto;i++) {
					hl_obj_proto *p = t->obj->proto + i;
					p->name = hl_read_ustring(r);
					p->hashed_name = hl_hash_gen(p->name,true);
					p->findex = UINDEX();
					p->pindex = INDEX();
				}
				for(i=0;i<nbindings;i++) {
					t->obj->bindings[i<<1] = UINDEX();
					t->obj->bindings[(i<<1)|1] = UINDEX();
				}
			}
		*/
	/*
		case HARRAY:
		case HTYPE:
	*/
	case HREF:
		t := new(hlTypeRef)
		t.tparam = b.index() // get_type
		return t
		/*
			t->tparam = hl_get_type(r);
		*/
	case HVIRTUAL:
		t := new(hlTypeVirtual)
		nField := b.index()
		t.field = make([]hlField, nField)
		for i := 0; i < nField; i++ {
			t.field[i].name = b.index() // read_ustring
			//t.field[i].hash = hl_hash_gen(f->name,true)
			t.field[i].hlType = b.index() // get_type
		}
		return t
		/*
			{
				int i;
				int nfields = UINDEX();
				t->virt = (hl_type_virtual*)hl_malloc(&r->code->alloc,sizeof(hl_type_virtual));
				t->virt->nfields = nfields;
				t->virt->fields = (hl_obj_field*)hl_malloc(&r->code->alloc,sizeof(hl_obj_field)*nfields);
				for(i=0;i<nfields;i++) {
					hl_obj_field *f = t->virt->fields + i;
					f->name = hl_read_ustring(r);
					f->hashed_name = hl_hash_gen(f->name,true);
					f->t = hl_get_type(r);
				}
			}
		*/
	/*
		case HDYNOBJ:
	*/
	case HABSTRACT:
		t := new(hlTypeAbstract)
		t.name = b.index() // read_ustring
		return t
		/*
			t->abs_name = hl_read_ustring(r);
		*/
	case HENUM:
		t := new(hlTypeEnum)
		i := int(b.index()) // read_ustring
		fmt.Printf("Name: %s\n", h.lString[i])
		t.globalValue = int(b.index())
		t.nConstructs = int(b.byte())
		fmt.Printf("enum: %d, globalValue: %d, numConstructs: %d\n", i, t.globalValue, t.nConstructs)
		for i := 0; i < t.nConstructs; i++ {
			name := int(b.index())
			nParam := int(b.index())
			fmt.Printf("Name: %s, Params: %d\n", h.lString[name], nParam)
			for j := 0; j < nParam; j++ {
				ty := int(b.index())
				fmt.Printf("Type: %d\n", ty)
			}
		}
		return t
		/*
			int i,j;
			t->tenum = hl_malloc(&r->code->alloc,sizeof(hl_type_enum));
			t->tenum->name = hl_read_ustring(r);
			t->tenum->global_value = (void**)(int_val)UINDEX();
			t->tenum->nconstructs = READ();
			t->tenum->constructs = (hl_enum_construct*)hl_malloc(&r->code->alloc, sizeof(hl_enum_construct)*t->tenum->nconstructs);
			for(i=0;i<t->tenum->nconstructs;i++) {
				hl_enum_construct *c = t->tenum->constructs + i;
				c->name = hl_read_ustring(r);
				c->nparams = UINDEX();
				c->params = (hl_type**)hl_malloc(&r->code->alloc,sizeof(hl_type*)*c->nparams);
				c->offsets = (int*)hl_malloc(&r->code->alloc,sizeof(int)*c->nparams);
				for(j=0;j<c->nparams;j++)
				c->params[j] = hl_get_type(r);
			}

		*/
	/*
		case HNULL:
	*/
	default:
		if kind >= HLAST {
			log.Fatal("Unknown type ", kind)
		}
		t := new(hlTypeGeneric)
		t.kind = kind
		return t
	}

	return nil
}

func main() {
	fmt.Printf("HL Dump\n")
	f, err := os.Open("data/deadcells.exe")
	f.Seek(369424, io.SeekStart)
	//f, err := os.Open("data/helloworld.hl")
	if err != nil {
		log.Fatal(err)
	}

	fs, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, fs.Size())
	_, err = f.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	dump := new(hldump)
	dump.init(buf)

	/*
		fmt.Printf("Version: %d, Flags: %x\n", version, flags)
		fmt.Printf("Ints: %d, Floats: %d, Strings: %d\n", numInt, numFloat, numString)
		fmt.Printf("Types: %d, Globals: %d, Natives: %d, Funcs: %d\n", numTypes, numGlobals, numNatives, numFunctions)
		fmt.Printf("Entry: %d\n", entry)
	*/
}
