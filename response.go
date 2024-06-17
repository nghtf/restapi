package restapi

type TResponse struct{}

type TResponseTemplate struct {
	Status string      `json:"status"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

const (
	StatusOK    = "Ok"
	StatusError = "Error"
)

func (r TResponse) OK() TResponseTemplate {
	return TResponseTemplate{
		Status: StatusOK,
	}
}

func (r TResponse) Data(v interface{}) TResponseTemplate {
	return TResponseTemplate{
		Status: StatusOK,
		Data:   v,
	}
}

func (r TResponse) Error(msg string) TResponseTemplate {
	return TResponseTemplate{
		Status: StatusError,
		Error:  msg,
	}
}
