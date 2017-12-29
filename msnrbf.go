package msnrbf

import (
	"fmt"
	"io"

	"github.com/MJKWoolnough/errors"
)

type File struct {
	serializedStreamHeader
	library interface{}
}

func Open(ra io.ReaderAt) (*File, error) {
	r := newReader(ra)
	if r.ReadByte() != 0 {
		return nil, ErrInvalidHeaderByte
	}
	sh := r.ReadSerializedStreamHeader()
	if sh.MajorVersion != 1 || sh.MinorVersion != 0 {
		return nil, ErrInvalidVersion
	}
	typ := r.ReadRecordTypeEnumeration()
	if r.Err != nil {
		return nil, r.Err
	}
	switch typ {
	case recordClassWithID:
		r.ReadClassWithID()
	case recordSystemClassWithMembers:
		r.ReadSystemClassWithMembers()
	case recordClassWithMembers:
		r.ReadClassWithMembers()
	case recordSystemClassWithMembersAndTypes:
		r.ReadSystemClassWithMembersAndTypes()
	case recordClassWithMembersAndTypes:
		r.ReadClassWithMembersAndTypes()
	case recordBinaryObjectString:
		r.ReadBinaryObjectString()
	case recordBinaryArray:
		r.ReadBinaryArray()
	case recordBinaryLibrary:
		r.ReadBinaryLibrary()
	case recordArraySinglePrimitive:
		r.ReadArraySinglePrimitive()
	case recordArraySingleObject:
		r.ReadArraySingleObject()
	case recordArraySingleString:
		r.ReadArraySingleString()
	case recordMethodCall:
		r.ReadMethodCall()
	case recordMethodReturn:
		r.ReadMethodReturn()
	default:
		return nil, fmt.Errorf("invalid type: %d", typ)
	}
	return &File{}, nil
}

var (
	ErrInvalidHeaderByte errors.Error = "invalid RecordTypeEnum: expecting 0"
	ErrInvalidVersion    errors.Error = "invalid version"
)
