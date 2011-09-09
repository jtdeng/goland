// Copyright 2011 James Deng.  All rights reserved.
// This is a demo of go channels and network channels
// 
// chat server can host multiple clients by assigning a unique ticket
// to each client and exporting communication channels base on this ticket. 
// client ask for ticket via ticket channel first, then communicate with
// server on channels base on ticket name.

package chat

import (
    "log"
    "net"
    "netchan"
    "fmt"
)

type Server struct {
    name        string
    broadcaster chan []byte
    receiver    chan []byte
    exp         *netchan.Exporter
    clients		map[string]net.Conn
}

func (s *Server) Serve (addr string) {
    //goroutine for tcp serve
    listener, err := net.Listen("tcp", addr)
    //defer netlisten.Close()
    if err != nil {log.Fatal(err)}
    
    conn, err := listener.Accept()
    s.clients[conn.RemoteAddr().String()] = conn
    go s.exp.ServeConn(conn)
    
   
    for {
	    select {
	        case  msg := <-s.receiver:
	            fmt.Print(msg)
	            s.broadcaster <- msg
	            //s.exp.Drain(1e9)
	    } 
    }    

    
    
}

func NewServer(name string) (*Server) {
    
    exp := netchan.NewExporter()
    
    in := make(chan []byte)
    err := exp.Export("broadcaster", in, netchan.Send)
    if err != nil { log.Fatal(err) }
        
    out := make(chan []byte)
    err = exp.Export("receiver", out, netchan.Recv)
    if err != nil { log.Fatal(err) }
    
    s := &Server{name, in, out, exp}
    
    return s
}
