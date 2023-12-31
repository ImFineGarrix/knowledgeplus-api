package models

import (
	"gorm.io/gorm"
)

type Category struct {
	CategoryID int64     `gorm:"column:category_id; primaryKey;autoIncrement" json:"category_id"`
	Name       string    `gorm:"column:name; not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	ImageUrl   string    `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
	Careers    []Careers `gorm:"many2many:categories_careers;foreignKey:CategoryID;joinForeignKey:CategoryID;References:CareerID;joinReferences:CareerID" json:"-"`
}

type Careers struct {
	CareerID    int64  `gorm:"column:career_id;primaryKey;autoIncrement;" json:"career_id"`
	Name        string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
	Description string `gorm:"column:description; default:NULL; type:LONGTEXT;"  json:"description" binding:"max=1500"`
}

func (Category) TableName() string {
	return "categories"
}

type UpdateCategoryModels struct {
	CategoryID int64  `gorm:"column:category_id; primaryKey;autoIncrement" json:"category_id"`
	Name       string `gorm:"column:name; not null; type:VARCHAR(255);" json:"name" binding:"max=255"`
	ImageUrl   string `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
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

// UpdateCategory updates a Category record by ID.
func UpdateCategory(db *gorm.DB, category *Category) (err error) {
	err = db.Save(category).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteCategoryById deletes a Category record by its ID from the database.
func DeleteCategoryById(db *gorm.DB, id int) (err error) {
	err = db.Where("category_id = ?", id).Delete(&Category{}).Error
	if err != nil {
		return err
	}
	return nil
}
