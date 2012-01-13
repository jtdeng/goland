package main

import "fmt"

type Printable interface {
	PrintIt()
}

type Inner struct {
	val int
}


func (i *Inner) PrintMe () {
	fmt.Println(i.val)
}

type Point struct {
	x int
	y int
	z *Inner
	p *Printable
}

func (p *Point) Area() int {
	//change p will change the receiver instance also, since p is a pointer to it
	area := p.x * p.y
	p.x = 100
	p.y = 100
	return area
}

func (p Point) Area2() int {
	// change p will not change the receiver instance, since p is a copy of it
	area := p.x * p.y
	p.x = 20
	p.y = 20
	return area
}

type Header map[string]string 

func (h Header) Append(key, val string) {
    h[key] = val, true
}


func main() {
	var pp = &Point{3,2,&Inner{7}, new(Printable)}
	fmt.Println(pp.Area())
	fmt.Println(pp)
	fmt.Println(pp.Area2())
	fmt.Println(pp)
	
	pp.z.PrintMe()
	(*pp.z).PrintMe()
	//you can not reference the method in Printable interface like below
	//pp.p.PrintIt()
	//seems that you can not use the pointer indirect (*p).method() ---->  p.method(), if method is from the interface
	
	//(*pp.p).PrintIt() // will fail because you are calling a nil instance


    var h = Header{}
    //or var h = make(Header)
    
    h.Append("James", "Deng")
    fmt.Println(h)
}