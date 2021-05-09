package responses

const (
	SuccessResponseStatus = "ok"
	ErrorResponseStatus   = "error"
)

type ServerResponse struct {
	Status  string      `json:"status"`
	Content interface{} `json:"content"`
}
