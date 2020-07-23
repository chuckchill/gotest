package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)
type Args struct{
	A, B int
}

type Quotient struct{
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error{
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error{
	if args.B == 0{
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main()  {
	arith := new(Arith)
	rpc.Register(arith)
	
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")//jsonrpc是基于TCP协议的，现在他还不支持http协议
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for{
		conn, err := listener.Accept()
		if err != nil{
			continue
		}
		fmt.Printf("%#v\n",conn.RemoteAddr().String())
		jsonrpc.ServeConn(conn)
	}
}