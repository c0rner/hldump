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

	hlb, err := hl.NewData(FindHLB(buf))
	if err != nil {
		log.Fatal(err)
	}

	hlb.Resolve()
	hlb.Dump()
}
