package bitsio

import (
	"unicode/utf8"
)

type BitWriter struct {
	buf    []byte
	ptr    int
	bitPtr byte
}

func NewBitWriter() *BitWriter {
	return &BitWriter{buf: []byte{0x00}}
}

func (bw *BitWriter) WriteBit(bit bool) {
	if bit {
		offset := 7 - bw.bitPtr
		bw.buf[bw.ptr] |= (1 << offset)
	}
	bw.moveBitPtr()
}

func (bw *BitWriter) moveBitPtr() {
	bw.bitPtr++
	if bw.bitPtr >= 8 {
		bw.buf = append(bw.buf, 0x00)
		bw.ptr++
		bw.bitPtr = 0
	}
}

func (bw *BitWriter) WriteByte(b byte) {
	for i := range 8 {
		offset := 7 - i
		val := b & (1 << offset)
		bw.WriteBit(!(val == 0))
	}
}

func (bw *BitWriter) WtiteRune(r rune) {
	if uint32(r) < utf8.RuneSelf {
		bw.WriteByte(byte(r))
		return
	}
	b := make([]byte, 4)
	n := utf8.EncodeRune(b, r)
	for i := range b[:n] {
		bw.WriteByte(b[i])
	}
}

func (bw *BitWriter) BitsLeftToByte() byte {
	return 8 - bw.bitPtr
}

func (bw *BitWriter) Bytes() []byte {
	if bw.bitPtr == 0 {
		return bw.buf[:len(bw.buf)-1]
	}
	return bw.buf
}

type BitReader struct {
	buf    []byte
	ptr    int
	bitPtr byte
}

func NewBitReader(b []byte) *BitReader {
	return &BitReader{buf: b}
}

func (br *BitReader) ReadBit() bool {
	// if br.IsEmpty() {
	// 	return false, fmt.Errorf("is empty")
	// }
	offset := 7 - br.bitPtr
	val := br.buf[br.ptr] & (1 << offset)
	br.moveBitPtr()
	return !(val == 0)
	// return !(val == 0), nil
}

func (br *BitReader) moveBitPtr() {
	br.bitPtr++
	if br.bitPtr >= 8 {
		br.ptr++
		br.bitPtr = 0
	}
}

func (br *BitReader) ReadByte() byte {
	bitWriter := NewBitWriter()
	for range 8 {
		bitWriter.WriteBit(br.ReadBit())
	}
	return bitWriter.Bytes()[0]
}

func (br *BitReader) ReadRune() rune {
	b := make([]byte, 0)
	b = append(b, br.ReadByte())
	for !utf8.Valid(b) {
		b = append(b, br.ReadByte())
	}
	r, _ := utf8.DecodeRune(b)
	return r
}

func (br *BitReader) FinishByte() {
	if br.bitPtr != 0 {
		br.ptr++
		br.bitPtr = 0
	}
}

func (br *BitReader) NumReadByte() int {
	return br.ptr
}

func (br *BitReader) IsLastByte() bool {
	return br.ptr+1 == len(br.buf)
}

func (br *BitReader) IsEmpty() bool {
	return br.ptr >= len(br.buf)
}
