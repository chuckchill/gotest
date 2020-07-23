package main

import (
	"context"
	"flag"
	"github.com/example"
	"github.com/smallnest/rpcx/client"
	"log"
)

var (
	addr       = flag.String("addr", "localhost:8972", "server address")
	consulAddr = flag.String("consulAddr", "10.8.0.107:8600", "consul address")
	basePath   = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	/*flag.Parse()
	kv, err := libkv.NewStore(store.CONSUL, []string{*consulAddr}, nil)
	if err != nil {
		panic(err)
	}
	ps, err := kv.List(*basePath + "/Arith")
	if err != nil {
		panic(err)
	}
	for _, p := range ps {
		fmt.Printf("%#s\n", p.LastIndex)
	}
	os.Exit(1)*/
	d := client.NewConsulDiscovery(*basePath, "Arith", []string{*consulAddr}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	args := &example.Args{
		A: 1,
		B: 2,
	}
	var reply int
	var call chan *client.Call
	called, err := xclient.Go(context.Background(), "Mul", args, &reply, call)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	<-called.Done
	log.Printf("%d * %d = %d", args.A, args.B, reply)
	
}
