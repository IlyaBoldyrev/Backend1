package shop

// Модель товара
type Item struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"` // цена в копейках/центах
	ImageLink   string `json:"image_link"`
}

// Модель юзера
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
