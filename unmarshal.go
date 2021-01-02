package msnrbf

import (
	"errors"
	"io"
	"time"
)

type call struct {
	*reader
	objects map[int32]interface{}
}

func newCall(r *reader) call {
	return call{
		reader:  r,
		objects: make(map[int32]interface{}),
	}
}

type class struct {
	classInfo
	*memberTypeInfo
	libraryID int32
}

type classWithData struct {
	class
	data []interface{}
}

func (c *call) readClass(te recordTypeEnumeration) classWithData {
	var (
		cl class
		id int32
	)
	switch te {
	case recordClassWithID:
		cid := c.ReadClassWithID()
		o, ok := c.objects[cid.MetadataID]
		if !ok {
			c.SetError(ErrInvalidObject)
			return classWithData{}
		}
		cd, ok := o.(classWithData)
		if !ok {
			c.SetError(ErrInvalidObject)
			return classWithData{}
		}
		cl = cd.class
		id = cid.ObjectID
	case recordClassWithMembers:
		cwm := c.ReadClassWithMembers()
		id = cwm.ClassInfo.ObjectID
		cl = class{
			classInfo: cwm.ClassInfo,
			libraryID: cwm.LibraryID,
		}
	case recordClassWithMembersAndTypes:
		cwmt := c.ReadClassWithMembersAndTypes()
		id = cwmt.ClassInfo.ObjectID
		cl = class{
			classInfo:      cwmt.ClassInfo,
			memberTypeInfo: &cwmt.MemberTypeInfo,
			libraryID:      cwmt.LibraryID,
		}
	case recordSystemClassWithMembers:
		scwm := c.ReadSystemClassWithMembers()
		id = scwm.ClassInfo.ObjectID
		cl = class{
			classInfo: scwm.ClassInfo,
			libraryID: -1,
		}
	case recordSystemClassWithMembersAndTypes:
		scwmt := c.ReadSystemClassWithMembersAndTypes()
		id = scwmt.ClassInfo.ObjectID
		cl = class{
			classInfo:      scwmt.ClassInfo,
			memberTypeInfo: &scwmt.MemberTypeInfo,
			libraryID:      -1,
		}
	default:
		c.SetError(ErrInvalidRecord)
		return classWithData{}
	}
	cwd := classWithData{
		cl,
		c.readMemberReferences(len(cl.MemberNames), func(n int) (binaryTypeEnumeration, primitiveTypeEnum) {
			bt := binaryTypeEnumeration(255)
			ai := primitiveTypeEnum(255)
			if cl.memberTypeInfo != nil {
				bt = cl.memberTypeInfo.BinaryTypeEnums[n]
			}
			switch bt {
			case binaryTypePrimitive:
				ai = cl.memberTypeInfo.AdditionalInfos[n].(primitiveTypeEnum)
			}
			return bt, ai
		}),
	}
	c.objects[id] = cwd
	return cwd
}

var id = 0

func (c *call) readArray(te recordTypeEnumeration) {
	switch te {
	case recordArraySingleObject:
		so := c.ReadArraySingleObject()
		values := c.readMemberReferences(int(pa.ArrayInfo.Length), func(_ int) (binaryTypeEnumeration, primitiveTypeEnum) {
			return 255, 255
		})
		_ = values
	case recordArraySinglePrimitive:
		pa := c.ReadArraySinglePrimitive()
		values := c.readPrimitiveArray(pa.PrimitiveTypeEnum)
		_ = values
	case recordArraySingleString:
		ass := c.ReadArraySingleString()
		values := make([]interface{}, ass.ArrayInfo.Length)
		for i := int32(0); i < ass.ArrayInfo.Length; i++ {
			switch te := c.ReadRecordTypeEnumeration(); te {
			case recordBinaryObjectString:
				values[i] = c.ReadBinaryObjectString()
			case recordMemberReference:
				values[i] = c.ReadMemberReference()
			case recordObjectNull:
				c.ReadObjectNull()
			case recordObjectNullMultiple256:
				i += int32(c.ReadObjectNullMultiple256()) - 1
			case recordObjectNullMultiple:
				i += int32(c.ReadObjectNullMultiple()) - 1
			default:
				c.SetError(ErrInvalidRecord)
			}
		}
		_ = values
	case recordBinaryArray:
		ba := c.ReadBinaryArray()
		length := 0
		for _, l := range ba.Lengths {
			length += int(l)
		}
		for _, l := range ba.LowerBounds {
			length -= int(l)
		}
		var values interface{}
		switch ba.TypeEnum {
		case binaryTypePrimitive, binaryTypePrimitiveArray:
			values = c.readPrimitiveArray(ba.AdditionTypeInfo.(primitiveTypeEnum))
		default:
			values = c.readMemberReferences(length, func(_ int) (binaryTypeEnumeration, primitiveTypeEnum) {
				return ba.TypeEnum, ai
			})
			_ = values
		}
	}
}

func (c *call) readPrimitiveArray(pte primitiveTypeEnum) interface{} {
	switch pte {
	case primitiveTypeBoolean:
		data := make([]bool, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadBool()
		}
		return data
	case primitiveTypeByte:
		data := make([]uint8, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadUint8()
		}
		return data
	case primitiveTypeChar:
		data := make([]rune, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadChar()
		}
		return data
	case primitiveTypeDecimal:
		data := make([]string, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadDecimal()
		}
		return data
	case primitiveTypeDouble:
		data := make([]float64, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadFloat64()
		}
		return data
	case primitiveTypeInt16:
		data := make([]int16, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadInt16()
		}
		return data
	case primitiveTypeInt32:
		data := make([]int32, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadInt32()
		}
		return data
	case primitiveTypeInt64:
		data := make([]int64, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadInt64()
		}
		return data
	case primitiveTypeSByte:
		data := make([]int8, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadInt8()
		}
		return data
	case primitiveTypeSingle:
		data := make([]float32, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadFloat32()
		}
		return data
	case primitiveTypeTimeSpan:
		data := make([]time.Duration, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadTimeSpan()
		}
		return data
	case primitiveTypeDateTime:
		data := make([]time.Time, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadDateTime()
		}
		return data
	case primitiveTypeUInt16:
		data := make([]uint16, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadUint16()
		}
		return data
	case primitiveTypeUInt32:
		data := make([]uint32, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadUint32()
		}
		return data
	case primitiveTypeUInt64:
		data := make([]uint64, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadUint64()
		}
		return data
	case primitiveTypeNull:
		return make([]struct{}, pa.ArrayInfo.Length)
	case primitiveTypeString:
		data := make([]string, pa.ArrayInfo.Length)
		for n := range data {
			data[n] = c.ReadString()
		}
		return data
	}
	return nil
}

func (c *call) readMemberReferences(length int, typeInfo func(n int) (binaryTypeEnumeration, primitiveTypeEnum)) []interface{} {
	data := make([]interface{}, length)
	for i := 0; i < length; i++ {
		mti, pte := typeInfo(i)
		switch mti {
		case binaryTypePrimitive:
			switch pte {
			case primitiveTypeBoolean:
				data[i] = c.ReadBool()
			case primitiveTypeByte:
				data[i] = c.ReadUint8()
			case primitiveTypeChar:
				data[i] = c.ReadChar()
			case primitiveTypeDecimal:
				data[i] = c.ReadDecimal()
			case primitiveTypeDouble:
				data[i] = c.ReadFloat64()
			case primitiveTypeInt16:
				data[i] = c.ReadInt16()
			case primitiveTypeInt32:
				data[i] = c.ReadInt32()
			case primitiveTypeInt64:
				data[i] = c.ReadInt64()
			case primitiveTypeSByte:
				data[i] = c.ReadInt8()
			case primitiveTypeSingle:
				data[i] = c.ReadFloat32()
			case primitiveTypeTimeSpan:
				data[i] = c.ReadTimeSpan()
			case primitiveTypeDateTime:
				data[i] = c.ReadDateTime()
			case primitiveTypeUInt16:
				data[i] = c.ReadUint16()
			case primitiveTypeUInt32:
				data[i] = c.ReadUint32()
			case primitiveTypeUInt64:
				data[i] = c.ReadUint64()
			case primitiveTypeNull:
			case primitiveTypeString:
				data[i] = c.ReadString()
			}
		//case binaryTypePrimitiveArray:
		//case binaryTypeSystemClass:
		//case binaryTypeClass:
		default:
			switch te := c.ReadRecordTypeEnumeration(); te {
			case recordMemberReference:
				data[i] = c.ReadMemberReference()
			case recordBinaryObjectString:
				data[i] = c.ReadBinaryObjectString()
			case recordObjectNull:
				c.ReadObjectNull()
			case recordObjectNullMultiple256:
				i += int(c.ReadObjectNullMultiple256()) - 1
			case recordObjectNullMultiple:
				i += int(c.ReadObjectNullMultiple()) - 1
			case recordBinaryLibrary:
				b := c.ReadBinaryLibrary()
				c.objects[b.LibraryID] = b
				te = c.ReadRecordTypeEnumeration()
				fallthrough
			case recordClassWithID, recordClassWithMembers, recordClassWithMembersAndTypes, recordSystemClassWithMembers, recordSystemClassWithMembersAndTypes:
				data[i] = c.readClass(te)
			default:
				c.SetError(ErrInvalidRecord)
				return nil
			}
		}
	}
	return data
}

// Unmarshal parses nrbf encoded data
func Unmarshal(r io.Reader) error {
	sr := newReader(r)
	if sr.ReadRecordTypeEnumeration() != recordSerializedStreamHeader {
		return ErrInvalidRecord
	}
	header := sr.ReadSerializedStreamHeader()
	if header.MajorVersion != 1 || header.MinorVersion != 0 {
		return ErrInvalidVersion
	}
	c := newCall(sr)
	hadCR := false
Loop:
	for {
		te := sr.ReadRecordTypeEnumeration()
		hasBL := false
		if te == recordBinaryLibrary {
			b := sr.ReadBinaryLibrary()
			c.objects[b.LibraryID] = b
			hasBL = true
			te = sr.ReadRecordTypeEnumeration()
		}
		switch te {
		case recordClassWithID, recordClassWithMembers, recordClassWithMembersAndTypes, recordSystemClassWithMembers, recordSystemClassWithMembersAndTypes:
			c.readClass(te)
		case recordArraySingleObject, recordArraySinglePrimitive, recordArraySingleString, recordBinaryArray:
			c.readArray(te)
		case recordBinaryObjectString:
			if hasBL {
				return ErrInvalidRecord
			}
			bos := sr.ReadBinaryObjectString()
			c.objects[bos.ObjectID] = bos
		case recordMethodCall:
			if hadCR {
				return ErrInvalidRecord
			}
			hadCR = true
			// TODO
		case recordMethodReturn:
			if hadCR {
				return ErrInvalidRecord
			}
			hadCR = true
			// TODO
		case recordMessageEnd:
			if hasBL {
				return ErrInvalidRecord
			}
			break Loop
		default:
			return ErrInvalidRecord
		}
	}
	return nil
}

// errors
var (
	ErrInvalidRecord  = errors.New("invalid record")
	ErrInvalidVersion = errors.New("invalid version")
	ErrInvalidObject  = errors.New("invalid object")
)
