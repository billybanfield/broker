package rpc

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"unsafe"
)

type FOpenRequest struct {
	Filename string
}

type FOpenResponse struct {
	FileDescriptor uintptr
}

type FOpener interface {
	OpenFile(FOpenRequest, *FOpenResponse)
}

type UnixConnReader interface {
	Read(b []byte) (int, error)
	ReadMsgUnix(b, oob []byte) (n, oobn, flags int, addr *net.UnixAddr, err error)
}

type UnixConnWriter interface {
	Write(b []byte) (int, error)
	WriteMsgUnix(b, oob []byte, addr *net.UnixAddr) (n, oobn int, err error)
}

type UnixConnReaderWriterCloser interface {
	UnixConnReader
	UnixConnWriter
	Close() error
}

func ReadRequest(conn io.Reader, req *FOpenRequest) error {
	msgLen := int64(0)
	err := binary.Read(conn, binary.BigEndian, &msgLen)
	if err != nil {
		return err
	}
	msgBuffer := make([]byte, msgLen)
	_, err = conn.Read(msgBuffer)
	if err != nil {
		return err
	}
	req.Filename = string(msgBuffer)
	return nil
}

func ReadResponse(conn UnixConnReader, resp *FOpenResponse) error {
	datalen := unsafe.Sizeof(uintptr(0))
	oob := make([]byte, syscall.CmsgSpace(int(datalen)))
	_, _, _, _, err := conn.ReadMsgUnix(nil, oob)
	if err != nil {
		return err
	}
	sockControlMsgs, err := syscall.ParseSocketControlMessage(oob)
	if err != nil {
		return err
	}
	if len(sockControlMsgs) != 1 {
		return fmt.Errorf("incorrect number of control messages returned. expected %d, received %d", 1, len(sockControlMsgs))
	}
	fds, err := syscall.ParseUnixRights(&sockControlMsgs[0])
	if err != nil {
		return err
	}
	if len(fds) != 1 {
		return fmt.Errorf("incorrect number of file descrpitors returned. expected %d, received %d", 1, len(fds))
	}

	f := os.NewFile(uintptr(fds[0]), "")
	resp.FileDescriptor = f.Fd()
	return nil
}

func WriteRequest(conn UnixConnWriter, req FOpenRequest) error {
	nameBytes := []byte(req.Filename)
	msgLen := int64(len(nameBytes))

	err := binary.Write(conn, binary.BigEndian, msgLen)
	if err != nil {
		return err
	}
	n, _, err := conn.WriteMsgUnix(nameBytes, nil, nil)
	if err != nil {
		return err
	}
	if n != len(nameBytes) {
		return fmt.Errorf("incorrect number of bytes written to connection. wrote %d, expected %d", n, len(nameBytes))
	}
	return nil
}

func WriteResponse(conn UnixConnWriter, resp FOpenResponse) error {
	rights := syscall.UnixRights(int(resp.FileDescriptor))
	_, oobn, err := conn.WriteMsgUnix(nil, rights, nil)
	if err != nil {
		return err
	}
	if oobn != len(rights) {
		return fmt.Errorf("incorrect number of out of band bytes written to connection. wrote %d, expected %d", oobn, len(rights))
	}
	return nil
}
