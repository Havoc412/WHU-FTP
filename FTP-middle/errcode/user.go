package errcode

const (
	ErrInvaildUser = iota + ErrUsr
	ErrPassword
	ErrSignIn
	ErrNotSignIn

	ErrCreateFile
	ErrSaveFile
)

func UserMsgInit(m msg) {
	m[ErrInvaildUser] = "无效用户名"
	m[ErrPassword] = "密码错误"
	m[ErrSignIn] = "已登录，无法执行此操作"
	m[ErrNotSignIn] = "未登录，无法执行此操作"

	m[ErrCreateFile] = "本地创建文件失败"
	m[ErrSaveFile] = "本地保存文件失败"
}
