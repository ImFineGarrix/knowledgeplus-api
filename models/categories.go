package models

import (
	"gorm.io/gorm"
)

type Category struct {
	CategoryID int      `gorm:"primaryKey" json:"category_id"`
	Name      string   `gorm:"not null" json:"name"`
	ImageUrl   string   `gorm:"default:NULL" json:"image_url"`
	Careers    []Career `gorm:"many2many:categories_careers;foreignKey:CategoryID;joinForeignKey:CategoryID;References:CareerID;joinReferences:CareerID"`
}

func (Category) TableName() string {
	return "Categories"
}

// GetCategories retrieves all Category records from the database.
func GetCategories(db *gorm.DB, Category *[]Category) (err error) {
	err = db.Preload("Careers").Find(Category).Error
	if err != nil {
		return err
	}
	return nil
}

// GetCategoryById retrieves a Category by its ID from the database.
func GetCategoryById(db *gorm.DB, Category *Category, id int) (err error) {
	err = db.Where("category_id = ?", id).Preload("Careers").First(Category).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateCareer creates a new Career record in the database.
func CreateCategory(db *gorm.DB, category *Category) (err error) {
	err = db.Create(category).Error
	if err != nil {
		return err
	}

	return nil
}
