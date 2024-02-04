// package models

// import "gorm.io/gorm"

// type SkillsLevels struct {
// 	SkillsLevelsID int64     `gorm:"column:skills_levels_id; primaryKey; autoIncrement" json:"skills_levels_id"`
// 	SkillID        int64     `gorm:"column:skill_id; not null" json:"skill_id"`
// 	KnowledgeDesc  string    `gorm:"column:knowledge_desc" json:"knowledge_desc"`
// 	AbilityDesc    string    `gorm:"column:ability_desc" json:"ability_desc"`
// 	LevelsID       int64     `gorm:"column:levels_id; not null" json:"levels_id"`
// 	CourseID       *int64    `gorm:"column:course_id" json:"course_id"`
// 	Skill          Skill     `gorm:"foreignKey:SkillID" json:"skill"`
// 	Levels         Level     `gorm:"foreignKey:LevelsID" json:"level"`
// 	Course         *Course   `gorm:"foreignKey:CourseID" json:"course,omitempty"`
// }

// // Skill represents the skill table
// type Skill struct {
// 	SkillID int64  `gorm:"column:skill_id; primaryKey; autoIncrement" json:"skill_id"`
// 	Name    string `gorm:"column:name; not null" json:"name"`
// 	// Add other fields as needed
// }

// // Level represents the levels table
// type Level struct {
// 	LevelID int64  `gorm:"column:level_id; primaryKey" json:"level_id"`
// 	Level   string `gorm:"column:level; not null" json:"level"`
// 	// Add other fields as needed
// }

// // Course represents the courses table
// type Course struct {
// 	CourseID     int64  `gorm:"column:course_id; primaryKey; autoIncrement" json:"course_id"`
// 	Name         string `gorm:"column:name; not null" json:"name"`
// 	Description  string `gorm:"column:description" json:"description"`
// 	LearnHours   string `gorm:"column:learn_hours" json:"learn_hours"`
// 	AcademicYear string `gorm:"column:academic_year" json:"academic_year"`
// 	CourseLink   string `gorm:"column:course_link" json:"course_link"`
// 	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
// 	OrganizationID int64 `gorm:"column:organization_id; not null" json:"-"`
// }

// // Organization represents the organizations table
// type Organization struct {
// 	OrganizationID int64  `gorm:"column:organization_id; primaryKey; autoIncrement" json:"organization_id"`
// 	Name           string `gorm:"column:name; not null" json:"name"`
// 	Description    string `gorm:"column:description" json:"description"`
// 	ImageURL       string `gorm:"column:image_url" json:"image_url"`
// 	// Add other fields as needed
// }

package models

import "gorm.io/gorm"

type SkillsLevels struct {
	SkillsLevelsID int64  `gorm:"column:skills_levels_id; primaryKey; autoIncrement;" json:"skills_levels_id"`
	SkillID        *int64 `gorm:"column:skill_id; not null;" json:"-"`
	KnowledgeDesc  string `gorm:"column:knowledge_desc;" json:"knowledge_desc"`
	AbilityDesc    string `gorm:"column:ability_desc;" json:"ability_desc"`
	LevelsID       int64  `gorm:"column:levels_id; not null" json:"levels_id"`
	CourseID       *int64 `gorm:"column:course_id; not null;" json:"-"`
}

func (SkillsLevels) Tablename() string {
	return "skills_levels"
}

// GetSkillsLevels retrieves all records from the skills_levels table.
func GetSkillsLevels(db *gorm.DB, skillsLevels *[]SkillsLevels) (err error) {
	err = db.Find(skillsLevels).Error
	if err != nil {
		return err
	}
	return nil
}

//อาจจะต้องทำ getSkillsLevelsBySkillId,getSkillsLevelsByCourseId
