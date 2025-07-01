package user

type Repository interface {
	CreateUser(login, password string) (User, error)
	GetUserByCredentials(login, password string) (User, error)
	GetUserByID(userID string) (User, error)
	Withdrawal(userID, orderNumber string, sum float64) error
}
