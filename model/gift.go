package model

type Gift struct {
	ID          int    `gorm:"id,primaryKey"`
	Name        string `gorm:"name"`
	Description string `gorm:"description"`
	Picture     string `gorm:"picture"`
	Price       int    `gorm:"price"`
	Count       int    `gorm:"count"`
}

func (g Gift) TableName() string {
	return "inventory"
}
