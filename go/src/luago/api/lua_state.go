package api

type LuaType = int

type LuaState interface {
    GetTop() int
    Type(idx int) LuaType
}
