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
