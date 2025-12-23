package todoctrl

import "gorm.io/gorm"

type Controller struct {
	db *gorm.DB
}

func NewTodoController(db *gorm.DB) *Controller {
	return &Controller{db: db}
}
