package msnrbf

type binaryTypeEnumeration byte

const (
	binaryPrimitive binaryTypeEnumeration = iota
	binaryString
	binaryObject
	binarySystemClass
	binaryClass
	binaryObjectArray
	binaryStringArray
	binaryPrimitiveArray
)
