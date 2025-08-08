package auth

type UserAuthRequest struct {
	Phone     string `json:"phone" validate:"omitempty,numeric,len=11"`
	SessionID string `json:"session_id"`
	Code      string `json:"code"`
}

type UserAuthResponse struct {
	SessionID string `json:"session_id"`
}

type Token struct {
	Token string `json:"token"`
}
