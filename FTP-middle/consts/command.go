package consts

const (
	EOF = '\n'

	LOGIN = "LOGIN "
	EXIT  = "EXIT "

	ECHO     = "ECHO "
	LIST     = "LIST "
	DOWNLOAD = "DOWN "
	UPLOAD   = "UPLOAD "

	PAUSE  = "PAUSE "
	RESUME = "RESUME "
)

const (
	STOP = "STOP"
)

const (
	FTP_CTRL_PREFIX = "CTRL:"
	FTP_DATA_PREFIX = "DATA:"

	FTP_PAUSE  = FTP_CTRL_PREFIX + "pause"
	FTP_RESUME = FTP_CTRL_PREFIX + "resume"
	FTP_STOP   = FTP_CTRL_PREFIX + "stop"
)
