// Copyright 2011 James Deng.  All rights reserved.
// This is a demo of go channels and network channels
package chat

import (
    //"fmt"
    "log"
    "netchan"
)


type Client struct {
    name    string  //name of the client
    input   chan<-  string //channel for client to send msg
    output  <-chan  string //channel for client to recv msg
    
}


//NewClient connects to the server and returns 2 channels for communication 
//Parameters
//  name: name of the client
//  server: server address, e.g. "localhost:9898"
//Returns
//  input: channel for sending msgs, 
//  output: channel for receiving msgs
func NewClient(name, server string) (input chan<- []byte, output <-chan []byte) {
    
    imp, err := netchan.Import("tcp", server)
    if err != nil { log.Fatal(err) }
    
    out := make(chan []byte)
    err = imp.Import("broadcaster", out, netchan.Recv, 1)
    if err != nil { log.Fatal(err) }
        
    in := make(chan []byte)
    err = imp.Import("receiver", in, netchan.Send, 1)
    if err != nil { log.Fatal(err) }
    
    return in, out
}
