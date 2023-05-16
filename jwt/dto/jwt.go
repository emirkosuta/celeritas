package dto

type Token struct {
	Type         string       `json:"type"`
	Token        string       `json:"token"`
	RefreshToken RefreshToken `json:"-"`
}

type RefreshToken struct {
	Value string `json:"refresh_token"`
	Iat   string `json:"iat"`
}
