package models

import (
	"gorm.io/gorm"
)

type Skills struct {
	SkillID     int    `gorm:"column:skill_id;primaryKey" json:"skill_id"`
	Label       string `gorm:"column:label;not null" json:"label"`
	Value       string `gorm:"column:value;not null" json:" value"`
	Description string `gorm:"column:description;default:NULL" json:"description"`
	ImageUrl    string `gorm:"column:image_url;default:NULL" json:"image_url"`
	LevelID     int
	Levels      Levels `gorm:"foreignKey:LevelID;references:LevelID"`
	// Levels     []Level   `gorm:"foreignKey:LevelID"`
}

func (Skills) TableName() string {
	return "Skills"
}

type Levels struct {
	LevelID int    `gorm:"column:level_id;primaryKey"`
	Name    string `gorm:"column:name;not null"`
}

func (Levels) TableName() string {
	return "Levels"
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

// CreateSkill creates a new Skill record in the database.
func CreateSkill(db *gorm.DB, skill *Skill) (err error) {
	// err = db.Omit("Categories").Create(Skill).Error
	err = db.Create(skill).Error
	if err != nil {
		return err
	}

	return nil
}
