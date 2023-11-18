package models

import (
	"gorm.io/gorm"
)

type Skills struct {
	SkillID     int    `gorm:"column:skill_id;primaryKey"`
	Label       string `gorm:"column:label;not null"`
	Value       string `gorm:"column:value;not null"`
	Description string `gorm:"column:description;default:NULL"`
	ImageUrl    string `gorm:"column:image_url;default:NULL"`
	LevelID     int
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
	return "levels"
}

//   var users []User
//   err := db.Model(&User{}).Preload("CreditCard").Find(&users).Error
//   return users, err

// GetSkills retrieves all Skill records from the database.
func GetSkills(db *gorm.DB, skills *[]Skills) error {
	err := db.Model(&Skills{}).Preload("Levels").Find(&skills).Error
	// err := db.
	// Table("skills s").
	// Select("s.skill_id, s.label, s.value, s.description, s.image_url, l.level_id, l.name").
	// Joins("LEFT JOIN levels l ON s.level_id = l.level_id").
	// Find(&skills).
	// Error
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
