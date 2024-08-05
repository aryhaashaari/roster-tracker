package providererror

type Error struct {
	Code    int    `json:"-"`
	Errors  any    `json:"errors,omitempty"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
