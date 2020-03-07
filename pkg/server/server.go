package server

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/billybanfield/broker/pkg/rpc"
)

type Server struct {
	Addr string

	listener *net.UnixListener
}

func (srv *Server) OpenFile(req rpc.FOpenRequest, resp *rpc.FOpenResponse) {
	f, err := os.Open(req.Filename)
	if err != nil {
		log.Printf("%v", err)
	}
	resp.FileDescriptor = f.Fd()
}

func (srv *Server) ListenAndServe() error {
	unixAddr, err := net.ResolveUnixAddr("unix", srv.Addr)
	if err != nil {
		return err
	}
	ln, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		return err
	}
	connChan := make(chan *net.UnixConn)
	go func() {
		for {
			conn, err := ln.AcceptUnix()
			if err != nil {
				log.Printf("error listening on new connection %s", err)
				continue
			}
			connChan <- conn
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	for {
		select {
		case <-sigChan:
			err := ln.Close()
			if err != nil {
				return err
			}
			f, err := ln.File()
			if err != nil {
				return err
			}
			return f.Close()
		case conn := <-connChan:
			go srv.handleConn(conn)
		}
	}
}

func (srv *Server) handleConn(conn rpc.UnixConnReaderWriterCloser) {
	defer conn.Close()
	req := rpc.FOpenRequest{}
	err := rpc.ReadRequest(conn, &req)
	if err != nil {
		log.Printf("error reading request from connection: %s", err)
		return
	}
	resp := rpc.FOpenResponse{}
	srv.OpenFile(req, &resp)
	err = rpc.WriteResponse(conn, resp)
	if err != nil {
		log.Printf("error writing request to connection: %s", err)
		return
	}
}
