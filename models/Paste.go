package models

type Paste struct {
	ShortID     string      `json:"shortId"`
	ClickCount  int         `json:"clickCount"`
    Content     string 	  	`json:"content"`
	CreatedAt   string      `json:"createdAt"`
}