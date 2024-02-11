package models

import (
	"gorm.io/gorm"
)

type Skill struct {
	SkillID     *int   `gorm:"column:skill_id;primaryKey" json:"skill_id"`
	Name        string `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description string `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl    string `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
	Type        string `gorm:"column:type; default:NULL; type:ENUM('SOFT','HARD');" json:"type" binding:"max=100"`
	// Careers     []CareerInSkills `gorm:"many2many:careers_skills;foreignKey:SkillID;joinForeignKey:SkillID;References:CareerID;joinReferences:CareerID" json:"careers"`
	SkillsLevels []SkillsLevelsInSkills `gorm:"foreignKey:SkillID; References:SkillID;" json:"skills_levels"`
	// Levels     []Level   `gorm:"foreignKey:LevelID"`
}

type UpdateSkillModels struct {
	Name         string         `gorm:"column:name; type:VARCHAR(255);" json:"name" binding:"max=255"`
	Description  string         `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl     string         `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
	Type         string         `gorm:"column:type; default:NULL; type:ENUM('SOFT','HARD');" json:"type" binding:"max=100"`
	SkillsLevels []SkillsLevels `gorm:"foreignKey:SkillID" json:"skills_levels"`
	// LevelID     int    `json:"level_id" binding:"required"`
}

type SkillsLevelsInSkills struct {
	SkillsLevelsID int    `gorm:"column:skills_levels_id; primaryKey; autoIncrement;" json:"skills_levels_id"`
	SkillID        *int   `gorm:"column:skill_id;" json:"skill_id"`
	KnowledgeDesc  string `gorm:"column:knowledge_desc;" json:"knowledge_desc"`
	AbilityDesc    string `gorm:"column:ability_desc;" json:"ability_desc"`
	LevelsID       int    `gorm:"column:levels_id; not null" json:"levels_id"`
	CourseID       *int   `gorm:"column:course_id; not null;" json:"-"`
	CareerID       *int   `gorm:"column:career_id; not null;" json:"-"`
}

type Levels struct {
	LevelID int    `gorm:"column:level_id; primaryKeyl" json:"level_id"`
	Name    string `gorm:"column:name; not null; type:VARCHAR(255);" json:"name"`
}

func (Skill) Tablename() string {
	return "skills"
}

func (SkillsLevelsInSkills) TableName() string {
	return "skills_levels"
}

// GetSkills retrieves all Skill records from the database with pagination.
func GetSkills(db *gorm.DB, page, limit int, skills *[]Skill) (pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Preload("SkillsLevels").Model(&Skill{}).
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
	err = db.Where("skill_id = ?", id).Preload("SkillsLevels").First(Skill).Error
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
	err = tx.Preload("SkillsLevels").First(existingSkill, "skill_id = ?", updatedSkill.SkillID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete all existing SkillsLevels records for the specified skill
	err = tx.Where("skill_id = ?", existingSkill.SkillID).Delete(&SkillsLevels{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the delete operation
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// Begin a new transaction for the update operation
	tx = db.Begin()

	// Save the updated skill
	err = tx.Save(updatedSkill).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction for the update operation
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
