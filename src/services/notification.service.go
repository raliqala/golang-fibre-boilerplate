package services

import (
	"log"

	"github.com/gobuffalo/plush"
	"github.com/raliqala/golang-fibre-boilerplate/src/config"
	"github.com/raliqala/golang-fibre-boilerplate/src/utils"
)

func EmailVerification(emailValues utils.EmailVerification) string {
	ctx := plush.NewContext()
	emailValues.EmailSignature = "SafePass"
	emailValues.SignatureGreeting = "regards"
	emailValues.GreetUseStyle = "Hi"
	ctx.Set("greetUser", emailValues.GreetUseStyle)
	ctx.Set("Username", emailValues.Username)
	ctx.Set("verifyLink", config.Config("APP_URL")+"/verify"+emailValues.VerifyLink)
	ctx.Set("greeting", emailValues.SignatureGreeting)
	ctx.Set("team", emailValues.EmailSignature)
	content, err := plush.Render(utils.LoadTemplates("email_validation"), ctx)
	if err != nil {
		log.Fatal(err)
	}

	return content
}
