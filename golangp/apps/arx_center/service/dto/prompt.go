package dto

type GetPromptRequest struct {
	Pagination
	Keyword string `form:"keyword"`
}

type CreatePromptRequest struct {
	Content       string `json:"content"`
	Keywords      string `json:"keywords"`
	Settings      string `json:"settings"`
	ReferenceFile string `json:"reference_file"`
}

type UpdatePromptRequest struct {
	Base
	Content       string `json:"content"`
	Keywords      string `json:"keywords"`
	Settings      string `json:"settings"`
	ReferenceFile string `json:"reference_file"`
}
