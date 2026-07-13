package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
)

func lowercaseValidator(val interface{}) error {
	if err := survey.Required(val); err != nil {
		return err
	}
	str := val.(string)
	if !regexp.MustCompile(`^[a-z]+$`).MatchString(str) {
		return errors.New("must contain only lowercase letters (no spaces or special characters)")
	}
	return nil
}

func serviceNameValidator(val interface{}) error {
	if err := survey.Required(val); err != nil {
		return err
	}
	str := val.(string)
	if !regexp.MustCompile(`^[a-z]+(-[a-z]+)*$`).MatchString(str) {
		return errors.New("must be in kebab-case (e.g. markitos-it-service-users)")
	}
	return nil
}

func titleValidator(val any) error {
	str := val.(string)
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9 ]*$`).MatchString(str) {
		return errors.New("must start with a letter, and contain only letters, numbers, or spaces")
	}
	return nil
}

func NewServiceAction() {
	fmt.Println("\n✨ Welcome to the Interactive Service Wizard ✨")

	var qs = []*survey.Question{
		{
			Name:     "singular",
			Prompt:   &survey.Input{Message: "Entity Singular (e.g. user):"},
			Validate: lowercaseValidator,
		},
		{
			Name:     "plural",
			Prompt:   &survey.Input{Message: "Entity Plural (e.g. users):"},
			Validate: lowercaseValidator,
		},
		{
			Name:     "serviceName",
			Prompt:   &survey.Input{Message: "Service Name (e.g. markitos-it-service-users):"},
			Validate: serviceNameValidator,
		},
	}

	answers := struct {
		Singular    string
		Plural      string
		ServiceName string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println("❌ Wizard cancelled.")
		return
	}

	confirm := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Ready to create '%s'. Proceed?", answers.ServiceName),
	}
	survey.AskOne(prompt, &confirm)

	if !confirm {
		fmt.Println("❌ Aborted.")
		return
	}

	entitySingular = answers.Singular
	entityPlural = answers.Plural

	cloneCmd.Run(cloneCmd, []string{answers.ServiceName})
}
