package payload

type EmailRequest struct {
	Email string `json:"email"`
}

type Email struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}
