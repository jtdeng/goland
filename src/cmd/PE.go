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
type ResourceHanlder interface {
	Create() //map to 'PUT'
	Set()    //map to 'POST'
	Get()    //map to 'GET'
	Delete() //map to 'DELETE'
}

//url router, routes url to registered resource handlers
type router struct {
	p       string
	cp      *regexp.Regexp
	handler ResourceHanlder
}

//add a url router to routers list
func (pe *ProvisioningEngine) AddRouter(pattern string, handler ResourceHanlder) {
	cp, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Invalid URL pattern: %q\n", pattern)
		return
	}

	//pe.Routers.Push(router{pattern, cp, handler})
	//pe.Routers[len(pe.Routers)] = router{pattern, cp, handler}
    pe.Routers = append(pe.Routers, router{pattern, cp, handler})
}

func (pe *ProvisioningEngine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	//TODO: set content type according to resource format json or xml
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")

	//URL should be like /Feature/ResourceName/ResourceId&format=[xml|json]
	//feature := req.URL.Path.Split()
	fmt.Println(req.URL.Path)

	for i := 0; i < len(pe.Routers); i++ {
        fmt.Println(pe.Routers[i].p)
	}
	
	_handler :=  pe.GetHandler()
	
	switch _method := req.Method; _method {
	case "GET": 
        //call the Get() method of the handler object	
	case "PUT":
	
	case "DELETE":
	
	case "POST":
	
	default:
	   
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
}

func (hlr *HLRSubscriberHandler) Delete() {
}

func (hlr *HLRSubscriberHandler) Get() {
}

func (hlr *HLRSubscriberHandler) Set() {
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
