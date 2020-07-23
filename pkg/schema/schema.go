package schema

type Request struct {
	Id      string `json:"id,omitemptye"`
	Message string `json:"message"`
}

// Response schema
type Response struct {
	Code    int    `json:"code,omitempty"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
