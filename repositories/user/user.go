package user

import (
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

type UserInterface interface {
	Register(newUser models.User) (models.User, error)
	Login(email string) (models.User, error)
	Get(userId int) (models.User, error)
	Update(newUser models.User, userId int) (models.User, error)
	Delete(userId int) (models.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Register(newUser models.User) (models.User, error) {
	err := ur.db.Save(&newUser).Error
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (ur *UserRepository) Login(email string) (models.User, error) {
	var user models.User
	var err = ur.db.Preload("Partner").First(&user, "email = ?", email).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (ur *UserRepository) Get(userId int) (models.User, error) {
	user := models.User{}
	if err := ur.db.Preload("Partner").First(&user, userId).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (ur *UserRepository) Update(newUser models.User, userId int) (models.User, error) {
	user := models.User{}
	if err := ur.db.First(&user, "id=?", userId).Error; err != nil {
		return newUser, err
	}
	ur.db.Model(&user).Updates(newUser)
	return newUser, nil
}

func (ur *UserRepository) Delete(userId int) (models.User, error) {
	user := models.User{}
	if err := ur.db.First(&user, "id=?", userId).Error; err != nil {
		return user, err
	}
	ur.db.Delete(&user)
	return user, nil
}
