package models

import "gorm.io/gorm"

type CareerCategory struct {
	CategoriesCareersID int `gorm:"primaryKey;autoIncrement" json:"categories_careers_id"`
	CareerID            int `gorm:"foreignKey:CareerID ;references:CareerID "`
	CategoryID          int `gorm:"foreignKey:CategoryID ;references:CategoryID "`
}

// CreateCareer creates a new Career record in the database.
func CreateCareerCategory(db *gorm.DB, careerCategory *CareerCategory) (err error) {
	err = db.Create(careerCategory).Error
	if err != nil {
		return err
	}
	return nil
}
