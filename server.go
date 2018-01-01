package msnrbf

import (
	"io"
	"net/rpc"
)

type serverCodec struct {
	conn io.ReadWriteCloser
}

func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	return &serverCodec{
		conn: conn,
	}
}

func (s *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	return nil
}

func (s *serverCodec) ReadRequestBody(v interface{}) error {
	return nil
}

func (s *serverCodec) WriteResponse(r *rpc.Response, v interface{}) error {
	return nil
}

func (s *serverCodec) Close() error {
	return nil
}
