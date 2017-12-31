package msnrbf

import "io"

type serverCodec struct {
	conn io.ReadWriteCloser
}

func NewServerCodec(conn io.ReadWriteCloser) rpc.NewServerCodec {
	return &serverCodec{
		conn: conn,
	}
}
