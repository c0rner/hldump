package main

import (
	"fmt"
	"log"
	"os"
)

type hlNative struct {
	lib    string
	name   string
	t      hxType
	findex int
}

type hlFunction struct {
	typePtr hxType
	funIdx  int
	lReg    []hxType
	lInst   []hxilInst

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
type hldumper interface {
	GetType(int) hxType
	GetString(int) string
	HasDebug() bool
}

type hldump struct {
	version    int
	flags      int
	lInt       []int
	lFloat     []float64
	lString    []string
	lDebugFile []string
	lType      []hxType
	lGlobal    []hxType
	lNative    []hlNative
	lFunction  []*hlFunction
	entryPoint int
}

func (h *hldump) GetString(i int) string {
	if i < 0 || i > len(h.lString) {
		return ""
	}
	return h.lString[i]
}

func (h *hldump) GetType(i int) hxType {
	if i < 0 || i > len(h.lType) {
		return nil
	}
	return h.lType[i]
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

	h.flags = b.index()
	nInt := b.index()
	nFloat := b.index()
	nString := b.index()
	nType := b.index()
	nGlobal := b.index()
	nNative := b.index()
	nFunction := b.index()
	h.entryPoint = b.index()

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
		nDebugFile := b.index()
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
		h.lType[i] = readType(h, &b)
	}

	h.lGlobal = make([]hxType, nGlobal)
	for i := 0; i < nGlobal; i++ {
		h.lGlobal[i] = h.GetType(b.index())
	}

	h.lNative = make([]hlNative, nNative)
	for i := 0; i < nNative; i++ {
		h.lNative[i].lib = h.GetString(b.index())
		h.lNative[i].name = h.GetString(b.index())
		h.lNative[i].t = h.GetType(b.index())
		h.lNative[i].findex = b.index()
		//fmt.Printf("Lib: %s, %s, %d\n", h.lNative[i].lib, h.lNative[i].name, h.lNative[i].findex)
	}

	h.lFunction = make([]*hlFunction, nFunction)
	for i := 0; i < nFunction; i++ {
		h.lFunction[i] = readFunction(h, &b)
	}

	return nil
}

func (h *hldump) HasDebug() bool {
	return h.flags&1 == 1
}

func readInstruction(b *hlBuf, inst *hxilInst) error {
	op := hxilOp(b.byte())
	if int(op) >= len(hxilOpCodes) {
		// BAD OP CODE FIXME ADD ERROR
		return nil
	}
	//fmt.Printf("Opcode: %.1x (%s)\n", op, hxilOpCodes[op].name)
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
			log.Fatal("Not implemented!")
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
	return nil
}

func readDebugInfo(b *hlBuf, nOp int) {
	//var file int
	var line int
	for i := 0; i < nOp; {
		c := int(b.byte())
		switch {
		case (c & 1) == 1:
			//file = ((c << 7) & 0x7f00) + int(b.byte())
			b.byte()
		case (c & 2) == 2:
			delta := c >> 6
			count := (c >> 2) & 0xf
			for j := 0; j < count; j++ {
				i++
			}
			line += delta
		case (c & 4) == 4:
			line += c >> 3
			i++
		default:
			b.byte()
			b.byte()
			i++
		}
	}
}

func readFunction(c hldumper, b *hlBuf) *hlFunction {
	f := new(hlFunction)
	f.typePtr = c.GetType(b.index())
	f.funIdx = b.index()
	nReg := b.index()
	nInst := b.index()
	fmt.Printf("New func [%d] %T (%d,%d)\n", f.funIdx, f.typePtr, nReg, nInst)
	f.lReg = make([]hxType, nReg)
	for i := 0; i < nReg; i++ {
		f.lReg[i] = c.GetType(b.index())
		fmt.Printf("%d: %T\n", i, f.lReg[i])
	}
	f.lInst = make([]hxilInst, nInst)
	for i := 0; i < nInst; i++ {
		readInstruction(b, &f.lInst[i])
	}

	if c.HasDebug() {
		readDebugInfo(b, nInst) // Use lInst?
	}
	return nil
}

func readType(c hldumper, b *hlBuf) hxType {
	typeId := hxitId(b.byte())
	t := NewHXType(typeId)
	if t == nil {
		log.Fatalf("Unknown type: %d\n", typeId)
	}

	switch t := t.(type) {
	//case *hxtVoid:
	//case *hxtUI8:
	//case *hxtUI16:
	//case *hxtI32:
	//case *hxtI64:
	//case *hxtF32:
	//case *hxtF64:
	//case *hxtBool:
	//case *hxtBytes:
	//case *hxtDyn:
	case *hxtFun:
		t.Unmarshal(c, b)
	case *hxtObj:
		t.Unmarshal(c, b)
	//case *hxtArray:
	//case *hxtType:
	case *hxtRef:
		t.Unmarshal(c, b)
	case *hxtVirtual:
		t.Unmarshal(c, b)
	//case *hxtDynObj:
	case *hxtAbstract:
		t.Unmarshal(c, b)
	case *hxtEnum:
		t.Unmarshal(c, b)
	case *hxtNull:
		t.Unmarshal(c, b)
	default:
	}

	return t
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
