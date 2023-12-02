package models

import (
	"gorm.io/gorm"
)

type Career struct {
	CareerID    int64        `gorm:"primaryKey;autoIncrement" json:"career_id"`
	Name        string       `gorm:"not null" json:"name" binding:"max=45"`
	Description string       `gorm:"default:NULL" json:"description" binding:"max=255"`
	ShortDesc   string       `gorm:"default:NULL" json:"short_desc" binding:"max=50"`
	Categories  []Categories `gorm:"many2many:categories_careers;foreignKey:CareerID;joinForeignKey:CareerID;References:CategoryID;joinReferences:CategoryID" json:"categories"`
	Skills      []Skills     `gorm:"many2many:careers_skills;foreignKey:CareerID;joinForeignKey:CareerID;References:SkillID;joinReferences:SkillID" json:"skills"`
}

type Categories struct {
	CategoryID int64  `gorm:"primaryKey" json:"category_id"`
	Name       string `gorm:"not null" json:"name"`
	ImageUrl   string `gorm:"default:NULL" json:"image_url"`
}

type Skills struct {
	SkillID     int    `gorm:"column:skill_id;primaryKey" json:"skill_id"`
	Name        string `gorm:"column:name;not null" json:"name" binding:"max=45"`
	Description string `gorm:"column:description;default:NULL" json:"description" binding:"max=200"`
	ImageUrl    string `gorm:"column:image_url;default:NULL" json:"image_url" binding:"max=255"`
	LevelID     int    `json:"-"`
	Levels      Levels `gorm:"foreignKey:LevelID;references:LevelID"`
}

func (Career) TableName() string {
	return "careers"
}

func (Skills) Tablename() string {
	return "skills"
}

// totals

// GetCareers retrieves all Career records from the database with pagination.
func GetCareers(db *gorm.DB, page, limit int) (careers []Career, pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Preload("Categories").Preload("Skills").Preload("Skills.Levels").
		Offset(offset).Limit(limit).
		Find(&careers).Error
	if err != nil {
		return nil, Pagination{}, err
	}

	// Calculate total pages
	var totalCount int64
	if err := db.Model(&Career{}).Count(&totalCount).Error; err != nil {
		return nil, Pagination{}, err
	}

	totalPages := int(totalCount)

	pagination = Pagination{
		Page:  page,
		Total: totalPages,
		Limit: limit,
	}

	return careers, pagination, nil
}

// GetCareersWithHaveCategories retrieves careers where category_id in many-to-many is not null with pagination.
func GetCareersWithHaveCategories(db *gorm.DB, page, limit int) (careers []Career, pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Joins("JOIN categories_careers ON careers.career_id = categories_careers.career_id").
		Where("categories_careers.category_id IS NOT NULL").
		Preload("Categories").Preload("Skills").
		Offset(offset).Limit(limit).
		Find(&careers).Error
	if err != nil {
		return nil, Pagination{}, err
	}

	// Calculate total pages
	var totalCount int64
	if err := db.Joins("JOIN categories_careers ON careers.career_id = categories_careers.career_id").
		Where("categories_careers.category_id IS NOT NULL").
		Model(&Career{}).Count(&totalCount).Error; err != nil {
		return nil, Pagination{}, err
	}

	totalPages := int(totalCount)

	pagination = Pagination{
		Page:  page,
		Total: totalPages,
		Limit: limit,
	}

	return careers, pagination, nil
}

// Pagination struct
type Pagination struct {
	Page  int `json:"page"`
	Total int `json:"total"`
	Limit int `json:"limit"`
}

// GetCareerById retrieves a Career by its ID from the database.
func GetCareerById(db *gorm.DB, Career *Career, id int) (err error) {
	err = db.Where("career_id = ?", id).Preload("Categories").Preload("Skills").First(Career).Error
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

// UpdateCareer updates an existing Career record in the database.
func UpdateCareer(db *gorm.DB, career *Career) (err error) {
	existingCareer := &Career{}

	// Begin a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the career record exists
	err = tx.Where("career_id = ?", career.CareerID).Preload("Categories").Preload("Skills.Levels").First(existingCareer).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update the fields of the existing career record
	existingCareer.Name = career.Name
	existingCareer.Description = career.Description
	existingCareer.ShortDesc = career.ShortDesc

	// Clear existing associations within the transaction
	err = tx.Model(existingCareer).Association("Categories").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(existingCareer).Association("Skills").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Replace existing categories with new ones
	err = tx.Model(existingCareer).Association("Categories").Replace(career.Categories)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Replace existing skills with new ones
	err = tx.Model(existingCareer).Association("Skills").Replace(career.Skills)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit().Error
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
