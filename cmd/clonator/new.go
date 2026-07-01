package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"

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

func titleValidator(val interface{}) error {
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

	// Bucle para recolectar Custom Fields
	customFields = []Field{}
	for {
		addMore := false
		survey.AskOne(&survey.Confirm{
			Message: "Do you want to add a custom model field?",
			Default: false,
		}, &addMore)

		if !addMore {
			break
		}

		var f Field
		survey.Ask([]*survey.Question{
			{
				Name:     "Name",
				Prompt:   &survey.Input{Message: "Field Name (e.g. age):"},
				Validate: lowercaseValidator,
			},
			{
				Name:     "Title",
				Prompt:   &survey.Input{Message: "Field Title (e.g. Age):"},
				Validate: titleValidator,
			},
			{
				Name: "Type",
				Prompt: &survey.Select{
					Message: "Field Type:",
					Options: []string{"string", "int", "int64", "float64", "bool", "time.Time", "enum"},
				},
			},
		}, &f)

		if f.Type == "enum" {
			survey.AskOne(&survey.Input{Message: "Enum Values (comma separated, e.g. admin,user):"}, &f.EnumValues, survey.WithValidator(survey.Required))
		}

		survey.Ask([]*survey.Question{
			{
				Name:   "DefaultValue",
				Prompt: &survey.Input{Message: "Default Value (optional):"},
			},
			{
				Name:   "Required",
				Prompt: &survey.Confirm{Message: "Is this field required?", Default: false},
			},
			{
				Name:   "Validation",
				Prompt: &survey.Input{Message: "Validation (e.g. min=18, or regex):"},
			},
		}, &f)

		customFields = append(customFields, f)
	}

	if len(customFields) > 0 {
		fmt.Println("\n📋 Custom Fields Summary:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tTITLE\tTYPE\tENUM VALUES\tDEFAULT\tREQUIRED\tVALIDATION")
		fmt.Fprintln(w, "----\t-----\t----\t-----------\t-------\t--------\t----------")

		for _, f := range customFields {
			reqStr := "No"
			if f.Required {
				reqStr = "Yes"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", f.Name, f.Title, f.Type, f.EnumValues, f.DefaultValue, reqStr, f.Validation)
		}
		w.Flush()
		fmt.Println()
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

	// Asignar a las variables globales de Clonator
	entitySingular = answers.Singular
	entityPlural = answers.Plural

	// Ejecutar la lógica de clonado del comando principal re-utilizando cloneCmd
	cloneCmd.Run(cloneCmd, []string{answers.ServiceName})
}
