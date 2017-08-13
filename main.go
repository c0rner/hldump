package main

import (
	"fmt"
	"log"
	"os"
)

type hlNative struct {
	lib    string
	name   string
	t      int
	findex int
}

type hlFunction struct {
	typeIdx int
	funIdx  int
	lReg    []int
	lOp     []byte

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
	version    int
	flags      int
	lInt       []int
	lFloat     []float64
	lString    []string
	lDebugFile []string
	lType      []hxType
	lGlobal    []int
	lNative    []hlNative
	lFunction  []*hlFunction
	entryPoint int
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
	}

	h.lFloat = make([]float64, nFloat)
	for i := 0; i < nFloat; i++ {
		h.lFloat[i] = b.float64()
	}

	h.lString = make([]string, nString)
	skip := int(b.int32())
	tmpBuf := b[:skip]
	b.skip(skip)
	for i := 0; i < nString; i++ {
		sz := b.index()
		h.lString[i] = string(tmpBuf[:sz])
		tmpBuf = tmpBuf[sz+1:]
	}

	if h.HasDebug() {
		nDebugFile := int(b.index())
		h.lDebugFile = make([]string, nDebugFile)
		skip := int(b.int32())
		tmpBuf := b[:skip]
		b.skip(skip)
		for i := 0; i < nDebugFile; i++ {
			sz := b.index()
			h.lDebugFile[i] = string(tmpBuf[:sz])
			tmpBuf = tmpBuf[sz+1:]
		}
	}

	h.lType = make([]hxType, nType)
	for i := 0; i < nType; i++ {
		h.lType[i] = h.readType(&b)
	}

	h.lGlobal = make([]int, nGlobal)
	for i := 0; i < nGlobal; i++ {
		h.lGlobal[i] = h.VerifyType(int(b.index())) // FIXME: get_type
	}

	h.lNative = make([]hlNative, nNative)
	for i := 0; i < nNative; i++ {
		h.lNative[i].lib = h.lString[int(b.index())]
		h.lNative[i].name = h.lString[int(b.index())]
		h.lNative[i].t = h.VerifyType(int(b.index())) // FIXME: get_type
		h.lNative[i].findex = int(b.index())
		fmt.Printf("Lib: %s, %s, %d\n", h.lNative[i].lib, h.lNative[i].name, h.lNative[i].findex)
	}

	h.lFunction = make([]*hlFunction, nFunction)
	for i := 0; i < nFunction; i++ {
		h.lFunction[i] = h.readFunction(&b)
	}

	return nil
}

// THIS SHOULD BE DONE SOME OTHER WAY
func (h *hldump) VerifyType(t int) int {
	if t < 0 || t > len(h.lType) {
		log.Fatal("Fucking shit! ", t)
	}
	return t
}

func (h *hldump) HasDebug() bool {
	return h.flags&1 == 1
}

func (h *hldump) readOpCode(b *hlBuf) {
	op := hxilOp(b.byte())
	if int(op) >= len(hxilOpCodes) {
		log.Fatal("Bad OP CODE: ", int(op))
	}
	fmt.Printf("Opcode: %.1x (%s)\n", op, hxilOpCodes[op].name)
	switch hxilOpCodes[op].args {
	case 0:
	case 1:
		b.index()
	case 2:
		b.index()
		b.index()
	case 3:
		b.index()
		b.index()
		b.index()
	case 4:
		b.index()
		b.index()
		b.index()
		b.index()
	case -1:
		switch op {
		case hxilCallN, hxilCallClosure, hxilCallMethod,
			hxilCallThis, hxilMakeEnum:
			b.index()
			b.index()
			i := b.byte()
			for ; i > 0; i-- {
				b.index()
			}
		case hxilSwitch:
			b.index()
			i := b.index()
			for ; i > 0; i-- {
				b.index()
			}
			b.index()
		default:
			log.Fatal("Default op shit!")
		}
	default:
		size := hxilOpCodes[op].args - 3
		b.index()
		b.index()
		b.index()
		for i := 0; i < size; i++ {
			b.index()
		}
	}
}

func (h *hldump) readDebugInfo(b *hlBuf, nOp int) {
	var file int
	var line int
	for i := 0; i < nOp; {
		c := int(b.byte())
		fmt.Printf("Debug: 0x%.02x\n", c)
		switch {
		case (c & 1) == 1:
			file = ((c << 7) & 0x7f00) + int(b.byte())
			fmt.Printf("File: %s\n", h.lDebugFile[file])
		case (c & 2) == 2:
			delta := c >> 6
			count := (c >> 2) & 0xf
			for j := 0; j < count; j++ {
				fmt.Printf("Op: %d Line: %d\n", i, line)
				i++
			}
			line += delta
		case (c & 4) == 4:
			line += c >> 3
			fmt.Printf("Op: %d Line: %d\n", i, line)
			i++
		default:
			b.byte()
			b.byte()
			fmt.Printf("Op: %d Line: %d\n", i, line)
			i++
		}
	}
}

func (h *hldump) readFunction(b *hlBuf) *hlFunction {
	f := new(hlFunction)
	f.typeIdx = h.VerifyType(b.index()) // FIXME: get_type
	f.funIdx = b.index()
	nReg := b.index()
	nOp := b.index()
	fmt.Printf("New func [%d] %s (%d,%d)\n", f.funIdx, h.lType[f.typeIdx].Kind(), nReg, nOp)
	f.lReg = make([]int, nReg)
	for i := 0; i < nReg; i++ {
		f.lReg[i] = h.VerifyType(b.index()) // FIXME: get_type
		fmt.Printf("%d: %s\n", i, h.lType[f.lReg[i]].Kind())
	}
	for i := 0; i < nOp; i++ {
		h.readOpCode(b)
	}

	if h.HasDebug() {
		h.readDebugInfo(b, nOp)
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

func (h *hldump) readType(b *hlBuf) hxType {
	typeId := hxitId(b.byte())
	switch typeId {
	case hxitFun: // FUN
		t := new(hxTypeFun)
		t.typeId = typeId
		nArgs := int(b.byte())
		t.args = make([]int, nArgs)
		for i := 0; i < nArgs; i++ {
			t.args[i] = h.VerifyType(int(b.index())) // FIXME: get_type
		}
		t.ret = h.VerifyType(int(b.index())) // FIXME: get_type
		return t
	case hxitObj: // OBJ TODO
		t := new(hxTypeObj)
		t.typeId = typeId
		t.name = h.lString[int(b.index())] // read_ustring
		//fmt.Printf("Name: %s\n", t.name)
		t.super = b.index() // FIXME - NULL super
		t.globalValue = b.index()
		nField := b.index()
		nProto := b.index()
		nBinding := b.index()

		t.field = make([]hlField, nField)
		for i := 0; i < nField; i++ {
			t.field[i].name = h.lString[b.index()] // read_ustring
			//t.field[i].hash = hl_hash_gen(f->name,true)
			t.field[i].typeIdx = h.VerifyType(b.index()) // FIXME: get_type
			//fmt.Printf("  Field: %s\n", t.field[i].name)
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
	case hxitRef: // REF
		t := new(hxTypeRef)
		t.typeId = typeId
		t.tparam = h.VerifyType(b.index()) // FIXME: get_type
		return t
	case hxitVirtual: // VIRTUAL
		t := new(hxTypeVirtual)
		t.typeId = typeId
		nField := b.index()
		t.field = make([]hlField, nField)
		for i := 0; i < nField; i++ {
			t.field[i].name = h.lString[b.index()] // read_ustring
			//t.field[i].hash = hl_hash_gen(f->name,true) // TODO
			t.field[i].typeIdx = h.VerifyType(b.index()) // FIXME: get_type
		}
		return t
	case hxitAbstract: // ABSTRACT
		t := new(hxTypeAbstract)
		t.typeId = typeId
		t.name = h.lString[b.index()] // read_ustring
		return t
	case hxitEnum: // ENUM TODO
		t := new(hxTypeEnum)
		t.typeId = typeId
		t.name = h.lString[int(b.index())] // read_ustring
		t.globalValue = int(b.index())
		t.nConstructs = int(b.byte())
		//fmt.Printf("Enum name: '%s', globalValue: %d, numConstructs: %d\n", t.name, t.globalValue, t.nConstructs)
		for i := 0; i < t.nConstructs; i++ {
			b.index() // read_ustring
			nParam := int(b.index())
			//fmt.Printf("Name: %s, Params: %d\n", h.lString[name], nParam)
			for j := 0; j < nParam; j++ {
				h.VerifyType(b.index()) // FIXME: get_type
				//fmt.Printf("Type: %d\n", tidx)
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
	case hxitNull: // NULL
		t := new(hxTypeNull)
		t.typeId = typeId
		t.tparam = h.VerifyType(b.index()) // FIXME: get_type
		return t
	default:
		if typeId >= hxitLast {
			log.Fatal("Unknown type ", typeId)
		}
		t := new(hxTypeGeneric)
		t.typeId = typeId
		return t
	}

	return nil
}

func ScanHLB(b []byte) hlBuf {
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(hlMagic); j++ {
			if b[i+j] != hlMagic[j] {
				break
			}
			if j == len(hlMagic)-1 {
				return hlBuf(b[i:])
			}
		}
	}
	return nil
}

func main() {
	fmt.Printf("HL Dump\n")
	f, err := os.Open("data/helloworld.hl")
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

	b := ScanHLB(buf)
	if b == nil {
		log.Fatal("No Hashlink detected.\n")
	}
	dump := new(hldump)
	err = dump.init(b)
	if err != nil {
		log.Fatal(err)
	}
}
