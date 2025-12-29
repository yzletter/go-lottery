package model

type Order struct {
	ID     int
	UserID int
	GiftID int
}

func (o Order) TableName() string {
	return "orders"
}
