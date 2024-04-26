package models

type Stats struct {
	TotalPastes int   `json:"totalPastes"`
	TotalClicks int	  `json:"totalClicks"`
	AvgClicks   int   `json:"avgClicks"`
}