package errcode

type msg map[int]string

var ErrMsg msg // info 构建，用于全局

func init() {
	ErrMsg = make(msg)
	GeneralMsgInit(ErrMsg)
	UserMsgInit(ErrMsg)
	FtpMsgInit(ErrMsg)
}

func GeneralMsgInit(m msg) {
	m[0] = ""
	m[ErrInvalidData] = "参数无效"
	m[ErrInternalError] = "内部服务器错误"
}
