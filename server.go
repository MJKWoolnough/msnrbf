package msnrbf

import (
	"io"
	"net/rpc"
	"sync"
)

type serverCodec struct {
	mu   sync.Mutex
	conn io.ReadWriteCloser
	num  int64
}

func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	return &serverCodec{
		conn: conn,
	}
}

func (s *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	s.mu.Lock()
	r.Seq = s.num
	s.num++
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
