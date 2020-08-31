package binchunk

//观察本机参数后确认
const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

//常量表 
const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type binaryChunk struct {
	header
	sizeUpvalues byte // ?
	mainFunc     *Prototype
}

type header struct {
	signature       [4]byte     //签名，ESC、L、u、a的ASCII码,快速识别文件格式
	version         byte        //major*16 + minor > hex
	format          byte
	luacData        [6]byte     //double check
	cintSize        byte        //接下来校验5种数据类型在chunk中占用的字节数
	sizetSize       byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luacInt         int64       //0x5678  检测大小端
	luacNum         float64     //f370.5  检测浮点数格式   例如IEEE 754
}

// function prototype
type Prototype struct {
	Source          string // debug
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte
	Code            []uint32
	Constants       []interface{}
	Upvalues        []Upvalue
	Protos          []*Prototype
	LineInfo        []uint32 // debug
	LocVars         []LocVar // debug
	UpvalueNames    []string // debug
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

func Undump(data []byte) *Prototype {
	rdr := &reader{data}
	rdr.checkHeader()
	rdr.readByte() // size_upvalues
	return rdr.readProto("")
}
