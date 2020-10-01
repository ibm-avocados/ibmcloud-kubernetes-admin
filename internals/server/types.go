package server

type AccountEmailBody struct {
	AccountID string   `json:"accountID"`
	Email     []string `json:"email"`
}

type StatusOK struct {
	Message string `json:"message"`
}

type AccountLogin struct {
	OTP string `json:"otp"`
}

type Bill struct {
	Bill string `json:"bill"`
}
