package responses

type Response struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Result  map[string]interface{} `json:"result"`
}

type DeleteResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type UpdateResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
