package mlflow

const (
	// ErrorUnauthorized indicates that the requester is unauthenticated or unauthorized to make the request.
	ErrorUnauthorized = "UNAUTHORIZED"

	// ErrorResourceAlreadyExists indicates that a resource with the given name already exists.
	ErrorResourceAlreadyExists = "RESOURCE_ALREADY_EXISTS"

	// ErrorResourceDoesNotExist indicates that the requested resource does not exist.
	ErrorResourceDoesNotExist = "RESOURCE_DOES_NOT_EXIST"
)

// Error represents an error returned by the MLflow API.
type Error struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}
