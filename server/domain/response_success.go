package domain

type ResponseSuccess struct {
	Success bool `json:"success"`
	Message string `json:"message,omitempty"`
	Data    interface{} `json:"data"`
	RequestId string `json:"reqId"`
}


