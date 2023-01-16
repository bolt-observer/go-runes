package runes

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"os"
)

// The purpose of all this is to extract (and be able to set) SHA-256 midstates

const (
	// The size of a chunk in bytes.
	CHUNK_SIZE = 64
	// The size of a SHA256 checksum in bytes.
	OUTPUT_SIZE = 32
	MAGIC256    = "sha\x03"
)

type Sha256 struct {
	hasher hash.Hash
	len    uint64
}

type Marshaller interface {
	UnmarshalBinary([]byte) error
	MarshalBinary() ([]byte, error)
}

type MidState struct {
	H   [8]uint32
	Len uint64
}

func (state *MidState) GetSum() [OUTPUT_SIZE]byte {
	var digest [OUTPUT_SIZE]byte

	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint32(digest[i*4:], state.H[i])
	}

	return digest
}

func NewSha256() *Sha256 {
	return &Sha256{hasher: sha256.New(), len: uint64(0)}
}

func (s *Sha256) Write(p []byte) (nn int, err error) {
	s.len += uint64(len(p))
	return s.hasher.Write(p)
}

func (s *Sha256) Reset() {
	s.len = uint64(0)
	s.hasher.Reset()
}

func (s *Sha256) GetMidState() *MidState {
	m := s.hasher.(Marshaller)
	if m == nil {
		return nil
	}
	b, err := m.MarshalBinary()
	if err != nil {
		return nil
	}

	ret := &MidState{}

	b = b[len(MAGIC256):]
	for i := 0; i < 8; i++ {
		b, ret.H[i] = consumeUint32(b)
	}

	b = b[len(b)-8:]
	_, ret.Len = consumeUint64(b)

	return ret
}

func (s *Sha256) SetMidState(state *MidState) error {
	if state == nil {
		return os.ErrInvalid
	}

	m := s.hasher.(Marshaller)
	if m == nil {
		return nil
	}

	s.Reset()

	marshaledSize := len(MAGIC256) + 8*4 + CHUNK_SIZE + 8
	b := make([]byte, 0, marshaledSize)
	b = append(b, MAGIC256...)

	for i := 0; i < 8; i++ {
		b = appendUint32(b, state.H[i])
	}

	// Critical assumption here is there are no leftover bytes
	b = appendUint64(b, state.Len)

	err := m.UnmarshalBinary(b)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sha256) AddPadding() error {
	len := s.len
	// Padding. Add a 1 bit and 0 bits until 56 bytes mod 64.
	var tmp [64 + 8]byte // padding + length buffer
	tmp[0] = 0x80
	var t uint64
	if len%64 < 56 {
		t = 56 - len%64
	} else {
		t = 64 + 56 - len%64
	}

	// Length in bits.
	len <<= 3
	padlen := tmp[:t+8]
	binary.BigEndian.PutUint64(padlen[t+0:], len)
	_, err := s.hasher.Write(padlen)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sha256) GetSum() [OUTPUT_SIZE]byte {
	return s.GetMidState().GetSum()
}

func appendUint64(b []byte, x uint64) []byte {
	var a [8]byte
	binary.BigEndian.PutUint64(a[:], x)
	return append(b, a[:]...)
}

func appendUint32(b []byte, x uint32) []byte {
	var a [4]byte
	binary.BigEndian.PutUint32(a[:], x)
	return append(b, a[:]...)
}

func consumeUint64(b []byte) ([]byte, uint64) {
	_ = b[7]
	x := uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	return b[8:], x
}

func consumeUint32(b []byte) ([]byte, uint32) {
	_ = b[3]
	x := uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	return b[4:], x
}
