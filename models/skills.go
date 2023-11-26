package models

import (
	"gorm.io/gorm"
)

type Skills struct {
	SkillID     int    `gorm:"column:skill_id;primaryKey" json:"skill_id"`
	Name        string `gorm:"column:name;not null" json:"name" binding:"required,len=45"`
	Description string `gorm:"column:description;default:NULL" json:"description" binding:"len=200"`
	ImageUrl    string `gorm:"column:image_url;default:NULL" json:"image_url" binding:"len=255"`
	LevelID     int    `json:"-"`
	Levels      Levels `gorm:"foreignKey:LevelID;references:LevelID"`
	// Levels     []Level   `gorm:"foreignKey:LevelID"`
}

func (Skills) TableName() string {
	return "skills"
}

type Levels struct {
	LevelID int    `gorm:"column:level_id;primaryKey"`
	Name    string `gorm:"column:name;not null"`
}

func (Levels) TableName() string {
	return "Levels"
}

// GetSkills retrieves all Skill records from the database.
func GetSkills(db *gorm.DB, skills *[]Skills) error {
	err := db.Model(&Skills{}).Preload("Levels").Find(&skills).Error
	if err != nil {
		return err
	}
	return nil
}

// GetSkillById retrieves a Skill by its ID from the database.
func GetSkillById(db *gorm.DB, Skill *Skills, id int) (err error) {
	err = db.Where("skill_id = ?", id).Preload("Levels").First(Skill).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateSkill creates a new Skill record in the database.
func CreateSkill(db *gorm.DB, skill *Skills) (err error) {
	// err = db.Omit("Categories").Create(Skill).Error
	err = db.Create(skill).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateSkill updates a Skill record by ID.
func UpdateSkill(db *gorm.DB, skill *Skills) (err error) {
	err = db.Save(skill).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteSkillById deletes a Skill record by its ID from the database.
func DeleteSkillById(db *gorm.DB, id int) (err error) {
	err = db.Where("skill_id = ?", id).Delete(&Skills{}).Error
	if err != nil {
		return err
	}
	return nil
}
