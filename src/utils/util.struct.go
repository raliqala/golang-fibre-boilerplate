package utils

type SignIn struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}
