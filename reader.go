package msnrbf

import (
	"io"
	"unicode/utf8"

	"github.com/MJKWoolnough/byteio"
	"github.com/MJKWoolnough/errors"
)

type reader struct {
	*byteio.StickyLittleEndianReader
	rs *readSeeker
}

type readSeeker struct {
	io.ReaderAt
	pos int64
}

func (r *readSeeker) Read(p []byte) (int, error) {
	n, err := r.ReadAt(p, r.pos)
	r.pos += int64(n)
	return n, err
}

func newReader(r io.ReaderAt) reader {
	nr := reader{
		rs: &readSeeker{ReaderAt: r},
	}
	nr.StickyLittleEndianReader = &byteio.StickyLittleEndianReader{Reader: nr.rs}
	return nr
}

const maxString = 16 * 1024 * 1024

func (r *reader) ReadString() string {
	length := r.ReadLength()
	if length == 0 {
		return ""
	} else if length > maxString {
		r.SetError(ErrStringTooLong)
		return ""
	}
	b := make([]byte, length)
	_, err := io.ReadFull(r, b)
	if err != nil {
		r.Err = err
		return ""
	}
	if !utf8.Valid(b) {
		r.SetError(ErrInvalidString)
		return ""
	}
	return string(b)
}

func (r *reader) ReadVarInt() int32 {
	var n int32
	for i := 0; i < 5; i++ {
		b := r.ReadByte()
		n |= int32(b&127) << uint(7*i)
		if b&128 == 0 {
			break
		}
	}
	if n < 0 {
		r.SetError(ErrInvalidVarInt)
	}
	return n
}

func (r *reader) RealBool() bool {
	switch r.ReadByte() {
	case 0:
		return false
	case 1:
		return true
	}
	r.SetError(ErrInvalidBool)
	return false
}

func (r *reader) ReadByte() byte {
	return r.ReadUint8()
}

func (r *reader) Goto(n uint32) {
	r.rs.pos = int64(n)
}

func (r *reader) SetError(err error) {
	if r.Err == nil {
		r.Err = err
	}
}

func (r *reader) Skip(n uint32) {
	r.rs.pos += int64(n)
}

func (r *reader) SkipByte() {
	r.Skip(1)
}

func (r *reader) SkipInt32() {
	r.Skip(4)
}

func (r *reader) SkipUint32() {
	r.Skip(4)
}

func (r *reader) SkipFloat32() {
	r.Skip(4)
}

func (r *reader) SkipString() {
	r.Skip(r.ReadUint32())
}

// Errors
const (
	ErrInvalidString errors.Error = "string is invalid"
	ErrInvalidVarInt errors.Error = "invalid variable integer"
	ErrInvalidLength errors.Error = "invalid length"
	ErrStringTooLong errors.Error = "string exceeds maximum length"
	ErrInvalidSeek   errors.Error = "invalid seek"
	ErrInvalidBool   errors.Error = "invalid boolean"
)
