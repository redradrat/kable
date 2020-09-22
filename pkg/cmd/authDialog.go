package cmd

import (
	"github.com/AlecAivazis/survey/v2"
)

func RunAuthDialog() (string, string, error) {
	username := ""
	inPrompt := &survey.Input{
		Message: "Please type your username",
	}
	err := survey.AskOne(inPrompt, &username)
	if err != nil {
		return "", "", err
	}

	password := ""
	pwPrompt := &survey.Password{
		Message: "Please type your password",
	}
	err = survey.AskOne(pwPrompt, &password)
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}
