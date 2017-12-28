package msnrbf

type classInfo struct {
	ObjectID    int32
	Name        string
	MemberNames []string
}

func (r *reader) ReadClassInfo() classInfo {
	var c classInfo
	c.ObjectID = r.ReadInt32()
	c.Name = r.ReadString()
	memberCount := r.ReadInt32()
	if memberCount < 0 {
		r.SetError(ErrInvalidLength)
		return classInfo{}
	}
	c.MemberNames = make([]string, memberCount)
	for n := range c.MemberNames {
		c.MemberNames[n] = r.ReadString()
	}
	return c
}

type classTypeInfo struct {
	TypeName  string
	LibraryID int32
}

func (r *reader) ReadClassTypeInfo() classTypeInfo {
	var c classTypeInfo
	c.TypeName = r.ReadString()
	c.LibraryID = r.ReadInt32()
	return c
}

type memberTypeInfo struct {
	BinaryTypeEnums []binaryTypeEnumeration
	AdditionalInfos []interface{}
}

func (r *reader) ReadMemberTypeInfo(l uint32) memberTypeInfo {
	m := memberTypeInfo{
		BinaryTypeEnums: make([]binaryTypeEnumeration, l),
		AdditionalInfos: make([]interface{}, l),
	}

	for n := range m.BinaryTypeEnums {
		m.BinaryTypeEnums[n] = r.ReadBinaryTypeEnumeration()
		switch m.BinaryTypeEnums[n] {
		case binaryTypePrimitive:
			m.AdditionalInfos[n] = r.ReadPrimitiveTypeEnum()
		case binaryTypeSystemClass:
			m.AdditionalInfos[n] = r.ReadString()
		case binaryTypeClass:
			m.AdditionalInfos[n] = r.ReadClassTypeInfo()
		case binaryTypePrimitiveArray:
			m.AdditionalInfos[n] = r.ReadPrimitiveTypeEnum()
		}
	}
	return m
}
