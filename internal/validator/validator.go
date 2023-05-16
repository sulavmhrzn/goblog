package validator

import "net/mail"

type Validator struct {
	Error map[string]interface{}
}

func New() *Validator {
	return &Validator{
		Error: make(map[string]interface{}),
	}
}

func (v *Validator) AddErrorMessage(title, message string) {
	if _, ok := v.Error[title]; !ok {
		v.Error[title] = message
	}
}

func (v *Validator) Check(ok bool, title string, message string) {
	if !ok {
		v.AddErrorMessage(title, message)
	}
}

func (v *Validator) IsValid() bool {
	return len(v.Error) == 0
}

func (v *Validator) ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
