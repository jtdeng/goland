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
    //"sync"
    "netchan"
    "strconv"
)

//the client representation is server side
type client struct {
	name		string
	recvChan	chan []byte	//chann to recv msg from client
	sendChan 	chan []byte //chan to send msg to client
}

type Server struct {
    //name        string
    //broadcaster chan []byte
    //receiver    chan []byte
    tickets		int64
    //ticketMutex	sync.Mutex	
    exp         *netchan.Exporter
    ticketReq	chan []byte
    ticketRes	chan int64
    clients		map[int64]*client
}

//handle ticket requests and assigning tickets to client
//then export the channel base on this ticket
func (s *Server) handleRequest() {
	for {
	    select { //no default, blocking
	        case  req := <-s.ticketReq:
	        	log.Println("Got ticket request from:", string(req))
	            //assign a ticket
	            s.tickets++
	            //export channels, create a new client
	            rChan, sChan := make(chan []byte), make(chan []byte)
	            log.Println("Exporting channel:", "recv-"+strconv.Itoa64(s.tickets))
	            s.exp.Export("recv-"+strconv.Itoa64(s.tickets), rChan, netchan.Recv)
	            log.Println("Exporting channel:", "send-"+strconv.Itoa64(s.tickets))
	            s.exp.Export("send-"+strconv.Itoa64(s.tickets), sChan, netchan.Send)
	            s.clients[s.tickets] = &client{string(req), rChan, sChan}
	    		//send ticket to client
	    		s.ticketRes <- s.tickets
	            s.exp.Drain(-1)
	            log.Println("Sent ticket:", s.tickets)
	            //ready to serve client
	            log.Println("Serving channel:", s.tickets)
	    		go s.handleSession(s.tickets)
	    } 
    }    
}

func (s *Server)handleSession(ticket int64) {
	for {
		select {
			case msg := <-s.clients[ticket].recvChan:
				log.Println(s.clients[ticket].name, msg)
				for k,v := range(s.clients) {
					if k!=ticket {
						log.Println("Forwarding message to", v.name)
						//v.sendChan <- msg
					}
				}
		}
	
	}	
}

func (s *Server) Serve (addr string) {
    //goroutine for tcp serve
    //listener, err := net.Listen("tcp", addr)
    //defer netlisten.Close()
    //if err != nil {log.Fatal(err)}
    //conn, err := listener.Accept()
    //s.clients[conn.RemoteAddr().String()] = conn
    //go s.exp.ServeConn(conn)
    
    go s.exp.ListenAndServe("tcp", addr)
    
    s.exp.Export("t-req", s.ticketReq, netchan.Recv)
    s.exp.Export("t-res", s.ticketRes, netchan.Send)
    s.handleRequest()

}

func NewServer() (*Server) {
    
    exp := netchan.NewExporter()
    
    //in := make(chan []byte)
    //err := exp.Export("broadcaster", in, netchan.Send)
    //if err != nil { log.Fatal(err) }
        
    //out := make(chan []byte)
    //err = exp.Export("receiver", out, netchan.Recv)
    //if err != nil { log.Fatal(err) }
    reqChan := make(chan []byte)
    resChan := make(chan int64)
    s := &Server{	tickets:0, 
    				exp:exp, 
    				ticketReq:reqChan, 
    				ticketRes:resChan,
    				clients: make(map[int64]*client)}
    return s
}
