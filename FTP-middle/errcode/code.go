package errcode

const (
	ErrGeneral = (iota + 1) * 100
	ErrUsr
	ErrFtp
)

const (
	ErrGeneralStart = ErrGeneral + iota
	ErrInvalidData
	ErrInternalError
)
