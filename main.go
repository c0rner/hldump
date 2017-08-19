package main

import (
	"fmt"
	"log"
	"os"
)

import (
	hl "github.com/c0rner/hldump/internal/hashlink"
)

func FindHLB(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(hl.Magic); j++ {
			if b[i+j] != hl.Magic[j] {
				break
			}
			if j == len(hl.Magic)-1 {
				return b[i:]
			}
		}
	}
	return nil
}

func main() {
	fmt.Printf("HL Dump\n")

	f, err := os.Open("data/deadcells.exe")
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

	hlb, err := hl.NewData(FindHLB(buf))
	if err != nil {
		log.Fatal(err)
	}

	for i, t := range hlb.Types {
		switch t := t.(type) {
		case *hxtObj:
			var extends string = "none"
			//if t.super != nil {
			//	extends = *t.super.name
			//}
			fmt.Printf("@%d Class: %s, Global: %d, Extends: %s\n", i, string(t.name), t.global, extends)
			fmt.Printf("\t%d fields\n", len(t.lField))
			for j := range t.lField {
				fmt.Printf("\t\t@%d %s %T\n", j, t.lField[j].name, t.lField[j].typePtr)
			}
			fmt.Printf("\t%d methods\n", len(t.lProto))
			for j := range t.lProto {
				fun := hlb.GetFunction(t.lProto[j].fIndex)
				fmt.Printf("\t\t@%d %s fun@%d() %T\n", j, t.lProto[j].name, t.lProto[j].fIndex, fun.typePtr)
			}
			fmt.Printf("\t%d bindings\n", len(t.lBinding))
			for j := range t.lBinding {
				fmt.Printf("\t\t@%d %d fun@%d\n", j, t.lBinding[j].index, t.lBinding[j].fIndex)
			}
		}
	}
}
