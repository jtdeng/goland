// Copyright 2012 James Deng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// an experimental RESTful/HTTP provisioning server.
//

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"regexp"
	"runtime"
	//"url"
)

type ProvisioningEngine struct {
	Routers []router
}

//every resource should be able to handle the 4 operations
type ResourceHandler interface {
	Create() //map to 'PUT'
	Set()    //map to 'POST'
	Get()    //map to 'GET'
	Delete() //map to 'DELETE'
}

//url router, routes url to registered resource handlers
type router struct {
	p       string
	cp      *regexp.Regexp
	handler ResourceHandler
}

//add a url router to routers list
func (pe *ProvisioningEngine) AddRouter(pattern string, handler ResourceHandler) {
	cp, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Invalid URL pattern: %q\n", pattern)
		return
	}
    pe.Routers = append(pe.Routers, router{pattern, cp, handler})
}

//lookup the handler with given url string
func (pe *ProvisioningEngine) LookupHandler(url string) ResourceHandler {
    for _, r := range pe.Routers {
        if r.cp.MatchString(url) {
            fmt.Println("Handler found for: ", r.p)
            return r.handler
        }
	}
	return nil 
}

func (pe *ProvisioningEngine) ServeHTTP(resp http.ResponseWriter, 
                                        req *http.Request) {
	//TODO: set content type according to resource format json or xml
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")

	//URL should be like /Feature/ResourceName/ResourceId&format=[xml|json]
	//feature := req.URL.Path.Split()
	fmt.Println(req.URL.Path)

	for _, r := range pe.Routers {
        fmt.Println(r.p)
	}
	
	//lookup the handler for the incoming request
	_handler :=  pe.LookupHandler(req.URL.Path)
	if _handler == nil {
	   http.NotFound(resp, req)
	   return
	}
	
	switch _method := req.Method; _method {
	case "GET": 
        _handler.Get()
	case "PUT":
	    _handler.Create()
	case "DELETE":
	    _handler.Delete()
	case "POST":
	    _handler.Set()
	default:
	    http.NotFound(resp, req)
	}
	
}

type Configuration map[string]string

func (pe *ProvisioningEngine) Initialize(cfg Configuration) {
	//read from config file and initialize PE
}

func (pe *ProvisioningEngine) Run(addr string) {
	pe.Initialize(Configuration{"Name": "test"})

	mux := http.NewServeMux()

	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/heap", http.HandlerFunc(pprof.Heap))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/", pe) //calls pe.ServeHTTP

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	err = http.Serve(l, mux)
	l.Close()
}

//////////////////////////////////////////////////////////////

type HLRSubscriberHandler struct {
}

func (hlr *HLRSubscriberHandler) Create() {
    fmt.Println("HLRSubscriberHandler created")
}

func (hlr *HLRSubscriberHandler) Delete() {
    fmt.Println("HLRSubscriberHandler deleted")
}

func (hlr *HLRSubscriberHandler) Get() {
    fmt.Println("HLRSubscriberHandler got")
}

func (hlr *HLRSubscriberHandler) Set() {
    fmt.Println("HLRSubscriberHandler set")
}

func init() {
	runtime.GOMAXPROCS(4)
	//TODO: load config from file
}

func main() {
	var pe = &ProvisioningEngine{}
	//register handlers for resources
	pe.AddRouter("/hlrsub", &HLRSubscriberHandler{})
	pe.Run(":8080")
}
