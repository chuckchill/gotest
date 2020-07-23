package main

import (
	"flag"
	"fmt"
	"github.com/example"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"log"
	"net"
	"time"
)

var (
	addr       = flag.String("addr", "localhost:8972", "server address")
	consulAddr = flag.String("consulAddr", "10.8.0.107:8600", "consul address")
	basePath   = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	flag.Parse()
	s := server.NewServer()
	addRegistryPlugin(s)
	s.RegisterName("Arith", new(example.Arith), "")
	s.RegisterName("Test", new(example.Test), "")
	fmt.Println(*addr)
	if err := s.Serve("tcp", *addr); err != nil {
		fmt.Println(err)
	}
}
func addRegistryPlugin(s *server.Server) {
	
	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		ConsulServers:  []string{*consulAddr},
		BasePath:       *basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(&ConnectionPlugin{})
	s.Plugins.Add(r)
}

type ConnectionPlugin struct {
}

func (p *ConnectionPlugin) HandleConnAccept(conn net.Conn) (net.Conn, bool) {
	log.Printf("client %v connected", conn.RemoteAddr().String())
	return conn, true
}
func (p *ConnectionPlugin) HandleConnClose(conn net.Conn) bool {
	log.Printf("client %v closed", conn.RemoteAddr().String())
	return true
}
