package binchunk

import "encoding/binary"
import "math"


type reader struct {
    data []byte
}

func (rd *reader) readByte() byte{
    b := rd.data[0]
    rd.data = rd.data[1:]
    return b
}

func (rd *reader) readBytes(n uint) []byte {
	bytes := rd.data[:n]
	rd.data = rd.data[n:]
	return bytes
}

func (rd *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(rd.data)
	rd.data = rd.data[4:]
	return i
}

func (rd *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(rd.data)
	rd.data = rd.data[8:]
	return i
}

func (rd *reader) readLuaInteger() int64 {
	return int64(rd.readUint64())
}

func (rd *reader) readLuaNumber() float64 {
	return math.Float64frombits(rd.readUint64())
}

func (rd *reader) readString() string {
	size := uint(rd.readByte())
	if size == 0 {
		return ""
	}
	if size == 0xFF {
		size = uint(rd.readUint64()) // size_t
	}
	bytes := rd.readBytes(size - 1)
	return string(bytes) // todo
}

func (rd *reader) checkHeader() {
	if string(rd.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	}
	if rd.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	}
	if rd.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	}
	if string(rd.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}
	if rd.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	}
	if rd.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch!")
	}
	if rd.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}
	if rd.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	}
	if rd.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}
	if rd.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	}
	if rd.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (rd *reader) readProto(parentSource string) *Prototype {
	source := rd.readString()
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     rd.readUint32(),
		LastLineDefined: rd.readUint32(),
		NumParams:       rd.readByte(),
		IsVararg:        rd.readByte(),
		MaxStackSize:    rd.readByte(),
		Code:            rd.readCode(),
		Constants:       rd.readConstants(),
		Upvalues:        rd.readUpvalues(),
		Protos:          rd.readProtos(source),
		LineInfo:        rd.readLineInfo(),
		LocVars:         rd.readLocVars(),
		UpvalueNames:    rd.readUpvalueNames(),
	}
}

func (rd *reader) readCode() []uint32 {
	code := make([]uint32, rd.readUint32())
	for i := range code {
		code[i] = rd.readUint32()
	}
	return code
}

func (rd *reader) readConstants() []interface{} {
	constants := make([]interface{}, rd.readUint32())
	for i := range constants {
		constants[i] = rd.readConstant()
	}
	return constants
}

func (rd *reader) readConstant() interface{} {
	switch rd.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return rd.readByte() != 0
	case TAG_INTEGER:
		return rd.readLuaInteger()
	case TAG_NUMBER:
		return rd.readLuaNumber()
	case TAG_SHORT_STR, TAG_LONG_STR:
		return rd.readString()
	default:
		panic("corrupted!") // todo
	}
}

func (rd *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, rd.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: rd.readByte(),
			Idx:     rd.readByte(),
		}
	}
	return upvalues
}

func (rd *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, rd.readUint32())
	for i := range protos {
		protos[i] = rd.readProto(parentSource)
	}
	return protos
}

func (rd *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, rd.readUint32())
	for i := range lineInfo {
		lineInfo[i] = rd.readUint32()
	}
	return lineInfo
}

func (rd *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, rd.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: rd.readString(),
			StartPC: rd.readUint32(),
			EndPC:   rd.readUint32(),
		}
	}
	return locVars
}

func (rd *reader) readUpvalueNames() []string {
	names := make([]string, rd.readUint32())
	for i := range names {
		names[i] = rd.readString()
	}
	return names
}
