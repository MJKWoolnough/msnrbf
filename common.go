package msnrbf

type arrayInfo struct {
	ObjectID, Length int32
}

func (r *reader) ReadArrayInfo() arrayInfo {
	var a arrayInfo
	a.ObjectID = r.ReadInt32()
	a.Length = r.ReadInt32()
	return a
}

type messageFlags uint32

func (r *reader) ReadMessageFlags() messageFlags {
	return messageFlags(r.ReadUint32())
}

type stringValueWithCode struct {
	PrimitiveTypeEnum primitiveTypeEnum
	StringValue       string
}

func (r *reader) ReadStringValueWithCode() stringValueWithCode {
	var s stringValueWithCode
	s.PrimitiveTypeEnum = r.ReadPrimitiveTypeEnum()
	s.StringValue = r.ReadString()
	return s
}

type valueWithCode interface{}

func (r *reader) ReadValueWithCode() valueWithCode {
	return nil
}

type arrayOfValueWithCode []valueWithCode

func (r *reader) ReadArrayOfValueWithCode() arrayOfValueWithCode {
	return nil
}
