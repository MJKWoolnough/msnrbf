package msnrbf

import (
	"io"
	"time"
	"unicode/utf8"

	"vimagination.zapto.org/byteio"
	"vimagination.zapto.org/errors"
)

type rtePeeker struct {
	io.Reader
	peekedRTE recordTypeEnumeration
	peeked    bool
}

func (r *rtePeeker) Read(p []byte) (int, error) {
	r.peeked = false
	return r.Reader.Read(p)
}

type reader struct {
	byteio.StickyLittleEndianReader
	rtePeeker rtePeeker
}

func newReader(r io.Reader) *reader {
	rr := &reader{
		rtePeeker: rtePeeker{
			Reader: r,
		},
	}
	rr.StickyLittleEndianReader.Reader = &rr.rtePeeker
	return rr
}

func (r *reader) PeekRTE() recordTypeEnumeration {
	if r.rtePeeker.peeked {
		return r.rtePeeker.peekedRTE
	}
	r.rtePeeker.peekedRTE = r.ReadRecordTypeEnumeration()
	r.rtePeeker.peeked = true
	return r.rtePeeker.peekedRTE
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
		b := r.ReadUint8()
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
	switch r.ReadUint8() {
	case 0:
		return false
	case 1:
		return true
	}
	r.SetError(ErrInvalidBool)
	return false
}

func (r *reader) SetError(err error) {
	if r.Err == nil {
		r.Err = err
	}
}

func (r *reader) ReadChar() rune {
	var char [4]byte
	char[0] = r.ReadUint8()
	var l int
	if char[0]&0x80 == 0 {
		l = 1
	} else if char[0]&0xc0 == 0x80 {
		// read 1 byte
		char[1] = r.ReadUint8()
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
	// Kind 2-bit 0 - No Time zone, 1 - UTC, 2 - Local - !! UNUSED !! Assumed to be 0
	d := r.ReadUint64()
	di := int64(d << 2)
	di *= 25
	t := time.Unix(di/1000000000, di%1000000000)
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
