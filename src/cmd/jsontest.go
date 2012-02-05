package main

import (
    "fmt"
    "encoding/json"
)


var json1 = []byte(`{  "id":"ejiaden",
                "firstName":"James",
                "lastName": "Deng",
                "family":[{"name":"Chun"}, {"name":"Yang"}],
                "contact":{"tel":"+24049452","mobile":"09018334499"}
}`)

//JSONObj = JSONQuery(json1)
//JSONObj.get("id") 
//JSONObj.get("family") //returns a list a map
//JSONObj.get("family[0]") //returns first map in list
//JSONObj.get("family[0].name") //returns name of first map 
//JSONObj.get("contact.mobile") 
//
////if exists, update, otherwise add
//JSONObj.set(json_query, value)
//JSONObj.remove(json_query) //returns the removed 
//

func main() {
    var f interface{}
    err := json.Unmarshal(json1, &f)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        m := f.(map[string]interface{})
        fmt.Println(m["id"])
        fmt.Println(m["family"])
        fmt.Println(m["family"].([]interface{})[0])
        fmt.Println(m["family"].([]interface{})[0].(map[string]interface{})["name"])
        fmt.Println(m["contact"].(map[string]interface{})["mobile"])
        
    }
}