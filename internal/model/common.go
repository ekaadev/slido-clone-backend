package model

type WebResponse struct {
	Data   interface{}         `json:"data,omitempty"`
	Paging *PaginationResponse `json:"paging,omitempty"`
	Error  string              `json:"errors,omitempty"`
}

type PaginationResponse struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
