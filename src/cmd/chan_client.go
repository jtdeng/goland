//an interactive chat client using chat package
package main

import (
    "chat"
    "fmt"
    "bufio"
    "os"
    "bytes"
    "flag"
)

func keyboard_listener(kb chan []byte) {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("you> ")
        input, _ := reader.ReadBytes('\n')
        if !bytes.Equal(input, []byte("\n")) {
            kb <- input
        }
    }
}


var client_name = flag.String("n", "Client", "Client name")

func main() {
	flag.Parse()
    in, out := chat.NewClient(*client_name, "localhost:9988")
    
    kb := make(chan []byte)
    go keyboard_listener(kb)
    
    for {
        select {
            case msg := <-out:
                fmt.Printf("BROADCAST: %s", msg)
            case msg := <-kb:
                in <- msg
        }
    }
}