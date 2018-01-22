package msnrbf

import (
	"io"
	"net/rpc"
	"sync"
)

type clientCodec struct {
	mu   sync.Mutex
	conn io.ReadWriteCloser
}

func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &clientCodec{
		conn: conn,
	}
}

func (c *clientCodec) WriteRequest(r *rpc.Request, v interface{}) error {
	c.mu.Lock()

	return nil
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	return nil
}

func (c *clientCodec) ReadResponseBody(v interface{}) error {
	return nil
}

func (c *clientCodec) Close() error {
	return nil
}
