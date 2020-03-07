package rpc

import (
	"bytes"
	"net"
	"testing"
)

type mockUnixReader struct {
	buf    *bytes.Buffer
	oobBuf *bytes.Buffer
}

type mockUnixWriter struct {
	buf    *bytes.Buffer
	oobBuf *bytes.Buffer
}

func (m *mockUnixReader) Read(b []byte) (int, error) {
	return m.buf.Read(b)
}

func (m *mockUnixReader) ReadMsgUnix(b, oob []byte) (n, oobn, flags int, addr *net.UnixAddr, err error) {
	n, err = m.buf.Read(b)
	if err != nil {
		return
	}
	oobn, err = m.oobBuf.Read(oob)
	if err != nil {
		return
	}
	return
}

func (m *mockUnixWriter) Write(b []byte) (int, error) {
	return m.buf.Write(b)
}

func (m *mockUnixWriter) WriteMsgUnix(b, oob []byte, addr *net.UnixAddr) (n, oobn int, err error) {
	n, err = m.buf.Write(b)
	if err != nil {
		return
	}
	oobn, err = m.oobBuf.Write(oob)
	if err != nil {
		return
	}
	return
}

func TestReadWriteRequest(t *testing.T) {
	fname := "name"
	b := &bytes.Buffer{}
	oob := &bytes.Buffer{}

	r := &mockUnixReader{
		buf:    b,
		oobBuf: oob,
	}
	w := &mockUnixWriter{
		buf:    b,
		oobBuf: oob,
	}

	err := WriteRequest(w, FOpenRequest{fname})
	if err != nil {
		t.Fatal(err)
	}
	req := &FOpenRequest{}
	err = ReadRequest(r, req)
	if err != nil {
		t.Fatal(err)
	}
	if req.Filename != fname {
		t.Fatalf("name did not match expected")
	}
}

func TestReadWriteResponse(t *testing.T) {
	fd := uintptr(100)

	b := &bytes.Buffer{}
	oob := &bytes.Buffer{}

	r := &mockUnixReader{
		buf:    b,
		oobBuf: oob,
	}
	w := &mockUnixWriter{
		buf:    b,
		oobBuf: oob,
	}

	err := WriteResponse(w, FOpenResponse{fd})
	if err != nil {
		t.Fatalf("write response err: %s", err)
	}
	resp := &FOpenResponse{}
	err = ReadResponse(r, resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.FileDescriptor != fd {
		t.Fatalf("file descriptor did not match expected")
	}
}
