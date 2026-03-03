package domain

type SucessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
