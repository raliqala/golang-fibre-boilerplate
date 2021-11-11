package utils

type SignIn struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type EmailVerification struct {
	GreetUseStyle     string `json:"greet_use_style"`
	Username          string `json:"username"`
	VerifyLink        string `json:"verify_link"`
	SignatureGreeting string `json:"signature_greeting"`
	EmailSignature    string `json:"email_signature"`
}

type Success struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type RefreshTokens struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

type AccessTokens struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

type AuthTokensObject struct {
	Access  AccessTokens  `json:"access"`
	Refresh RefreshTokens `json:"refresh"`
}

type ForgotPass struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type ResetPass struct {
	Password string `json:"password"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}
