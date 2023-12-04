package models

import (
	"gorm.io/gorm"
)

type Skill struct {
	SkillID     int       `gorm:"column:skill_id;primaryKey" json:"skill_id"`
	Name        string    `gorm:"column:name;not null" json:"name" binding:"required,max=45"`
	Description string    `gorm:"column:description;default:NULL" json:"description" binding:"max=200"`
	ImageUrl    string    `gorm:"column:image_url;default:NULL" json:"image_url" binding:"max=255"`
	LevelID     int       `json:"level_id"`
	Levels      Levels    `gorm:"foreignKey:LevelID;references:LevelID" json:"-"`
	Careers     []Careers `gorm:"many2many:careers_skills;foreignKey:SkillID;joinForeignKey:SkillID;References:CareerID;joinReferences:CareerID" json:"-"`
	// Levels     []Level   `gorm:"foreignKey:LevelID"`
}

func (Skill) Tablename() string {
	return "skills"
}

type Levels struct {
	LevelID int    `gorm:"column:level_id;primaryKey" json:"level_id"`
	Name    string `gorm:"column:name;not null" json:"name"`
}

func (Levels) Tablename() string {
	return "levels"
}

// GetSkills retrieves all Skill records from the database with pagination.
func GetSkills(db *gorm.DB, page, limit int, skills *[]Skill) (pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Model(&Skill{}).Preload("Levels").
		Offset(offset).Limit(limit).
		Find(&skills).Error
	if err != nil {
		return Pagination{}, err
	}

	// Calculate total pages
	var totalCount int64
	if err := db.Model(&Skill{}).Count(&totalCount).Error; err != nil {
		return Pagination{}, err
	}

	totalPages := int(totalCount)

	pagination = Pagination{
		Page:  page,
		Total: totalPages,
		Limit: limit,
	}

	return pagination, nil
}

// GetSkillById retrieves a Skill by its ID from the database.
func GetSkillById(db *gorm.DB, Skill *Skill, id int) (err error) {
	err = db.Where("skill_id = ?", id).Preload("Levels").First(Skill).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateSkill creates a new Skill record in the database.
func CreateSkill(db *gorm.DB, skill *Skill) (err error) {
	// err = db.Omit("Categories").Create(Skill).Error
	err = db.Create(skill).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateSkill updates a Skill record by ID.
func UpdateSkill(db *gorm.DB, skill *Skill) (err error) {
	err = db.Save(skill).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteSkillById deletes a Skill record by its ID from the database.
func DeleteSkillById(db *gorm.DB, id int) (err error) {
	err = db.Where("skill_id = ?", id).Delete(&Skill{}).Error
	if err != nil {
		return err
	}
	return nil
}
