package user

import "fmt"

const (
	UniqueConstraintUserEmail = "users_email_key"
)

type EmailAlreadyExistsError struct {
	Email string
}

func (e *EmailAlreadyExistsError) Error() string {
	return fmt.Sprintf("Email %s already exists", e.Email)
}
