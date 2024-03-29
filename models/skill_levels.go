package models

import "gorm.io/gorm"

type SkillsLevels struct {
	SkillsLevelsID int               `gorm:"column:skills_levels_id; primaryKey; autoIncrement;" json:"skills_levels_id"`
	SkillID        *int              `gorm:"column:skill_id;" json:"skill_id"`
	KnowledgeDesc  string            `gorm:"column:knowledge_desc;" json:"knowledge_desc"`
	AbilityDesc    string            `gorm:"column:ability_desc;" json:"ability_desc"`
	LevelID        int               `gorm:"column:level_id; not null" json:"level_id"`
	Skill          SkillInCourses    `gorm:"foreignKey:SkillID;references:SkillID" json:"skill"`
	Careers        []CareerInCourses `gorm:"many2many:careers_skills_levels;foreignKey:SkillsLevelsID;joinForeignKey:SkillsLevelsID;References:CareerID;joinReferences:CareerID" json:"-"`
	Courses        []CourseInCareers `gorm:"many2many:courses_skills_levels;foreignKey:SkillsLevelsID;joinForeignKey:SkillsLevelsID;References:CourseID;joinReferences:CourseID" json:"-"`
}

func (SkillsLevels) TableName() string {
	return "skills_levels"
}

// GetAllSkillsLevels retrieves all SkillsLevels records from the database with pagination.
func GetAllSkillsLevels(db *gorm.DB, page, limit int, skillsLevels *[]SkillsLevels) (pagination Pagination, err error) {
	offset := (page - 1) * limit

	// Preload the related Skill and Careers
	err = db.Preload("Skill").Preload("Careers").Preload("Courses").Model(&SkillsLevels{}).
		Offset(offset).Limit(limit).
		Find(skillsLevels).Error
	if err != nil {
		return Pagination{}, err
	}

	// Calculate total pages
	var totalCount int64
	if err := db.Model(&SkillsLevels{}).Count(&totalCount).Error; err != nil {
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
