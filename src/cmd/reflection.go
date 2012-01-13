package main

import ("reflect"
        "fmt")

type MyInt int


func main() {
    var i MyInt = 3
    v := reflect.ValueOf(i)
    fmt.Println("Type:", v.Type(), "Kind:", v.Kind())    

}