package errcode

const (
	ErrServer = iota + ErrFtp

	ErrDownloadFileLength
	ErrDownloadFileData

	ErrOpenLocalFile

	ErrGetRides
	ErrUploadInterrupt
	ErrGetNewLink

	ErrTaskPause
	ErrTaskResume
)

func FtpMsgInit(m msg) {
	m[ErrServer] = "FTP服务器不在连接中，请联系管理员"
	m[ErrDownloadFileLength] = "获取文件长度错误，请重试"
	m[ErrDownloadFileData] = "获取文件数据错误，请重试或重传"
	m[ErrOpenLocalFile] = "打开本地文件错误"

	m[ErrGetRides] = "任务已失效"
	m[ErrUploadInterrupt] = "任务异常中断，请重连"
	m[ErrGetNewLink] = "获取新连接失败"

	m[ErrTaskPause] = "FTP服务器任务暂停失败"
	m[ErrTaskResume] = "FTP服务器任务恢复失败"
}
