package order

type Repository interface {
	CreateOrder(userID, number string) (Order, error)
	GetOrdersByUserID(userID string) ([]Order, error)
	GetOrderByNumber(number string) (Order, error)
	UpdateOrderStatus(number string, status Status) error
}
