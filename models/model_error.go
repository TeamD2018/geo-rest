package models

type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string {
	return e.Message
}

// List of Error
var (
	CourierAlreadyExists              = Error{"Courier already exists", 10}
	ServerError                       = Error{"Sorry, server error", 20}
	OneOfParametersNotFound           = Error{"One of parameters not found", 30}
	OneOfParameterHaveIncorrectFormat = Error{"One of parameter have incorrect format", 40}
	EntityNotFound                    = Error{"Entity with such id not found", 50}
)
