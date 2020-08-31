package main

import "fmt"
import "io/ioutil"
import "os"
import "luago/binchunk"

func main(){
    fmt.Println("----- start -----");
    if len(os.Args) > 1 {
        data,err := ioutil.ReadFile(os.Args[1])
        if err != nil {panic(err)}

        proto := binchunk.Undump(data)
        fmt.Println(proto)
        list(proto)
    }
}

func list(f * binchunk.Prototype) {
    printHeader(f)
    //printCode(f)
    //printDetail(f)
    for _,p := range f.Protos {
        list(p)
    }
}

func printHeader(f *binchunk.Prototype){
    funcType := "main"
    if f.LineDefined > 0 { funcType = "function" }

    varargFlag := ""

    if f.IsVararg > 0 { varargFlag = "+" }

    fmt.Printf ( " \n %s <%s:%d,%d> ( %d instructions )\n" ,
                funcType,f.Source,
                f.LineDefined,f.LastLineDefined,len(f.Code))
    fmt.Printf( "%d %s params , %d slots , %d upvalues,",
                f.NumParams,varargFlag,f.MaxStackSize,len(f.Upvalues))
    fmt.Printf( "\n%d locals , %d constants , %d functions \n",
                len(f.LocVars),len(f.Constants),len(f.Protos))
}
