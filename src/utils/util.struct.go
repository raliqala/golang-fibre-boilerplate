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
