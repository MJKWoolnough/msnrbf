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

type binaryArrayTypeEnumeration byte

const (
	binaryArrayTypeSingle binaryArrayTypeEnumeration = iota
	binaryArrayTypeJagged
	binaryArrayTypeRectangular
	binaryArrayTypeSingleOffset
	binaryArrayTypeJaggedOffset
	binaryArrayTypeRectangularOffset
)
