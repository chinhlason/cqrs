package command

import (
	"gorm.io/gorm"
	"time"
)

type IRepository interface {
	InsertUser(name, email string) error
	UpdateUser(id int64, name, email string) error

	InsertOrder(userId int64, product string) error
	UpdateOrder(orderId int64, product string) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{db}
}

func (r *Repository) InsertUser(name, email string) error {
	return r.db.Create(&User{Name: name, Email: email, CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
}

func (r *Repository) UpdateUser(id int64, name, email string) error {
	updates := make(map[string]interface{})
	if name != "" {
		updates["name"] = name
	}
	if email != "" {
		updates["email"] = email
	}
	return r.db.Model(&User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) InsertOrder(userId int64, product string) error {
	order := Order{
		UserId:    userId,
		Product:   product,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := r.db.Create(&order).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateOrder(orderId int64, product string) error {
	newOrder := Order{
		Product:   product,
		UpdatedAt: time.Now(),
	}
	if err := r.db.Model(&Order{}).Where("id = ?", orderId).Updates(&newOrder).Error; err != nil {
		return err
	}
	return nil
}
