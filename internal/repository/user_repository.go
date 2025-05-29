package repository

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    "github.com/moha/kaafipay-backend/internal/models"
)

type UserRepository interface {
    Create(user *models.User) error
    FindByID(id uuid.UUID) (*models.User, error)
    FindByPhone(phone string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uuid.UUID) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
    var user models.User
    if err := r.db.First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
    var user models.User
    if err := r.db.First(&user, "phone = ?", phone).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
    return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
    return r.db.Delete(&models.User{}, "id = ?", id).Error
} 