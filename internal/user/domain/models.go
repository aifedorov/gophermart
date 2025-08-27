package domain

type User struct {
	ID       string
	Login    string
	Password string
}

type RegisterRequest struct {
	Login    string
	Password string
}

type LoginRequest struct {
	Login    string
	Password string
}
