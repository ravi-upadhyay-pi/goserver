package ferror

type Code string

const (
	InternalServerError   Code = "InternalServerError"
	NotFound              Code = "NotFound"
	Unauthorized          Code = "Unauthorized"
	InvalidCredential     Code = "InvalidCredential"
	UsernameNotAvailable  Code = "UsernameNotAvailable"
	ArticleIdNotAvailable Code = "ArticleIdNotAvailable"
)

type Error struct {
	HttpStatus int
	ErrorCode  Code
	Message    string
}

func (err *Error) Error() string {
	return err.Message
}

func New(status int, errorCode Code, message string) *Error {
	return &Error{
		HttpStatus: status,
		ErrorCode:  errorCode,
		Message:    message,
	}
}
