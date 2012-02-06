package main

import "Greeting"

//this is to demo how to use a package in go
func main() {
	Greeting.SayHello(getName())
}

