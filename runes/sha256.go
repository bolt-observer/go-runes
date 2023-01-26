package runes

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"os"
)

// The purpose of all this is to extract (and be able to set) SHA-256 midstates

const (
	// ChunkSize is chunk size in bytes
	ChunkSize = 64
	// OutputSize is size of SHA256 checksum in bytes.
	OutputSize = 32
	// Magic256 is the magic for SHA256
	Magic256 = "sha\x03"
)

// Sha256 struct
type Sha256 struct {
	hasher hash.Hash
	len    uint64
}

// Marshaller is the interface for marshalling
type Marshaller interface {
	UnmarshalBinary([]byte) error
	MarshalBinary() ([]byte, error)
}

// MidState struct
type MidState struct {
	H   [8]uint32
	Len uint64
}

// GetSum gets SHA-256 hash
func (state *MidState) GetSum() [OutputSize]byte {
	var digest [OutputSize]byte

	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint32(digest[i*4:], state.H[i])
	}

	return digest
}

// NewSha256 - construct new instance
func NewSha256() *Sha256 {
	return &Sha256{hasher: sha256.New(), len: uint64(0)}
}

// Write - add bytes
func (s *Sha256) Write(p []byte) (nn int, err error) {
	s.len += uint64(len(p))
	return s.hasher.Write(p)
}

// Reset - reset instance
func (s *Sha256) Reset() {
	s.len = uint64(0)
	s.hasher.Reset()
}

// GetMidState - get the internal state
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

	b = b[len(Magic256):]
	for i := 0; i < 8; i++ {
		b, ret.H[i] = consumeUint32(b)
	}

	b = b[len(b)-8:]
	_, ret.Len = consumeUint64(b)

	return ret
}

// SetMidState - updates internal state
func (s *Sha256) SetMidState(state *MidState) error {
	if state == nil {
		return os.ErrInvalid
	}

	m := s.hasher.(Marshaller)
	if m == nil {
		return nil
	}

	s.Reset()

	marshaledSize := len(Magic256) + 8*4 + ChunkSize + 8
	b := make([]byte, 0, marshaledSize)
	b = append(b, Magic256...)

	for i := 0; i < 8; i++ {
		b = appendUint32(b, state.H[i])
	}

	// Critical assumption here is there are no leftover bytes
	b = appendUint64(b, state.Len)

	// Extend
	b = b[:marshaledSize]

	err := m.UnmarshalBinary(b)
	if err != nil {
		return err
	}

	return nil
}

// SetLen sets the internal length
func (s *Sha256) SetLen(len uint64) {
	s.len = len
}

// GetLen gets the internal length
func (s *Sha256) GetLen() uint64 {
	return s.len
}

// AddPadding adds necessary padding at the end of a chunk
func (s *Sha256) AddPadding() error {
	l := s.len

	// Padding. Add a 1 bit and 0 bits until 56 bytes mod 64.
	var tmp [64 + 8]byte // padding + length buffer
	tmp[0] = 0x80
	var t uint64
	if l%64 < 56 {
		t = 56 - l%64
	} else {
		t = 64 + 56 - l%64
	}

	// Length in bits.
	l <<= 3
	padlen := tmp[:t+8]

	binary.BigEndian.PutUint64(padlen[t+0:], l)
	//fmt.Printf("T is %d, len is %d, l is %d,  padlen %s len %d\n", t, s.len, l, hex.EncodeToString(padlen), len(padlen))

	_, err := s.hasher.Write(padlen)
	if err != nil {
		return err
	}

	s.len += uint64(len(padlen))

	return nil
}

// GetSum returns the SHA-256 hash
func (s *Sha256) GetSum() [OutputSize]byte {
	ret := s.GetMidState().GetSum()
	return ret
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
