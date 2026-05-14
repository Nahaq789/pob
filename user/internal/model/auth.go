package model

type Auth struct {
	UserName      string
	PasswordPlane string
}

func NewAuth(userName, passwordPlane string) Auth {
	return Auth{
		UserName:      userName,
		PasswordPlane: passwordPlane,
	}
}
