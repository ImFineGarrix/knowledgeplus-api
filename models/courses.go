package models

import (
	"gorm.io/gorm"
)

type Course struct {
	CourseID        int                     `gorm:"column:course_id;primaryKey" json:"course_id"`
	Name            string                  `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description     string                  `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	LearnHours      string                  `gorm:"column:learn_hours; default:NULL; type:VARCHAR(45);" json:"learn_hours"`
	AcademicYear    string                  `gorm:"column:academic_year; default:NULL; type:VARCHAR(45);" json:"academic_year"`
	CourseLink      string                  `gorm:"column:course_link; default:NULL; type:LONGTEXT;" json:"course_link"`
	LearningOutcome string                  `gorm:"column:learinig_outcome; default:NULL; type:LONGTEXT;" json:"learning_outcome"`
	Organization    OrganizationInCourses   `gorm:"foreignKey:OrganizationID;references:OrganizationID" json:"organization"`
	OrganizationID  int                     `gorm:"column:organization_id" json:"organization_id"`
	SkillsLevels    []SkillsLevelsInCourses `gorm:"foreignKey:CourseID; References:CourseID;" json:"skills_levels"`
}

type CourseWithoutSkillLevels struct {
	CourseID        int                     `gorm:"column:course_id;primaryKey" json:"course_id"`
	Name            string                  `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description     string                  `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	LearnHours      string                  `gorm:"column:learn_hours; default:NULL; type:VARCHAR(45);" json:"learn_hours"`
	AcademicYear    string                  `gorm:"column:academic_year; default:NULL; type:VARCHAR(45);" json:"academic_year"`
	CourseLink      string                  `gorm:"column:course_link; default:NULL; type:LONGTEXT;" json:"course_link"`
	LearningOutcome string                  `gorm:"column:learinig_outcome; default:NULL; type:LONGTEXT;" json:"learning_outcome"`
	Organization    OrganizationInCourses   `gorm:"foreignKey:OrganizationID;references:OrganizationID" json:"organization"`
	OrganizationID  int                     `gorm:"column:organization_id" json:"organization_id"`
	SkillsLevels    []SkillsLevelsInCourses `gorm:"foreignKey:CourseID; References:CourseID;" json:"-"`
}

type OrganizationInCourses struct {
	OrganizationID int    `gorm:"column:organization_id; primaryKey;" json:"-"`
	Name           string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"max=255"`
	Description    string `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl       string `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
}

type SkillsLevelsInCourses struct {
	SkillsLevelsID int            `gorm:"column:skills_levels_id; primaryKey; autoIncrement;" json:"skills_levels_id"`
	SkillID        *int           `gorm:"column:skill_id;" json:"skill_id"`
	KnowledgeDesc  string         `gorm:"column:knowledge_desc;" json:"knowledge_desc"`
	AbilityDesc    string         `gorm:"column:ability_desc;" json:"ability_desc"`
	LevelID        int            `gorm:"column:level_id; not null" json:"level_id"`
	CourseID       *int           `gorm:"column:course_id; not null;" json:"-"`
	CareerID       *int           `gorm:"column:career_id; not null;" json:"career_id"`
	Skill          SkillInCourses `gorm:"foreignKey:SkillID;references:SkillID" json:"skill"`
	Career         CareersInGroup `gorm:"foreignKey:CareerID;references:CareerID" json:"career"`
}

type SkillInCourses struct {
	SkillID      *int                   `gorm:"column:skill_id;primaryKey" json:"-"`
	Name         string                 `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description  string                 `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl     string                 `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
	Type         string                 `gorm:"column:type; default:NULL; type:ENUM('SOFT','HARD');" json:"type" binding:"max=100"`
	SkillsLevels []SkillsLevelsInSkills `gorm:"foreignKey:SkillID; References:SkillID;" json:"-"`
}

type CareerInCourse struct {
	CareerID     int                     `gorm:"column:career_id;primaryKey;autoIncrement;" json:"-"`
	Name         string                  `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
	Description  string                  `gorm:"column:description; default:NULL; type:LONGTEXT;"  json:"description" binding:"max=1500"`
	Groups       []GroupsInCareers       `gorm:"many2many:groups_careers;foreignKey:CareerID;joinForeignKey:CareerID;References:GroupID;joinReferences:GroupID" json:"-"`
	SkillsLevels []SkillsLevelsInCareers `gorm:"foreignKey:CareerID; References:CareerID;" json:"-"`
}

func (Course) TableName() string {
	return "courses"
}

func (CourseWithoutSkillLevels) TableName() string {
	return "courses"
}

func (OrganizationInCourses) TableName() string {
	return "organizations"
}

func (SkillsLevelsInCourses) TableName() string {
	return "skills_levels"
}

func (SkillInCourses) TableName() string {
	return "skills"
}

func (CareerInCourse) TableName() string {
	return "careers"
}

// GetCourses retrieves all Course records from the database with pagination.
func GetCourses(db *gorm.DB, page, limit int, courses *[]Course) (pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Career").Model(&Course{}).
		Offset(offset).Limit(limit).
		Find(&courses).Error
	if err != nil {
		return Pagination{}, err
	}

	// Calculate total pages
	var totalCount int64
	if err := db.Model(&Course{}).Count(&totalCount).Error; err != nil {
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

// GetCoursesBySkillId retrieves courses based on the provided SkillID with pagination.
func GetCoursesBySkillId(db *gorm.DB, skillID, page, limit int, courses *[]Course) (pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.
		Preload("Organization").
		Preload("SkillsLevels.Skill").
		Preload("SkillsLevels.Career").
		Joins("JOIN skills_levels ON courses.course_id = skills_levels.course_id").
		Where("skills_levels.skill_id = ?", skillID).
		Model(&Course{}).
		Offset(offset).Limit(limit).
		Find(courses).Error
	if err != nil {
		return Pagination{}, err
	}

	// Count total courses for pagination
	var totalCount int64
	if err := db.Model(&Course{}).
		Preload("Organization").
		Preload("SkillsLevels.Skill").
		Preload("SkillsLevels.Career").
		Joins("JOIN skills_levels ON courses.course_id = skills_levels.course_id").
		Where("skills_levels.skill_id = ?", skillID).
		Count(&totalCount).Error; err != nil {
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

// GetCoursesByCareerId retrieves courses based on the provided CareerID with pagination.
func GetCoursesByCareerId(db *gorm.DB, careerID, page, limit int, courses *[]CourseWithoutSkillLevels) (pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.
		Preload("Organization").
		Preload("SkillsLevels.Skill").
		Preload("SkillsLevels.Career").
		Joins("JOIN skills_levels ON courses.course_id = skills_levels.course_id").
		Where("skills_levels.career_id = ?", careerID).
		Model(&CourseWithoutSkillLevels{}).
		Offset(offset).Limit(limit).
		Find(courses).Error
	if err != nil {
		return Pagination{}, err
	}

	// Count total courses for pagination
	var totalCount int64
	if err := db.Model(&CourseWithoutSkillLevels{}).
		Preload("Organization").
		Preload("SkillsLevels.Skill").
		Preload("SkillsLevels.Career").
		Joins("JOIN skills_levels ON courses.course_id = skills_levels.course_id").
		Where("skills_levels.career_id = ?", careerID).
		Count(&totalCount).Error; err != nil {
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

// CreateCourse creates a new Course record in the database.
func CreateCourse(db *gorm.DB, course *Course) (err error) {
	err = db.Create(course).Error
	if err != nil {
		return err
	}
	return nil
}

// GetCourseById retrieves a Course by its ID from the database.
func GetCourseById(db *gorm.DB, course *Course, id int) (err error) {
	err = db.Where("course_id = ?", id).Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Career").
		First(course).
		Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateCourse updates an existing Course record in the database.
func UpdateCourse(db *gorm.DB, updatedCourse *Course) (err error) {
	// Begin a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the course record exists
	existingCourse := &Course{}
	err = tx.Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Career").First(existingCourse, "course_id = ?", updatedCourse.CourseID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete all existing SkillsLevels records for the specified course
	err = tx.Where("course_id = ?", existingCourse.CourseID).Delete(&SkillsLevels{}).Error
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

	// Save the updated course
	err = tx.Save(updatedCourse).Error
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

// DeleteCourseById deletes a Course record by its ID from the database.
func DeleteCourseById(db *gorm.DB, id int) (err error) {
	err = db.Where("course_id = ?", id).Delete(&Course{}).Error
	if err != nil {
		return err
	}
	return nil
}
