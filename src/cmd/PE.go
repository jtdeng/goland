package main

import (
	"fmt"
	"container/vector"
	"http"
	"http/pprof"
	"regexp"
	"net"
	"log"
)



type ProvisioningEngine struct {
	Routers	vector.Vector
	
}

//every resource should be able to handle the 4 operations
type ResourceHanlder interface {
	Create() //map to 'PUT'
	Set()	//map to 'POST'
	Get()	//map to 'GET'
	Delete()	//map to 'DELETE'
}

type router struct {
	p       string
    cp      *regexp.Regexp
    handler ResourceHanlder
}

//add a router to routers list
func (pe *ProvisioningEngine) AddRouter(pattern string, handler ResourceHanlder) {
	cp, err := regexp.Compile(pattern)
	if err != nil {
        fmt.Printf("Invalid URL pattern: %q\n", pattern)
        return
    }
    
   pe.Routers.Push(router{pattern, cp, handler})
}

func (pe *ProvisioningEngine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

}

type Configuration map[string]string
func (pe *ProvisioningEngine) Initialize(cfg Configuration){
	//read config file and initialize PE
}



func (pe *ProvisioningEngine) Run (addr string) {
	pe.Initialize(Configuration{"Name":"test"})

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
}

func (hlr *HLRSubscriberHandler) Delete() {
}

func (hlr *HLRSubscriberHandler) Get() {
}

func (hlr *HLRSubscriberHandler) Set() {
}

func main() {
	var pe = &ProvisioningEngine{}
	//register handlers for resources
	pe.AddRouter("/hlrsub", &HLRSubscriberHandler{})
	pe.Run(":8080")
}