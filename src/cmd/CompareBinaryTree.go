package main

import "tour/tree"
import "fmt"

//type Tree struct {
//	Left  *Tree
//	Value int
//	Right *Tree
//}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
    var F func(*tree.Tree)
    
    F = func(t *tree.Tree) {
    	if t != nil {
        	F(t.Left)
		    ch <- t.Value
		    F(t.Right)
    	} 
    
    }
    // we have to close ch after all recursive all,
    // otherwise the receiver will wait it forever
    // so we put all recursive task in a closure F()
    F(t)
    close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
    c1 := make(chan int)
    c2 := make(chan int)
    
    go Walk(t1, c1)
    go Walk(t2, c2)

    
    for v1 := range c1 {
        v2 := <- c2
        if v1 == v2 {
            continue    
        } else {
            return false
        }
    }
    
    return true
}

func main() {
    t1 := tree.New(1)
    t2 := tree.New(2)
    //c1 := make(chan int)
    //go Walk(t1, c1)
    //for v1 := range c1 {
    //	fmt.Print(v1, ",")
    //}
    if r := Same(t1,t2); r {
        fmt.Println("Trees are same")
    } else {
        fmt.Println("Trees are different")
    }
}
