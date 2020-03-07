package server

import (
	"testing"
)

func TestNew(t *testing.T) {
	srv := &Server{
		Addr: "/tmp/listener.sock",
	}
	err := srv.ListenAndServe()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(srv.listener)
}
