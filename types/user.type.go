package types

type RegisterUser struct {
	Name     string `json:"name" validate:"required" name:"name"`
	Email    string `json:"email" validate:"required,email" name:"email"`
	Username string `json:"username" validate:"required" name:"username"`
}

type AuthHeader struct {
	Authorization string `header:"Authorization"`
}
