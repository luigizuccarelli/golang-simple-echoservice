package schema

// ShcemaInterface - acts as an interface wrapper for our profile schema
// All the go microservices will using this schema
type SchemaInterface struct {
	ID         string `json:"_id,omitempty"`
	LastUpdate int64  `json:"lastupdate,omitempty"`
	MetaInfo   string `json:"metainfo,omitempty"`
}

// Response schema
type Response struct {
	Code       int             `json:"code,omitempty"`
	StatusCode string          `json:"statuscode"`
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Payload    SchemaInterface `json:"payload"`
}
