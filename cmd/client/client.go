package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/billybanfield/broker/pkg/client"
	"github.com/billybanfield/broker/pkg/rpc"
)

var (
	fname    = flag.String("fname", "/etc/hostname", "")
	sockAddr = flag.String("sockAddr", "/tmp/test.sock", "")
)

func main() {
	flag.Parse()
	cli, err := client.New(*sockAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	resp := rpc.FOpenResponse{}
	err = cli.OpenFile(rpc.FOpenRequest{Filename: *fname}, &resp)
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}
	f := os.NewFile(resp.FileDescriptor, *fname)
	defer f.Close()
	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		log.Fatal(err)
	}
}
