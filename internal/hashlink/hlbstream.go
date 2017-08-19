package hashlink

import (
	"encoding/binary"
)

// Hashlink Byte Stream
type hlbStream []byte

// Skip will advance the stream ptr by i bytes
func (b *hlbStream) skip(i int) {
	if i > len(*b) {
		i = len(*b)
	}
	*b = (*b)[i:]
}

// byte returns the next available byte and advances stream ptr
func (b *hlbStream) byte() byte {
	res := (*b)[0]
	*b = (*b)[1:]
	return res
}

func (b *hlbStream) int32() int32 {
	res := int32(binary.LittleEndian.Uint32(*b))
	*b = (*b)[4:]
	return res
}

// Inconclusive how doubles are handled by Hashlib
// For reference see hl_read_double()
// https://github.com/HaxeFoundation/hashlink/blob/master/src/code.c
func (b *hlbStream) float64() float64 {
	var res float64
	// FIXME
	*b = (*b)[8:]
	return res
}

// The index type is encoded to store integers
// in 1,2 or 4 bytes.  The encoding requires
// the data to be read in big endian notation.
func (b *hlbStream) index() int {
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
