package client

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	cli, err := New("/tmp/uds.sock")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cli)
}
