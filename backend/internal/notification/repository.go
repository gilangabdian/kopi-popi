package notification

import (
	"gorm.io/gorm"
)

type Repository interface {
	SaveNotification(notification *Notification) error
	GetNotificationsByUser(userID string) ([]Notification, error)
	MarkAsRead(id string, userID string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) SaveNotification(notification *Notification) error {
	return r.db.Create(notification).Error
}

func (r *repository) GetNotificationsByUser(userID string) ([]Notification, error) {
	var notifications []Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *repository) MarkAsRead(id string, userID string) error {
	// Memastikan hanya notif milik user tersebut yang bisa ditandai terbaca
	return r.db.Model(&Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true).Error
}
