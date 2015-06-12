package lib

type SysError struct {
	ErrorCode    int
	ErrorMessage string
	Err          error
}

func (e *SysError) Error() string {
	return string(e.ErrorCode) + " " + e.ErrorMessage + ": " + e.Err.Error()
}
