package model

type WebResponse struct {
	Data   interface{}         `json:"data,omitempty"`
	Paging *PaginationResponse `json:"paging,omitempty"`
	Error  string              `json:"errors,omitempty"`
}

type PaginationResponse struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}
