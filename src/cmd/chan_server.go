package main

import (
    "chat"
)

func main() {
    s := chat.NewServer()
    s.Serve("localhost:9988")
}