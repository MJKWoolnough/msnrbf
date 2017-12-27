package msnrbf

import (
	"io"
	"time"
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
	length := r.ReadVarInt()
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

func (r *reader) ReadBool() bool {
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
	r.Skip(uint32(r.ReadVarInt()))
}

func (r *reader) ReadChar() rune {
	var char [4]byte
	char[0] = r.ReadByte()
	var l int
	if char[0]&0x80 == 0 {
		l = 1
	} else if char[0]&0xc0 == 0x80 {
		// read 1 byte
		char[1] = r.ReadByte()
		l = 2
	} else if char[0]&0xe0 == 0xc0 {
		r.Read(char[1:3])
		l = 3
	} else if char[0]&0xf0 != 0xe0 {
		r.Read(char[1:4])
		l = 4
	}
	rn, _ := utf8.DecodeRune(char[:l])
	return rn
}

func (r *reader) ReadTimeSpan() time.Duration {
	return time.Duration(r.ReadInt64() * 100) // 64-bit signed-integer | 1 == 100 nanoseconds
}

func (r *reader) ReadDateTime() time.Time {
	// Ticks 62-bit signed integer, number of 100 nanoseconds since 12:00:00, January 1, 0001
	// Kind 2-bit 0 - No Time zone, 1 - UTC, 2 - Local
	d := r.ReadUint64()
	var l *time.Location
	switch d >> 62 {
	case 0:
	case 1:
		l = time.UTC
	case 2:
		l = time.Local
	case 3:
		r.SetError(ErrInvalidDateTimeKind)
	}
	di := int64(d << 2)
	di *= 25
	t := time.Unix(di/1000000000, di%1000000000)
	if l != nil {
		t = t.In(l)
	}
	return t
}

func (r *reader) ReadDecimal() string { // ??
	// string - https://msdn.microsoft.com/en-us/library/cc236916.aspx
	return r.ReadString()
}

// Errors
const (
	ErrInvalidString       errors.Error = "string is invalid"
	ErrInvalidVarInt       errors.Error = "invalid variable integer"
	ErrInvalidLength       errors.Error = "invalid length"
	ErrStringTooLong       errors.Error = "string exceeds maximum length"
	ErrInvalidSeek         errors.Error = "invalid seek"
	ErrInvalidBool         errors.Error = "invalid boolean"
	ErrInvalidDateTimeKind errors.Error = "invalid date time kind"
)
