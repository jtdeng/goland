// Copyright 2011 James Deng.  All rights reserved.
// This is a demo of go channels and network channels
package chat

import (
    //"fmt"
    "log"
    "netchan"
    "strconv"
)


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
    
    req := make(chan []byte)
    err = imp.Import("t-req", req, netchan.Send, 1)
    if err != nil { log.Fatal(err) }
        
    res := make(chan int64)
    err = imp.Import("t-res", res, netchan.Recv, 1)
    if err != nil { log.Fatal(err) }
    
    req <- []byte(name)
    ticket := <-res
    
    log.Println("Got ticket:", ticket)
    //imp.Hangup("t-req")
    //imp.Hangup("t-res")
    //close(req) //ticket has been closed by Hangup?
    //close(res)
    
    in, out := make(chan []byte), make(chan []byte)
    imp.Import("recv-"+strconv.Itoa64(ticket), out, netchan.Send, 1)
    imp.Import("send-"+strconv.Itoa64(ticket), in, netchan.Recv, 1)
    return in, out
}
