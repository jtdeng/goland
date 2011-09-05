package main

import (
	"rpc"
	"rpc/jsonrpc"
	"net"
	"fmt"
	"os"
	"http"
	"io/ioutil"
	"sync"
	//"time"
)

type Args struct {
	Timezone string
}

type Reply struct {
	Time    string
	Counter int64
}

type TimeService struct {
	//to protect counter
	//Google says:
	//Regarding mutexes, the sync package implements them, but we hope Go programming 
	//style will encourage people to try higher-level techniques. In particular, consider 
	//structuring your program so that only one goroutine at a time is ever responsible 
	//for a particular piece of data.
	//in short, Don't communicate by sharing memory; share memory by communicating. 
	locker  sync.Mutex 
	counter int64
}

func (ts *TimeService) GetTime(args *Args, reply *Reply) os.Error {
	ts.locker.Lock()
	defer ts.locker.Unlock()
	
	ts.counter++
	resp, _ := http.Get("http://json-time.appspot.com/time.json?tz=GMT")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	reply.Time = string(body)
	reply.Counter = ts.counter
	return nil
}

func init() {
	rpc.Register(new(TimeService))
}
//Start a RPC server that can return time
//Time is retrived from http://json-time.appspot.com/time.json?tz=GMT

func test_rpc() {
	fmt.Println("Testing rpc with TimeService...")
	cli, srv := net.Pipe()
	go rpc.ServeConn(srv)

	client := rpc.NewClient(cli)
	defer client.Close()

	// Synchronous calls
	args := &Args{"GMT"}
	reply := new(Reply)
	for i := 0; i < 10; i++ {
		err := client.Call("TimeService.GetTime", args, reply)
		if err != nil {
			fmt.Errorf("TimeService.GetTime: expected no error but got string %q", err.String())
		}

		fmt.Printf("time:%s\n rpc.counter:%d\n", reply.Time, reply.Counter)
	}
}

func test_jsonrpc() {
	fmt.Println("Testing jsonrpc with TimeService...")

	cli, srv := net.Pipe()
	go jsonrpc.ServeConn(srv)

	client := jsonrpc.NewClient(cli)
	defer client.Close()

	// Synchronous calls
	args := &Args{"GMT"}
	reply := new(Reply)
	for i := 0; i < 10; i++ {
		err := client.Call("TimeService.GetTime", args, reply)
		if err != nil {
			fmt.Errorf("TimeService.GetTime: expected no error but got string %q", err.String())
		}

		fmt.Printf("time:%s\n jsonrpc.counter:%d\n", reply.Time, reply.Counter)
	}
}

func main() {
    //test rpc in go routine, demo the use of sync.Mutex
	go test_rpc()
    test_jsonrpc()

	fmt.Println("One thing you can find is that rpc and jsonrpc share the same TimeService instance")
}
