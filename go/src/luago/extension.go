package main

import "fmt"
import "luago/binchunk"

func list(f * binchunk.Prototype) {
    printHeader(f)
    printCode(f)
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

func printCode( f * binchunk.Prototype){
    for pc,c := range f.Code {
        line := "-"
        if len(f.LineInfo) > 0 {
            line = fmt.Sprintf("%d",f.LineInfo[pc])
        }
        fmt.Printf("\t %d \t [%s] \t 0x%08X\n", pc+1 , line , c)
    }
}

