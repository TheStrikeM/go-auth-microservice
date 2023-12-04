package modules

var (
	ErrValidationAuth     = "Validation error"
	ErrUserExists         = "User already exists"
	ErrUserAdding         = "Error in AddUser"
	ErrInvalidCredentials = "Invalid credentials"
	ErrPasswordHash       = "Error in password hash"
)

type ModuleInitializer struct {
}
