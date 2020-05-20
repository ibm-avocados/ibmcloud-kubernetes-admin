package server

type AccountEmailBody struct {
	AccountID string   `json:"accountID"`
	Email     []string `json:"email"`
}
