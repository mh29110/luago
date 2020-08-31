package api

type LuaVM interface {
    LuaState
    PC() int            //return current Program Counter; (only for test)
    AddPC (n int)       //modify PC; (jump)
    Fetch() uint32      //取出当前指令；将PC指向下一条指令。
    GetConst (idx int)  //将指定常量推入栈顶
    GetRK(rk int)       //将指定常量或栈值推入栈顶
}
