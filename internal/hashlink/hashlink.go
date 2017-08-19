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
	strings    StringContainer
	types      []hlType
	globals    []hlType
	natives    []*hlNative
	functions  []*hlFunction
	funLookup  []int
	debugFiles []LineFile
}

func (d *Data) LookupInt(i int) int       { return d.ints[i] }
func (d *Data) LookupFloat(i int) float64 { return d.floats[i] }

//func (d *Data) LookupString(i int) []byte { return d.strings[i] }
func (d *Data) LookupType(i int) hlType   { return d.types[i] }
func (d *Data) LookupGlobal(i int) hlType { return d.globals[i] }
func (d *Data) LookupFun(i int) int       { return d.funLookup[i] }
func (d *Data) LookupFunction(i int) *hlFunction {
	if i < len(d.functions) {
		return d.functions[i]
	}
	return nil
}
func (d *Data) LookupNative(i int) *hlNative {
	if i < len(d.natives) {
		fmt.Printf("Native: %d\n", d.natives[i])
		return d.natives[i]
	}
	return nil
}

func (d *Data) Dump() {
	for _, t := range d.types {
		switch t := t.(type) {
		case *ObjType:
			var extends = "none"
			if t.super > 0 {
				super := d.LookupType(t.super).(*ObjType)
				extends = d.strings.GoString(super.nameIdx)
			}
			fmt.Printf("@%d Class: %s, Global: %d, Extends: %s\n", 0, d.strings.String(t.nameIdx), t.global, extends)
			fmt.Printf("\t%d fields\n", len(t.lField))
			for j := range t.lField {
				fmt.Printf("\t\t@%d %s %T\n", j, d.strings.String(t.lField[j].nameIdx), t.lField[j].typePtr)
			}
			fmt.Printf("\t%d methods\n", len(t.lProto))
			for j := range t.lProto {
				var fname string
				idx := d.LookupFun(t.lProto[j].funIdx)
				if idx < len(d.functions) {
					fun := d.LookupFunction(idx)
					fname = fmt.Sprintf("%T", fun)
				} else {
					fun := d.LookupNative(idx)
					fname = d.strings.GoString(fun.libIdx) + "." + d.strings.GoString(fun.nameIdx)
				}

				//fmt.Printf("\t\t@%d %s fun@%d() %T\n", j, t.lProto[j].name, t.lProto[j].fIndex, fun.typePtr)
				fmt.Printf("\t\t@%d %s fun@%d() %s\n", j, d.strings.String(t.lProto[j].nameIdx), t.lProto[j].funIdx, fname)
			}
			fmt.Printf("\t%d bindings\n", len(t.lBinding))
			for j := range t.lBinding {
				fmt.Printf("\t\t@%d %d fun@%d\n", j, t.lBinding[j].funIdx, t.lBinding[j].funIdx)
			}
		}
	}
}

type hhlString []byte

type StringContainer struct {
	data  []byte
	index [][]byte
}

func (s *StringContainer) Append(b []byte) {
	s.data = append(s.data, b...)
	s.index = append(s.index, s.data[len(s.data)-len(b):])
}

func (s *StringContainer) String(i int) []byte {
	if i >= len(s.index) {
		return nil
	}
	return s.index[i]
}

func (s *StringContainer) GoString(i int) string {
	return string(s.String(i))
	return string(s.String(i))
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
	nStrings := b.index()
	d.types = make([]hlType, b.index())
	d.globals = make([]hlType, b.index())
	d.natives = make([]*hlNative, b.index())
	d.functions = make([]*hlFunction, b.index())
	d.funLookup = make([]int, len(d.natives)+len(d.functions))
	d.entryPoint = b.index()

	for i := range d.ints {
		d.ints[i] = int(b.int32())
	}

	fmt.Printf("Version: %d\nFlags: %x\n", d.version, d.flags)
	fmt.Printf("Ints: %d\nFloats: %d\nStrings: %d\nTypes: %d\nGlobals: %d\nNatives: %d\nFunctions: %d\n",
		len(d.ints), len(d.floats), nStrings, len(d.types), len(d.globals), len(d.natives), len(d.functions))

	for i := range d.floats {
		d.floats[i] = b.float64()
	}

	skip := int(b.int32())
	tmpBuf := b[:skip]
	b.skip(skip)
	for i := 0; i < nStrings; i++ {
		sz := b.index()
		d.strings.Append(tmpBuf[:sz])
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
		n := new(hlNative)
		n.libIdx = b.index()
		n.nameIdx = b.index()
		n.t = d.LookupType(b.index())
		n.funIdx = b.index()
		d.natives[i] = n
		d.funLookup[n.funIdx] = i + len(d.functions)
		//fmt.Printf("Lib: %s, %s, %d=%d\n", d.strings.String(n.libIdx), d.strings.String(n.nameIdx), n.funIdx, i+len(d.functions))
	}

	for i := range d.functions {
		f := new(hlFunction)
		f.typePtr = d.LookupType(b.index())
		f.funIdx = b.index()
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
		d.functions[i] = f
		d.funLookup[f.funIdx] = i
	}

	return d, nil
}

func readType(ctx *Data, b *hlbStream) hlType {
	typeId := HdtId(b.byte())
	t := typeId.NewType()
	if t == nil {
		log.Fatalf("Unknown type: %d\n", typeId)
	}

	switch t := t.(type) {
	//case *VoidType:
	//case *UI8Type:
	//case *UI16Type:
	//case *I32Type:
	//case *I64Type:
	//case *F32Type:
	//case *F64Type:
	//case *BoolType:
	//case *BytesType:
	//case *DynType:
	case *FunType:
		t.Unmarshal(ctx, b)
	case *ObjType:
		t.Unmarshal(ctx, b)
	//case *ArrayType:
	//case *TypeType:
	case *RefType:
		t.Unmarshal(ctx, b)
	case *VirtualType:
		t.Unmarshal(ctx, b)
	//case *DynObjType:
	case *AbstractType:
		t.Unmarshal(ctx, b)
	case *EnumType:
		t.Unmarshal(ctx, b)
	case *NullType:
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
