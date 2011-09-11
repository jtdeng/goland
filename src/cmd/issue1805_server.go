package main

import (
        "fmt"
        "netchan"
)

func runner(in, out chan string) {
        for {
                value := <-in
                fmt.Printf("Received %q\n", value)
                sentValue := "Someone sent " + value
                out <- sentValue
                fmt.Printf("Sent %q\n", sentValue)
        }
}

func main() {
        exporter := netchan.NewExporter()
        in := make(chan string)
        out := make(chan string)
        exporter.Export("input", in, netchan.Recv)
        exporter.Export("output", out, netchan.Send)
        exporter.ListenAndServe("tcp", "localhost:1234")
        runner(in, out)
}