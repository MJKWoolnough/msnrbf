package msnrbf

type primitiveTypeEnum byte

const (
	primitiveTypeBoolean primitiveTypeEnum = iota + 1
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

func (p primitiveTypeEnum) String() string {
	switch p {
	case primitiveTypeBoolean:
		return "boolean"
	case primitiveTypeByte:
		return "byte"
	case primitiveTypeChar:
		return "char"
	case primitiveTypeDecimal:
		return "decimal"
	case primitiveTypeDouble:
		return "double"
	case primitiveTypeInt16:
		return "int16"
	case primitiveTypeInt32:
		return "int32"
	case primitiveTypeInt64:
		return "int64"
	case primitiveTypeSByte:
		return "int8"
	case primitiveTypeSingle:
		return "float32"
	case primitiveTypeTimeSpan:
		return "duration"
	case primitiveTypeDateTime:
		return "datetime"
	case primitiveTypeUInt16:
		return "uint16"
	case primitiveTypeUInt32:
		return "uint32"
	case primitiveTypeUInt64:
		return "uint64"
	case primitiveTypeNull:
		return "null"
	case primitiveTypeString:
		return "string"
	default:
		return "unknown primitive"
	}
}
