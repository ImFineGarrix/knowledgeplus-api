package models

type User struct {
	ID       uint   `gorm:"column:user_id;primaryKey;autoIncrement;" json:"user_id"`
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
	ID    uint   `gorm:"column:user_id;primaryKey;autoIncrement;" json:"user_id"`
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
