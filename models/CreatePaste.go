package models

type CreatePaste struct {
	Content string           `json:"content"`
	Expiry  string           `json:"expiry"`
}