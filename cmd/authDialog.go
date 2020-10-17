package cmd

import (
	"github.com/AlecAivazis/survey/v2"
)

// RunAuthDialog runs a dialog to ask for additionally needed auth info. Returns a bool whether auth is required.
func RunAuthDialog(user, pass string) (bool, string, string, error) {

	if user != "" && pass != "" {
		return false, user, pass, nil
	}

	if user == "" && pass == "" {
		do := true
		doPrompt := &survey.Confirm{
			Message: "Does this repository require basic authentication?",
		}
		err := survey.AskOne(doPrompt, &do)
		if err != nil {
			return false, "", "", err
		}
	}

	if user == "" {
		inPrompt := &survey.Input{
			Message: "Please type your username",
		}
		err := survey.AskOne(inPrompt, &user)
		if err != nil {
			return false, "", "", err
		}
	}

	if pass == "" {
		pwPrompt := &survey.Password{
			Message: "Please type your password",
		}
		err := survey.AskOne(pwPrompt, &pass)
		if err != nil {
			return false, "", "", err
		}
	}

	return true, user, pass, nil
}
