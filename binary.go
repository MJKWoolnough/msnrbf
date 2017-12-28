package msnrbf

type binaryTypeEnumeration byte

const (
	binaryTypePrimitive binaryTypeEnumeration = iota
	binaryTypeString
	binaryTypeObject
	binaryTypeSystemClass
	binaryTypeClass
	binaryTypeObjectArray
	binaryTypeStringArray
	binaryTypePrimitiveArray
)

func (r *reader) ReadBinaryTypeEnumeration() binaryTypeEnumeration {
	return binaryTypeEnumeration(r.ReadByte())
}

func (b binaryTypeEnumeration) String() string {
	switch b {
	case binaryTypePrimitive:
		return "primitive"
	case binaryTypeString:
		return "string"
	case binaryTypeObject:
		return "object"
	case binaryTypeSystemClass:
		return "system class"
	case binaryTypeClass:
		return "type class"
	case binaryTypeObjectArray:
		return "object array"
	case binaryTypeStringArray:
		return "string array"
	case binaryTypePrimitiveArray:
		return "primitive array"
	default:
		return "unknown binary type"
	}
}

type binaryArrayTypeEnumeration byte

const (
	binaryArrayTypeSingle binaryArrayTypeEnumeration = iota
	binaryArrayTypeJagged
	binaryArrayTypeRectangular
	binaryArrayTypeSingleOffset
	binaryArrayTypeJaggedOffset
	binaryArrayTypeRectangularOffset
)

func (r *reader) ReadBinaryArrayTypeEnumeration() binaryArrayTypeEnumeration {
	return binaryArrayTypeEnumeration(r.ReadByte())
}

func (b binaryArrayTypeEnumeration) String() string {
	switch b {
	case binaryArrayTypeSingle:
		return "single"
	case binaryArrayTypeJagged:
		return "jagged"
	case binaryArrayTypeRectangular:
		return "rectangular"
	case binaryArrayTypeSingleOffset:
		return "single offset"
	case binaryArrayTypeJaggedOffset:
		return "jagged offset"
	case binaryArrayTypeRectangularOffset:
		return "rectangular offset"
	default:
		return "unknown binary array type"
	}
}
