package models

import (
	"gorm.io/gorm"
)

type Skill struct {
	SkillID     int    `gorm:"column:skill_id;primaryKey" json:"skill_id"`
	Name        string `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description string `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl    string `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
	Type        string `gorm:"column:type; default:NULL; type:ENUM('SOFT','HARD');" json:"type" binding:"max=100"`
	// Careers     []CareerInSkills `gorm:"many2many:careers_skills;foreignKey:SkillID;joinForeignKey:SkillID;References:CareerID;joinReferences:CareerID" json:"careers"`
	Careers      []CareerInSkills `gorm:"many2many:careers_skills;foreignKey:SkillID;joinForeignKey:SkillID;References:CareerID;joinReferences:CareerID" json:"careers"`
	SkillsLevels []SkillsLevels   `gorm:"foreignKey:SkillID; References:SkillID;" json:"skills_levels"`
	// Levels     []Level   `gorm:"foreignKey:LevelID"`
}

type CareerInSkills struct {
	CareerID    int64  `gorm:"column:career_id;primaryKey;autoIncrement;" json:"career_id"`
	Name        string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
	Description string `gorm:"column:description; default:NULL; type:LONGTEXT;"  json:"description" binding:"max=1500"`
}

type UpdateSkillModels struct {
	Name         string           `gorm:"column:name; type:VARCHAR(255);" json:"name" binding:"max=255"`
	Description  string           `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl     string           `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
	Type         string           `gorm:"column:type; default:NULL; type:ENUM('SOFT','HARD');" json:"type" binding:"max=100"`
	Careers      []CareerInSkills `gorm:"many2many:careers_skills;foreignKey:SkillID;joinForeignKey:SkillID;References:CareerID;joinReferences:CareerID" json:"careers"`
	SkillsLevels []SkillsLevels   `gorm:"foreignKey:SkillID" json:"skills_levels"`
	// LevelID     int    `json:"level_id" binding:"required"`
}

type Levels struct {
	LevelID int    `gorm:"column:level_id; primaryKeyl" json:"level_id"`
	Name    string `gorm:"column:name; not null; type:VARCHAR(255);" json:"name"`
}

func (Skill) Tablename() string {
	return "skills"
}

func (CareerInSkills) TableName() string {
	return "careers"
}

// GetSkills retrieves all Skill records from the database with pagination.
func GetSkills(db *gorm.DB, page, limit int, skills *[]Skill) (pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Preload("Careers").Preload("SkillsLevels").Model(&Skill{}).
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
	err = db.Where("skill_id = ?", id).Preload("Careers").Preload("SkillsLevels").First(Skill).Error
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

// ยังไม่ได้เช็คอัพเดตว่าใช้ได้มั้ย
// UpdateSkill updates an existing Skill record in the database.
func UpdateSkill(db *gorm.DB, updatedSkill *Skill) (err error) {
	existingSkill := &Skill{}

	// Begin a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the skill record exists
	err = tx.Preload("Careers").Preload("SkillsLevels").First(existingSkill, "skill_id = ?", updatedSkill.SkillID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// // Update only the specified fields if they are not empty
	// if updatedSkill.Name != "" {
	// 	existingSkill.Name = updatedSkill.Name
	// }

	// if updatedSkill.Description != "" {
	// 	existingSkill.Description = updatedSkill.Description
	// }

	// if updatedSkill.ImageUrl != "" {
	// 	existingSkill.ImageUrl = updatedSkill.ImageUrl
	// }

	// if updatedSkill.Type != "" {
	// 	existingSkill.Type = updatedSkill.Type
	// }

	// fmt.Println(existingSkill.SkillID)
	// // existingSkill.SkillID = updatedSkill.SkillID

	// err = tx.Model(existingSkill).Association("SkillLevels").Clear()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	db.Save(updatedSkill)

	// Clear existing associations within the transaction
	err = tx.Model(existingSkill).Association("Careers").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update existing groups with the new one (if provided)
	if len(updatedSkill.Careers) > 0 {
		err = tx.Model(existingSkill).Association("Careers").Append(updatedSkill.Careers)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Clear existing associations within the transaction
	err = tx.Model(existingSkill).Association("SkillsLevels").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update existing groups with the new ones (if provided)
	if len(updatedSkill.SkillsLevels) > 0 {
		err = tx.Model(existingSkill).Association("SkillsLevels").Append(updatedSkill.SkillsLevels)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
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
