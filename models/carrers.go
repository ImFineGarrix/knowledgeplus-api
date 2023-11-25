package models

import (
	"gorm.io/gorm"
)

type Career struct {
	CareerID    int64        `gorm:"primaryKey;autoIncrement" json:"career_id"`
	Name        string       `gorm:"not null" json:"name"`
	Description string       `gorm:"default:NULL" json:"description"`
	ShortDesc   string       `gorm:"default:NULL" json:"short_desc"`
	Categories  []Categories `gorm:"many2many:categories_careers;foreignKey:CareerID;joinForeignKey:CareerID;References:CategoryID;joinReferences:CategoryID" json:"categories"`
}

type Categories struct {
	CategoryID int64  `gorm:"primaryKey" json:"category_id"`
	Name       string `gorm:"not null" json:"name"`
	ImageUrl   string `gorm:"default:NULL" json:"image_url"`
}

func (Career) TableName() string {
	return "careers"
}

// GetCareers retrieves all Career records from the database.
func GetCareers(db *gorm.DB, Careers *[]Career) (err error) {
	err = db.Preload("Categories").Find(Careers).Error
	if err != nil {
		return err
	}
	return nil
}

// GetCareerById retrieves a Career by its ID from the database.
func GetCareerById(db *gorm.DB, Career *Career, id int) (err error) {
	err = db.Where("career_id = ?", id).Preload("Categories").First(Career).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateCareer creates a new Career record in the database.
func CreateCareer(db *gorm.DB, career *Career) (err error) {
	// err = db.Omit("Categories").Create(career).Error
	err = db.Create(career).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateCareer updates a Career record in the database.
func UpdateCareer(db *gorm.DB, career *Career) (err error) {
	err = db.Save(career).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteCareerById deletes a Career by its ID from the database.
func DeleteCareerById(db *gorm.DB, id int) (err error) {
	err = db.Where("career_id = ?", id).Delete(&Career{}).Error
	if err != nil {
		return err
	}
	return nil
}
