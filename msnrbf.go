package msnrbf

import (
	"io"
	"net"
	"net/rpc"
)

func Dial(network, address string) (*rpc.Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn), nil
}

func NewClient(conn io.ReadWriteCloser) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn))
}

func ServeConn(conn io.ReadWriteCloser) {
	rpc.ServeConn(NewServerCodec(conn))
}
