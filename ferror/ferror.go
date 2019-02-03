package ferror

type Code string

const (
	InternalServerError  Code = "InternalServerError"
	NotFound             Code = "NotFound"
	InvalidCredential    Code = "InvalidCredential"
	UsernameNotAvailable Code = "Username is not available"
)

type Error struct {
	HttpStatusCode int
	ErrorCode      Code
	Message        string
}

func (err *Error) Error() string {
	return err.Message
}

func New(status int, errorCode Code, message string) *Error {
	return &Error{
		HttpStatusCode: status,
		ErrorCode:      errorCode,
		Message:        message,
	}
}
