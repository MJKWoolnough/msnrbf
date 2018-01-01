package msnrbf

type recordTypeEnumeration byte

const (
	recordSerializedStreamHeader recordTypeEnumeration = iota
	recordClassWithID
	recordSystemClassWithMembers
	recordClassWithMembers
	recordSystemClassWithMembersAndTypes
	recordClassWithMembersAndTypes
	recordBinaryObjectString
	recordBinaryArray
	recordMemberPrimitiveTyped
	recordMemberReference
	recordObjectNull
	recordMessageEnd
	recordBinaryLibrary
	recordObjectNullMultiple256
	recordObjectNullMultiple
	recordArraySinglePrimitive
	recordArraySingleObject
	recordArraySingleString
	recordMethodCall
	recordMethodReturn
)

func (r *reader) ReadRecordTypeEnumeration() recordTypeEnumeration {
	if r.rtePeeker.peeked {
		r.rtePeeker.peeked = false
		return r.rtePeeker.peekedRTE
	}
	return recordTypeEnumeration(r.ReadByte())
}

type serializedStreamHeader struct {
	RootID, HeaderID, MajorVersion, MinorVersion int32
}

func (r *reader) ReadSerializedStreamHeader() serializedStreamHeader {
	var s serializedStreamHeader
	s.RootID = r.ReadInt32()
	s.HeaderID = r.ReadInt32()
	s.MajorVersion = r.ReadInt32()
	s.MinorVersion = r.ReadInt32()
	return s
}

type classWithID struct {
	ObjectID, MetadataID int32
}

func (r *reader) ReadClassWithID() classWithID {
	var c classWithID
	c.ObjectID = r.ReadInt32()
	c.MetadataID = r.ReadInt32()
	return c
}

type systemClassWithMembers struct {
	ClassInfo classInfo
}

func (r *reader) ReadSystemClassWithMembers() systemClassWithMembers {
	return systemClassWithMembers{r.ReadClassInfo()}
}

type classWithMembers struct {
	ClassInfo classInfo
	LibraryID int32
}

func (r *reader) ReadClassWithMembers() classWithMembers {
	var c classWithMembers
	c.ClassInfo = r.ReadClassInfo()
	c.LibraryID = r.ReadInt32()
	return c
}

type systemClassWithMembersAndTypes struct {
	ClassInfo      classInfo
	MemberTypeInfo memberTypeInfo
}

func (r *reader) ReadSystemClassWithMembersAndTypes() systemClassWithMembersAndTypes {
	var s systemClassWithMembersAndTypes
	s.ClassInfo = r.ReadClassInfo()
	s.MemberTypeInfo = r.ReadMemberTypeInfo(uint32(len(s.ClassInfo.MemberNames)))
	return s
}

type classWithMembersAndTypes struct {
	ClassInfo      classInfo
	MemberTypeInfo memberTypeInfo
	LibraryID      int32
}

func (r *reader) ReadClassWithMembersAndTypes() classWithMembersAndTypes {
	var c classWithMembersAndTypes
	c.ClassInfo = r.ReadClassInfo()
	c.MemberTypeInfo = r.ReadMemberTypeInfo(uint32(len(c.ClassInfo.MemberNames)))
	c.LibraryID = r.ReadInt32()
	return c
}

type binaryObjectString struct {
	ObjectID int32
	Value    string
}

func (r *reader) ReadBinaryObjectString() binaryObjectString {
	var b binaryObjectString
	b.ObjectID = r.ReadInt32()
	b.Value = r.ReadString()
	return b
}

type binaryArray struct {
	ObjectID         int32
	ArrayTypeEnum    binaryArrayTypeEnumeration
	Rank             int32
	Lengths          []int32
	LowerBounds      []int32
	TypeEnum         binaryTypeEnumeration
	AdditionTypeInfo interface{}
}

func (r *reader) ReadBinaryArray() binaryArray {
	var b binaryArray
	b.ObjectID = r.ReadInt32()
	b.ArrayTypeEnum = r.ReadBinaryArrayTypeEnumeration()
	b.Rank = r.ReadInt32()
	b.Lengths = make([]int32, b.Rank)
	for n := range b.Lengths {
		b.Lengths[n] = r.ReadInt32() // >= 0
	}
	switch b.ArrayTypeEnum {
	case binaryArrayTypeSingleOffset, binaryArrayTypeJaggedOffset, binaryArrayTypeRectangularOffset:
		b.LowerBounds = make([]int32, b.Rank)
		for n := range b.LowerBounds {
			b.LowerBounds[n] = r.ReadInt32()
		}
	}
	b.TypeEnum = r.ReadBinaryTypeEnumeration()
	switch b.TypeEnum {
	case binaryTypePrimitive, binaryTypePrimitiveArray:
		b.AdditionTypeInfo = r.ReadPrimitiveTypeEnum()
	case binaryTypeSystemClass:
		b.AdditionTypeInfo = r.ReadString()
	case binaryTypeClass:
		b.AdditionTypeInfo = r.ReadClassInfo()
	}
	return b
}

type memberPrimitiveTyped interface{}

func (r *reader) ReadMemberPrimitiveTyped() memberPrimitiveTyped {
	switch r.ReadPrimitiveTypeEnum() {
	case primitiveTypeBoolean:
		return r.ReadBool()
	case primitiveTypeByte:
		return r.ReadByte()
	case primitiveTypeChar:
		return r.ReadChar()
	case primitiveTypeDecimal:
		return r.ReadDecimal()
	case primitiveTypeDouble:
		return r.ReadFloat64()
	case primitiveTypeInt16:
		return r.ReadInt16()
	case primitiveTypeInt32:
		return r.ReadInt32()
	case primitiveTypeInt64:
		return r.ReadInt64()
	case primitiveTypeSByte:
		return r.ReadInt8()
	case primitiveTypeSingle:
		return r.ReadFloat32()
	case primitiveTypeTimeSpan:
		return r.ReadTimeSpan()
	case primitiveTypeDateTime:
		return r.ReadDateTime()
	case primitiveTypeUInt16:
		return r.ReadUint16()
	case primitiveTypeUInt32:
		return r.ReadUint32()
	case primitiveTypeUInt64:
		return r.ReadUint64()
	}
	// error
	return nil
}

type memberReference int32

func (r *reader) ReadMemberReference() memberReference {
	return memberReference(r.ReadInt32())
}

type objectNull struct{}

func (*reader) ReadObjectNull() objectNull {
	return objectNull{}
}

type messageEnd struct{}

func (*reader) ReadMessageEnd() messageEnd {
	return messageEnd{}
}

type binaryLibrary struct {
	LibraryID   int32
	LibraryName string
}

func (r *reader) ReadBinaryLibrary() binaryLibrary {
	var b binaryLibrary
	b.LibraryID = r.ReadInt32()
	b.LibraryName = r.ReadString()
	return b
}

type objectNullMultiple256 byte

func (r *reader) ReadObjectNullMultiple256() objectNullMultiple256 {
	return objectNullMultiple256(r.ReadByte())
}

type objectNullMultiple int32

func (r *reader) ReadObjectNullMultiple() objectNullMultiple {
	return objectNullMultiple(r.ReadInt32())
}

type arraySinglePrimitive struct {
	ArrayInfo         arrayInfo
	PrimitiveTypeEnum primitiveTypeEnum
}

func (r *reader) ReadArraySinglePrimitive() arraySinglePrimitive {
	var a arraySinglePrimitive
	a.ArrayInfo = r.ReadArrayInfo()
	a.PrimitiveTypeEnum = r.ReadPrimitiveTypeEnum()
	return a
}

type arraySingleObject struct {
	ArrayInfo arrayInfo
}

func (r *reader) ReadArraySingleObject() arraySingleObject {
	var a arraySingleObject
	a.ArrayInfo = r.ReadArrayInfo()
	return a
}

type arraySingleString struct {
	ArrayInfo arrayInfo
}

func (r *reader) ReadArraySingleString() arraySingleString {
	var a arraySingleString
	a.ArrayInfo = r.ReadArrayInfo()
	return a
}

type binaryMethodCall struct {
	MessageEnum                       messageFlags
	MethodName, TypeName, CallContext stringValueWithCode
	Args                              arrayOfValueWithCode
}

func (r *reader) ReadMethodCall() binaryMethodCall {
	var b binaryMethodCall
	b.MessageEnum = r.ReadMessageFlags()
	b.MethodName = r.ReadStringValueWithCode()
	b.TypeName = r.ReadStringValueWithCode()
	b.CallContext = r.ReadStringValueWithCode()
	b.Args = r.ReadArrayOfValueWithCode()
	return b
}

type binaryMethodReturn struct {
	MessageEnum messageFlags
	ReturnValue valueWithCode
	CallContext stringValueWithCode
	Args        arrayOfValueWithCode
}

func (r *reader) ReadMethodReturn() binaryMethodReturn {
	var b binaryMethodReturn
	b.MessageEnum = r.ReadMessageFlags()
	b.ReturnValue = r.ReadValueWithCode()
	b.CallContext = r.ReadStringValueWithCode()
	b.Args = r.ReadArrayOfValueWithCode()
	return b
}
