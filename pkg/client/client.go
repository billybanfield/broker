package client

import (
	"net"

	"github.com/billybanfield/broker/pkg/rpc"
)

type Client struct {
	conn rpc.UnixConnReaderWriterCloser
}

func New(srvAddr string) (*Client, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", srvAddr)
	if err != nil {
		return nil, err
	}
	sock, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: sock,
	}, nil
}

func (s *Client) OpenFile(req rpc.FOpenRequest, resp *rpc.FOpenResponse) error {
	err := rpc.WriteRequest(s.conn, req)
	if err != nil {
		return err
	}
	err = rpc.ReadResponse(s.conn, resp)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()

}
