package models

import "gorm.io/gorm"

type User struct {
	ID       int    `gorm:"column:user_id;primaryKey;autoIncrement;" json:"user_id"`
	Name     string `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Email    string `gorm:"column:email;not null; type:VARCHAR(100);" json:"email" binding:"required,max=255"`
	Password string `gorm:"column:password;not null; type:VARCHAR(200);" json:"password" binding:"max=255"`
	Role     string `gorm:"column:role;not null; type:ENUM('owner','admin');" json:"role" binding:"required,max=255"`
}

type UserLogin struct {
	Email    string `gorm:"column:email;not null; type:VARCHAR(100);" json:"email" binding:"max=255"`
	Password string `gorm:"column:password;not null; type:VARCHAR(200);" json:"password" binding:"max=255"`
}

type UserResponse struct {
	ID    int    `gorm:"column:user_id;primaryKey;autoIncrement;" json:"user_id"`
	Name  string `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Email string `gorm:"column:email;not null; type:VARCHAR(100);" json:"email" binding:"required,max=255"`
	Role  string `gorm:"column:role;not null; type:ENUM('owner','admin');" json:"role" binding:"required,max=255"`
}

func (User) TableName() string {
	return "INT371_02.users"
}

func (UserLogin) TableName() string {
	return "INT371_02.users"
}

func (UserResponse) TableName() string {
	return "INT371_02.users"
}

// GetAllUsers retrieves all user records from the database.
func GetAllUsers(db *gorm.DB, users *[]UserResponse) (err error) {
	err = db.Find(users).Error
	if err != nil {
		return err
	}
	return nil
}

// GetUserById retrieves a user by its ID from the database.
func GetUserById(db *gorm.DB, user *UserResponse, id uint) (err error) {
	err = db.Where("user_id = ?", id).First(user).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateUser creates a new user record in the database.
func CreateUser(db *gorm.DB, user *User) (err error) {
	err = db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates a user record by ID.
func UpdateUser(db *gorm.DB, user *UserResponse) (err error) {
	// Find the existing user by ID
	existingUser := User{}
	if err := db.First(&existingUser, user.ID).Error; err != nil {
		return err
	}

	// Update the fields
	// can't update password
	existingUser.Name = user.Name
	existingUser.Email = user.Email
	existingUser.Role = user.Role

	// Save the changes
	if err := db.Save(&existingUser).Error; err != nil {
		return err
	}

	return nil
}

// DeleteUserById deletes a user record by its ID from the database.
func DeleteUserById(db *gorm.DB, id uint) (err error) {
	err = db.Where("user_id = ?", id).Delete(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}
