package models

type UpdatePaste struct {
	ShortID string                  `json:"shortId"`
	Content string             `json:"content"`
}