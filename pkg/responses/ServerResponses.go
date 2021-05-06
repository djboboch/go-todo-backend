package responses

type ServerResponse struct {
	Status  string      `json:"status"`
	Content interface{} `json:"content"`
}

const (
	SuccessResponseStatus = "ok"
	ErrorResponseStatus   = "error"
)
