//an interactive chat client using chat package
package main

import (
    "chat"
    "fmt"
    "bufio"
    "os"
    "bytes"
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


func main() {
    in, out := chat.NewClient("James", "localhost:9988")
    
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