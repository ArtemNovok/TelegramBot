package telegram

type UpdateResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result,omitempty"`
}

type Update struct {
	ID      int    `json:"update_id"`
	Message string `json:"message"`
}
