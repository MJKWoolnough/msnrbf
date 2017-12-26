package msnrbf

type primitiveTypeEnum byte

const (
	primitiveTypeBoolean primitiveTypeEnum = iota
	primitiveTypeByte
	primitiveTypeChar
	_
	primitiveTypeDecimal
	primitiveTypeDouble
	primitiveTypeInt16
	primitiveTypeInt32
	primitiveTypeInt64
	primitiveTypeSByte
	primitiveTypeSingle
	primitiveTypeTimeSpan
	primitiveTypeDateTime
	primitiveTypeUInt16
	primitiveTypeUInt32
	primitiveTypeUInt64
	primitiveTypeNull
	primitiveTypeString
)

func (r *reader) ReadPrimitiveTypeEnum() primitiveTypeEnum {
	return primitiveTypeEnum(r.ReadByte())
}
