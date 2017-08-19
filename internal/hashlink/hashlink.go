package hashlink

import (
	"fmt"
	"log"
)

const (
	Magic = "HLB"
)

type Data struct {
	version    int
	flags      Flags
	entryPoint int
	ints       []int
	floats     []float64
	strings    []hlString
	types      []hlType
	globals    []hlType
	natives    []*hlNative
	functions  []*hlFunction
	debugFiles []LineFile
}

func (d *Data) LookupInt(i int) int         { return d.ints[i] }
func (d *Data) LookupFloat(i int) float64   { return d.floats[i] }
func (d *Data) LookupString(i int) hlString { return d.strings[i] }
func (d *Data) LookupType(i int) hlType     { return d.types[i] }
func (d *Data) LookupGlobal(i int) hlType   { return d.globals[i] }
func (d *Data) LookupFunction(i int) *hlFunction {
	if i > len(d.natives)+len(d.functions) {
		log.Fatal("Function out of range!")
	}
	if i < len(d.functions) {
		return d.functions[i]
	}
	return nil
}

type LineFile string

func NewData(b hlbStream) (*Data, error) {
	d := new(Data)

	// Verify existence of the magic HLB identifier
	if Magic != string(b[0:3]) {
		return nil, ErrNotValidHLB
	}
	b.skip(len(Magic))

	// Bail on fast on unsupported HLB version
	d.version = int(b.byte())
	if d.version != 2 {
		return nil, ErrUnsupported
	}

	d.flags = Flags(b.index())
	d.ints = make([]int, b.index())
	d.floats = make([]float64, b.index())
	d.strings = make([]hlString, b.index())
	d.types = make([]hlType, b.index())
	d.globals = make([]hlType, b.index())
	d.natives = make([]*hlNative, b.index())
	d.functions = make([]*hlFunction, b.index())
	d.entryPoint = b.index()

	for i := range d.ints {
		d.ints[i] = int(b.int32())
	}

	fmt.Printf("Version: %d\nFlags: %x\n", d.version, d.flags)
	fmt.Printf("Ints: %d\nFloats: %d\nStrings: %d\nTypes: %d\nGlobals: %d\nNatives: %d\nFunctions: %d\n",
		len(d.ints), len(d.floats), len(d.strings), len(d.types), len(d.globals), len(d.natives), len(d.functions))

	for i := range d.floats {
		d.floats[i] = b.float64()
	}

	skip := int(b.int32())
	tmpBuf := b[:skip]
	b.skip(skip)
	for i := range d.strings {
		sz := b.index()
		copy(d.strings[i], tmpBuf[:sz])
		tmpBuf = tmpBuf[sz+1:]
	}

	if d.flags.HasDebug() {
		nDebugFile := b.index()
		d.debugFiles = make([]LineFile, nDebugFile)
		skip := int(b.int32())
		tmpBuf := b[:skip]
		b.skip(skip)
		for i := range d.debugFiles {
			sz := b.index()
			d.debugFiles[i] = LineFile(tmpBuf[:sz])
			tmpBuf = tmpBuf[sz+1:]
		}
	}

	for i := range d.types {
		d.types[i] = readType(d, &b)
	}

	for i := range d.globals {
		d.globals[i] = d.LookupType(b.index())
		//fmt.Printf("\t@%d %T\n", i, d.lGlobal[i])
	}

	for i := range d.natives {
		d.natives[i].lib = d.LookupString(b.index())
		d.natives[i].name = d.LookupString(b.index())
		d.natives[i].t = d.LookupType(b.index())
		d.natives[i].findex = b.index()
		//fmt.Printf("Lib: %s, %s, %d\n", d.lNative[i].lib, d.lNative[i].name, d.lNative[i].findex)
	}

	for i := range d.functions {
		f := d.functions[i]
		f.typePtr = d.LookupType(b.index())
		f.fIndex = b.index()
		nReg := b.index()
		nInst := b.index()

		//fmt.Printf("New func [%d] %T (%d,%d)\n", f.funIdx, f.typePtr, nReg, nInst)
		f.lReg = make([]hlType, nReg)
		for i := 0; i < nReg; i++ {
			f.lReg[i] = d.LookupType(b.index())
		}
		f.lInst = make([]hxilInst, nInst)
		for i := 0; i < nInst; i++ {
			readInstruction(&b, &f.lInst[i])
		}

		if d.flags.HasDebug() {
			readDebugInfo(&b, nInst) // Use lInst?
		}
	}

	return d, nil
}

func readType(ctx *Data, b *hlbStream) hlType {
	typeId := HdtId(b.byte())
	t := typeId.New()
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
		t.Unmarshal(ctx, b)
	case *hxtObj:
		t.Unmarshal(ctx, b)
	//case *hxtArray:
	//case *hxtType:
	case *hxtRef:
		t.Unmarshal(ctx, b)
	case *hxtVirtual:
		t.Unmarshal(ctx, b)
	//case *hxtDynObj:
	case *hxtAbstract:
		t.Unmarshal(ctx, b)
	case *hxtEnum:
		t.Unmarshal(ctx, b)
	case *hxtNull:
		t.Unmarshal(ctx, b)
	default:
	}

	return t
}

func readInstruction(b *hlbStream, inst *hxilInst) error {
	inst.op = hxilOp(b.byte())
	if int(inst.op) >= len(hxilOpCodes) {
		return ErrBadOpCode
	}

	nArg := hxilOpCodes[inst.op].args
	switch nArg {
	case 0:
		return nil
	case -1:
		inst.arg = make([]int, 3)
		switch inst.op {
		case hxilCallN, hxilCallClosure, hxilCallMethod,
			hxilCallThis, hxilMakeEnum:
			inst.arg[0] = b.index()
			inst.arg[1] = b.index()
			inst.arg[2] = int(b.byte())
			inst.extra = make([]int, inst.arg[2])
			for i := 0; i < inst.arg[2]; i++ {
				inst.extra[i] = b.index()
			}
		case hxilSwitch:
			inst.arg[0] = b.index()
			inst.arg[1] = b.index()
			inst.extra = make([]int, inst.arg[1])
			for i := 0; i < inst.arg[1]; i++ {
				inst.extra[i] = b.index()
			}
			inst.arg[2] = b.index()
		default:
			log.Fatal("Not implemented!")
		}
	default:
		inst.arg = make([]int, nArg)
		for i := 0; i < nArg; i++ {
			inst.arg[i] = b.index()
		}
	}
	return nil
}

func readDebugInfo(b *hlbStream, nOp int) {
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
