package main

import (
        "fmt"
        "os"
        "netchan"
)

func main() {
        importer, err := netchan.Import("tcp", "localhost:1234")
        if err != nil {
                fmt.Printf("Error calling Import(): %v\n", err)
                return
        }

        sender := make(chan string)
        receiver := make(chan string)
        importer.Import("input", sender, netchan.Send, 1)
        importer.Import("output", receiver, netchan.Recv, 1)

        go func(errs chan os.Error){
                for {
                        fmt.Printf("Got error in send: %v\n", <-errs)
                }
        }(importer.Errors())

        sender <-"foo"
        fmt.Printf("Got %q\n", <-receiver)
}
