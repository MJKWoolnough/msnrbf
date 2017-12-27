package msnrbf

import (
	"fmt"
	"io"

	"github.com/MJKWoolnough/errors"
)

type File struct {
	RootID int32
}

func Open(r io.ReaderAt) (*File, error) {
	rs := newReader(r)
	if rs.ReadByte() != 0 {
		return nil, ErrInvalidHeaderByte
	}
	rootID := rs.ReadInt32()
	rs.SkipInt32() // headerId
	majorVer := rs.ReadInt32()
	minorVer := rs.ReadInt32()
	if majorVer != 1 || minorVer != 0 {
		return nil, ErrInvalidVersion
	}
	for {
		typ := recordTypeEnumeration(rs.ReadByte())
		if rs.Err != nil {
			return nil, rs.Err
		}
		switch typ {
		//case recordClassWithId:
		//case recordSystemClassWithMembers:
		//case recordClasswithMembers:
		//case recordSystemClassWithMembersAndTypes:
		//case recordClassWithMembersAndTypes:
		//case recordBinaryObjectString:
		//case recordBinaryArray:
		//case recordMemberPrimitiveTyped:
		//case recordMemberReference:
		//case recordObjectNull:
		//case recordMessageEnd:
		case recordBinaryLibrary:
			rs.ReadBinaryLibrary()
		//case recordObjectNullMultiple256:
		//case recordObjectNullMultiple:
		//case recordArraySinglePrimitive:
		//case recordArraySingleObject:
		//case recordArraySingleString:
		//case recordMethodCall:
		//case recordMethodReturn:
		default:
			return nil, fmt.Errorf("unhandled type: %d", typ)
		}
	}
	return &File{
		RootID: rootID,
	}, nil
}

var (
	ErrInvalidHeaderByte errors.Error = "invalid RecordTypeEnum: expecting 0"
	ErrInvalidVersion    errors.Error = "invalid version"
)
