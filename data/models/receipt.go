package models

type Receipt struct {
	ID           string
	UUID         string `json:"uuid"`
	Details      string `json:"details"`
	IsProcessing bool   `json:"isProcessing"`
	HasError     bool   `json:"hasError"`
	Error        string `json:"error"`
}
