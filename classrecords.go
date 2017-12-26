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

type memberTypeInfo []binaryTypeEnumeration

func (r *reader) ReadMemberTypeInfo(l uint32) memberTypeInfo {
	m := make(memberTypeInfo, l)
	for n := range m {
		m[n] = r.ReadByte()
	}
	return m
}
