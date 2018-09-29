package openapi

type Error string

// List of Error
const (
	COURIER_ALREADY_EXISTS Error = "Courier already exists"
	SORRY_SERVER_ERROR Error = "Sorry, server error"
	ONE_OF_PARAMETERS_NOT_FOUND Error = "One of parameters not found"
	ONE_OF_PARAMETER_HAVE_INCORRECT_FORMAT Error = "One of parameter have incorrect format"
)