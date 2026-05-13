package user

type UserRegistration struct {
	UserName string
	Password string
}

func NewUserRegistration(username, password string) UserRegistration {
	return UserRegistration{
		UserName: username,
		Password: password,
	}
}
