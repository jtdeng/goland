package main

import (
    "chat"
)

func main() {
    s := chat.NewServer("Chat Server1")
    s.Serve("localhost:9988")
}