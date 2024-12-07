package errors

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

var (
	ErrUserNotFound       = APIError{404, "user not found"}
	ErrUserExists         = APIError{409, "user already exists"}
	ErrInvalidInput       = APIError{400, "invalid input"}
	ErrInternalServer     = APIError{500, "internal server error"}
	ErrInvalidCredentials = APIError{401, "invalid email or password"}
)
