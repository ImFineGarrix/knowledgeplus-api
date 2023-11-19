package models

import (
	"gorm.io/gorm"
)

type Career struct {
	CareerID    int        `gorm:"primaryKey;autoIncrement" json:"career_id"`
	Label       string     `gorm:"not null" json:"label"`
	Value       string     `gorm:"not null" json:"value"`
	Description string     `gorm:"default:NULL" json:"description"`
	ShortDesc   string     `gorm:"default:NULL" json:"short_desc"`
	Categories  []Category `gorm:"many2many:categories_careers;foreignKey:CareerID;joinForeignKey:CareerID;References:CategoryID;joinReferences:CategoryID" json:"categories"`
}

func (Career) TableName() string {
	return "Careers"
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
