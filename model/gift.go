package model

type Gift struct {
	ID          int
	Name        string
	Description string
	Picture     string
	Price       int
	Count       int
}

func (g Gift) TableName() string {
	return "inventory"
}
