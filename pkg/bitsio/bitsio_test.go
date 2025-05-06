package bitsio_test

import (
	"testing"
	"unicode/utf8"

	. "github.com/PritOriginal/cryptolabs-back/pkg/bitsio"
)

func TestBitWriter_WriteBit(t *testing.T) {
	bw := NewBitWriter()
	bw.WriteBit(true)
	if b := bw.Bytes()[0]; b != 0b10000000 {
		t.Fatalf("BitWriter = %b; want %b", b, 0b10000000)
	}
	bw.WriteBit(true)
	if b := bw.Bytes()[0]; b != 0b11000000 {
		t.Fatalf("BitWriter = %b; want %b", b, 0b11000000)
	}
}

func TestBitWriter_WriteByte(t *testing.T) {
	bw := NewBitWriter()
	var want byte = 0x7a
	bw.WriteByte(want)
	if b := bw.Bytes()[0]; b != want {
		t.Fatalf("BitWriter = %b; want %b", b, want)
	}
}

func TestBitWriter_WriteRune_OneByte(t *testing.T) {
	bw := NewBitWriter()
	want := 'a'
	bw.WtiteRune(want)
	r, _ := utf8.DecodeRune(bw.Bytes())
	if r != want {
		t.Fatalf("BitWriter = %v; want %v", r, want)
	}
}

func TestBitWriter_WriteRune_TwoByte(t *testing.T) {
	bw := NewBitWriter()
	want := 'а'
	bw.WtiteRune(want)
	r, _ := utf8.DecodeRune(bw.Bytes())
	if r != want {
		t.Fatalf("BitWriter = %v; want %v", r, want)
	}
}

func TestBitWriter_BitsLeftToByte(t *testing.T) {
	bw := NewBitWriter()
	if bits := bw.BitsLeftToByte(); bits != 8 {
		t.Fatalf("BitWriter = %v; want %v", bits, 8)
	}

	bw.WriteBit(false)
	if bits := bw.BitsLeftToByte(); bits != 7 {
		t.Fatalf("BitWriter = %v; want %v", bits, 7)
	}

	for range 6 {
		bw.WriteBit(true)
	}
	if bits := bw.BitsLeftToByte(); bits != 1 {
		t.Fatalf("BitWriter = %v; want %v", bits, 1)
	}

	bw.WriteBit(false)
	if bits := bw.BitsLeftToByte(); bits != 8 {
		t.Fatalf("BitWriter = %v; want %v", bits, 8)
	}
}
