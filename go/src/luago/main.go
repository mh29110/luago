package main

import "fmt"
//import "io/ioutil"
//import "os"

import . "luago/api"
import "luago/state"

func main(){
    fmt.Println("----- start -----");
    /*
    if len(os.Args) > 1 {
        data,err := ioutil.ReadFile(os.Args[1])
        if err != nil {panic(err)}

        proto := binchunk.Undump(data)
        fmt.Println(proto)
        list(proto)
    }
    */

    ls := state.New()
    ls.PushNil()
    ls.PushBoolean(true)
    ls.PushString("hello")
    printStack(ls)
}

func printStack(ls LuaState) int {
    top := ls.GetTop()

    for i:=1; i<=top; i++ {
        t := ls.Type(i)
        switch t {
            case LUA_TBOOLEAN:
                fmt.Printf("[%t]", ls.ToBoolean(i))
            case LUA_TNUMBER:
                fmt.Printf("[%g]", ls.ToNumber(i))
            case LUA_TSTRING:
                fmt.Printf("[%q]", ls.ToString(i))
            default: fmt.Printf("[%s]", ls.TypeName(t))
        }
    }
    return 0

}
