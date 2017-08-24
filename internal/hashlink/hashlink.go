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
	funcLookup []int
	debugFiles []LineFile
}

func (d *Data) LookupInt(i int) int       { return d.ints[i] }
func (d *Data) LookupFloat(i int) float64 { return d.floats[i] }

//func (d *Data) LookupString(i int) []byte { return d.strings[i] }
func (d *Data) LookupType(i int) hlType   { return d.types[i] }
func (d *Data) LookupGlobal(i int) hlType { return d.globals[i] }
func (d *Data) LookupFunction(i int) Function {
	idx := d.funcLookup[i]
	if idx < len(d.functions) {
		return d.functions[idx]
	}
	idx -= len(d.functions)
	if idx < len(d.natives) {
		return d.natives[idx]
	}
	return nil
}

func (d *Data) Dump() {
	for i := range d.functions {
		f := d.functions[i]
		fmt.Printf("fun@%d", f.funcIdx)
		fobj := d.LookupType(f.typeIdx).(*FunType)
		for j := 0; j < len(fobj.argPtr); j++ {
			fmt.Printf("%s,", fobj.argPtr[j].Id())
		}
		fmt.Printf("%s\n", fobj.retPtr.Id())
		for j := range f.inst {
			f.inst[j].Print(d)
		}
	}
	for i := range d.types {
		switch t := d.types[i].(type) {
		case *ObjType:
			var extIdx int
			if t.superIdx > 0 {
				super := d.LookupType(t.superIdx).(*ObjType)
				extIdx = super.nameIdx
			}
			fmt.Printf("@%d Class: %s, Global: %d, Extends: %s\n", i, t.namePtr, t.global, d.strings.String(extIdx))
			fmt.Printf("\t%d fields\n", len(t.lField))
			for j := range t.lField {
				fmt.Printf("\t\t@%d %s %T\n", j+t.offset, d.strings.String(t.lField[j].nameIdx), d.LookupType(t.lField[j].typeIdx))
			}
			fmt.Printf("\t%d methods\n", len(t.lProto))
			for j := range t.lProto {
				var fname string
				f := d.LookupFunction(t.lProto[j].funcIdx)
				switch f := f.(type) {
				case *hlFunction:
					fname = fmt.Sprintf("%T", f)
				case *hlNative:
					fname = d.strings.String(f.libIdx) + "." + d.strings.String(f.nameIdx)
				}

				//fmt.Printf("\t\t@%d %s fun@%d() %T\n", j, t.lProto[j].name, t.lProto[j].fIndex, fun.typePtr)
				fmt.Printf("\t\t@%d %s fun@%d[%d] %s\n", j, d.strings.String(t.lProto[j].nameIdx), t.lProto[j].funcIdx, t.lProto[j].override, fname)
			}
			fmt.Printf("\t%d bindings\n", len(t.lBinding))
			for j := range t.lBinding {
				var fname string
				f := d.LookupFunction(t.lBinding[j].funcIdx)
				switch f := f.(type) {
				case *hlFunction:
					fname = fmt.Sprintf("%T", f)
				case *hlNative:
					fname = d.strings.String(f.libIdx) + "." + d.strings.String(f.nameIdx)
				}
				fmt.Printf("\t\t@%d %d fun@%d (%s)\n", j, t.lBinding[j].fldIdx, t.lBinding[j].funcIdx, fname)
			}
		}
	}
}

func (d *Data) Resolve() {
	for i := range d.functions {
		f := d.functions[i]
		d.funcLookup[f.funcIdx] = i
	}
	for i := range d.natives {
		n := d.natives[i]
		d.funcLookup[n.funcIdx] = i + len(d.functions)
		n.libPtr = d.strings.String(n.libIdx)
		n.namePtr = d.strings.String(n.nameIdx)
	}

	for i := range d.types {
		switch t := d.types[i].(type) {
		case *FunType:
			t.argPtr = make([]hlType, len(t.argIdx))
			for j := 0; j < len(t.argIdx); j++ {
				t.argPtr[j] = d.LookupType(t.argIdx[j])
			}
			t.retPtr = d.LookupType(t.retIdx)
		case *ObjType:
			t.namePtr = d.strings.Bytes(t.nameIdx)
			// TODO: Init global value
			for j := 0; j < len(t.lProto); j++ {
				p := t.lProto[j]
				f := d.LookupFunction(p.funcIdx).(*hlFunction)
				f.obj = t
				f.field = d.strings.Bytes(p.nameIdx)
			}
		}
	}
}

type StringContainer struct {
	data  []byte
	index [][]byte
}

func (s *StringContainer) Append(b []byte) {
	s.data = append(s.data, b...)
	s.index = append(s.index, s.data[len(s.data)-len(b):])
}

func (s *StringContainer) Bytes(i int) []byte {
	if i >= len(s.index) {
		return nil
	}
	return s.index[i]
}

func (s *StringContainer) String(i int) string {
	return string(s.Bytes(i))
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
	d.funcLookup = make([]int, len(d.natives)+len(d.functions))
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
		n.typeIdx = b.index()
		n.funcIdx = b.index()
		d.natives[i] = n
		//fmt.Printf("Lib: %s, %s, %d=%d\n", d.strings.String(n.libIdx), d.strings.String(n.nameIdx), n.funcIdx, i+len(d.functions))
	}

	for i := range d.functions {
		f := new(hlFunction)
		f.typeIdx = b.index()
		f.funcIdx = b.index()
		nReg := b.index()
		nInst := b.index()

		//fmt.Printf("New func [%d] %T (%d,%d)\n", f.funcIdx, f.typePtr, nReg, nInst)
		f.regIdx = make([]int, nReg)
		for i := 0; i < nReg; i++ {
			f.regIdx[i] = b.index()
		}
		f.inst = make([]HilInst, nInst)
		for i := 0; i < nInst; i++ {
			readInstruction(&b, &f.inst[i])
		}

		if d.flags.HasDebug() {
			readDebugInfo(&b, nInst) // Use lInst?
		}
		d.functions[i] = f
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

func readInstruction(b *hlbStream, inst *HilInst) error {
	inst.op = HilOp(b.byte())
	if int(inst.op) >= len(OpCodes) {
		return ErrBadOpCode
	}

	nArg := OpCodes[inst.op].args
	switch nArg {
	case 0:
		return nil
	case -1:
		inst.arg = make([]int, 3)
		switch inst.op {
		case OpCallN, OpCallClosure, OpCallMethod,
			OpCallThis, OpMakeEnum:
			inst.arg[0] = b.index()
			inst.arg[1] = b.index()
			inst.arg[2] = int(b.byte())
			inst.extra = make([]int, inst.arg[2])
			for i := 0; i < inst.arg[2]; i++ {
				inst.extra[i] = b.index()
			}
		case OpSwitch:
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

// TODO finish implementation
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
