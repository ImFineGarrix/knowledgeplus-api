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
	CourseType      string                  `gorm:"column:course_type; default:NULL; type:VARCHAR(45);" json:"course_type"`
	LearningOutcome string                  `gorm:"column:learinig_outcome; default:NULL; type:LONGTEXT;" json:"learning_outcome"`
	Organization    OrganizationInCourses   `gorm:"foreignKey:OrganizationID;references:OrganizationID" json:"organization"`
	OrganizationID  int                     `gorm:"column:organization_id" json:"organization_id"`
	SkillsLevels    []SkillsLevelsInCourses `gorm:"many2many:courses_skills_levels;foreignKey:CourseID;joinForeignKey:CourseID;References:SkillsLevelsID;joinReferences:SkillsLevelsID" json:"skills_levels"`
}

type CourseWithoutSkillLevels struct {
	CourseID        int                   `gorm:"column:course_id;primaryKey" json:"course_id"`
	Name            string                `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description     string                `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	LearnHours      string                `gorm:"column:learn_hours; default:NULL; type:VARCHAR(45);" json:"learn_hours"`
	AcademicYear    string                `gorm:"column:academic_year; default:NULL; type:VARCHAR(45);" json:"academic_year"`
	CourseLink      string                `gorm:"column:course_link; default:NULL; type:LONGTEXT;" json:"course_link"`
	LearningOutcome string                `gorm:"column:learinig_outcome; default:NULL; type:LONGTEXT;" json:"learning_outcome"`
	Organization    OrganizationInCourses `gorm:"foreignKey:OrganizationID;references:OrganizationID" json:"organization"`
	OrganizationID  int                   `gorm:"column:organization_id" json:"organization_id"`
	// SkillsLevels    []SkillsLevelsInCourses `gorm:"foreignKey:CourseID; References:CourseID;" json:"-"`
}

type OrganizationInCourses struct {
	OrganizationID int    `gorm:"column:organization_id; primaryKey;" json:"-"`
	Name           string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"max=255"`
	Description    string `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl       string `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=2000"`
}

type SkillsLevelsInCourses struct {
	SkillsLevelsID int    `gorm:"column:skills_levels_id; primaryKey; autoIncrement;" json:"skills_levels_id"`
	SkillID        *int   `gorm:"column:skill_id;" json:"skill_id"`
	KnowledgeDesc  string `gorm:"column:knowledge_desc;" json:"knowledge_desc"`
	AbilityDesc    string `gorm:"column:ability_desc;" json:"ability_desc"`
	LevelID        int    `gorm:"column:level_id; not null" json:"level_id"`
	// CourseID       *int             `gorm:"column:course_id;" json:"-"`
	// CareerID       *int             `gorm:"column:career_id;" json:"career_id"`
	Skill   SkillInCourses    `gorm:"foreignKey:SkillID;references:SkillID" json:"skill"`
	Careers []CareerInCourses `gorm:"many2many:careers_skills_levels;foreignKey:SkillsLevelsID;joinForeignKey:SkillsLevelsID;References:CareerID;joinReferences:CareerID" json:"careers"`
	// Career         *CareerInCourses `gorm:"foreignKey:CareerID; References:CareerID;" json:"career"`
}

type SkillInCourses struct {
	SkillID      *int                   `gorm:"column:skill_id;primaryKey" json:"-"`
	Name         string                 `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description  string                 `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl     string                 `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
	Type         string                 `gorm:"column:type; default:NULL; type:ENUM('SOFT','HARD');" json:"type" binding:"max=100"`
	SkillsLevels []SkillsLevelsInSkills `gorm:"foreignKey:SkillID; References:SkillID;" json:"-"`
}

type CareerInCourses struct {
	CareerID    int    `gorm:"column:career_id;primaryKey;autoIncrement;" json:"career_id"`
	Name        string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
	Description string `gorm:"column:description; default:NULL; type:LONGTEXT;"  json:"description" binding:"max=1500"`
	// Groups      []GroupsInCareers `gorm:"many2many:groups_careers;foreignKey:CareerID;joinForeignKey:CareerID;References:GroupID;joinReferences:GroupID" json:"-"`
	// SkillsLevels []SkillsLevelsInCareers `gorm:"foreignKey:CareerID; References:CareerID;" json:"-"`
}

type UpdateCourseDTO struct {
	CourseID        int                           `gorm:"column:course_id;primaryKey" json:"course_id"`
	Name            string                        `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description     string                        `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	LearnHours      string                        `gorm:"column:learn_hours; default:NULL; type:VARCHAR(45);" json:"learn_hours"`
	AcademicYear    string                        `gorm:"column:academic_year; default:NULL; type:VARCHAR(45);" json:"academic_year"`
	CourseLink      string                        `gorm:"column:course_link; default:NULL; type:LONGTEXT;" json:"course_link"`
	CourseType      string                        `gorm:"column:course_type; default:NULL; type:VARCHAR(45);" json:"course_type"`
	LearningOutcome string                        `gorm:"column:learning_outcome; default:NULL; type:LONGTEXT;" json:"learning_outcome"`
	Organization    OrganizationInCourses         `gorm:"foreignKey:OrganizationID;references:OrganizationID" json:"organization"`
	OrganizationID  int                           `gorm:"column:organization_id" json:"organization_id"`
	SkillsLevels    []SkillsLevelsInCoursesUpdate `gorm:"many2many:courses_skills_levels;foreignKey:CourseID;joinForeignKey:CourseID;References:SkillsLevelsID;joinReferences:SkillsLevelsID" json:"skills_levels"`
}

type SkillsLevelsInCoursesUpdate struct {
	SkillsLevelsID int    `gorm:"column:skills_levels_id; primaryKey; autoIncrement;" json:"skills_levels_id"`
	SkillID        *int   `gorm:"column:skill_id;" json:"skill_id"`
	KnowledgeDesc  string `gorm:"column:knowledge_desc;" json:"knowledge_desc"`
	AbilityDesc    string `gorm:"column:ability_desc;" json:"ability_desc"`
	LevelID        int    `gorm:"column:level_id; not null" json:"level_id"`
}

func (Course) TableName() string {
	return "courses"
}

func (UpdateCourseDTO) TableName() string {
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

func (SkillsLevelsInCoursesUpdate) TableName() string {
	return "skills_levels"
}

func (SkillInCourses) TableName() string {
	return "skills"
}

func (CareerInCourses) TableName() string {
	return "careers"
}

// GetCourses retrieves all Course records from the database with pagination.
func GetCourses(db *gorm.DB, page, limit int, courses *[]Course) (pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Careers").Preload("SkillsLevels").Model(&Course{}).
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

// GetAllCoursesWithFilter retrieves Course records from the database with filtering and pagination.
func GetAllCoursesWithFilter(db *gorm.DB, page, limit int, search string) (courses []Course, pagination Pagination, err error) {
	offset := (page - 1) * limit

	// Create a query builder with preloads and filters
	query := db.Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Careers").Preload("SkillsLevels").Offset(offset).Limit(limit)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	err = query.Find(&courses).Error
	if err != nil {
		return nil, Pagination{}, err
	}

	// Calculate total pages
	var totalCount int64
	if err := query.Model(&Course{}).Count(&totalCount).Error; err != nil {
		return nil, Pagination{}, err
	}

	totalPages := int(totalCount)

	pagination = Pagination{
		Page:  page,
		Total: totalPages,
		Limit: limit,
	}

	return courses, pagination, nil
}

// GetCoursesBySkillId retrieves courses based on the provided SkillID with pagination.
func GetCoursesBySkillId(db *gorm.DB, skillID, page, limit int, courses *[]Course) (pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Careers").Preload("SkillsLevels").
		Joins("JOIN courses_skills_levels ON courses.course_id = courses_skills_levels.course_id").
		Joins("JOIN skills_levels ON courses_skills_levels.skills_levels_id = skills_levels.skills_levels_id").
		Where("skills_levels.skill_id = ?", skillID).
		Distinct().
		Model(&Course{}).
		Offset(offset).Limit(limit).
		Find(courses).Error
	if err != nil {
		return Pagination{}, err
	}

	// Count total courses for pagination
	var totalCount int64
	if err := db.Model(&Course{}).Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Careers").Preload("SkillsLevels").
		Joins("JOIN courses_skills_levels ON courses.course_id = courses_skills_levels.course_id").
		Joins("JOIN skills_levels ON courses_skills_levels.skills_levels_id = skills_levels.skills_levels_id").
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
	err = db.Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Careers").Preload("SkillsLevels").
		Joins("JOIN courses_skills_levels ON courses.course_id = courses_skills_levels.course_id").
		Joins("JOIN careers_skills_levels ON careers_skills_levels.skills_levels_id = courses_skills_levels.skills_levels_id").
		Where("careers_skills_levels.career_id = ?", careerID).
		Distinct().
		Model(&CourseWithoutSkillLevels{}).
		Offset(offset).Limit(limit).
		Find(courses).Error
	if err != nil {
		return Pagination{}, err
	}

	// Count total courses for pagination
	var totalCount int64
	if err := db.Model(&CourseWithoutSkillLevels{}).Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Careers").Preload("SkillsLevels").
		Joins("JOIN courses_skills_levels ON courses.course_id = courses_skills_levels.course_id").
		Joins("JOIN careers_skills_levels ON careers_skills_levels.skills_levels_id = courses_skills_levels.skills_levels_id").
		Where("careers_skills_levels.career_id = ?", careerID).
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
	err = db.Where("course_id = ?", id).Preload("Organization").Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Careers").Preload("SkillsLevels").
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
	// Check if the course record exists
	existingCourse := &Course{}
	err = tx.Preload("Organization").Preload("SkillsLevels.Skill").Preload("SkillsLevels.Careers").Preload("SkillsLevels").First(existingCourse, "course_id = ?", updatedCourse.CourseID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Save the updated course
	err = tx.Save(updatedCourse).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Clear existing associations within the transaction
	err = tx.Model(existingCourse).Association("SkillsLevels").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}

	// // Clear existing associations within the transaction
	// err = tx.Model(existingCourse).Association("SkillsLevels.Careers").Clear()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// Update existing skillslevels with the new one (if provided)
	if len(updatedCourse.SkillsLevels) > 0 {
		err = tx.Model(existingCourse).Association("SkillsLevels").Append(updatedCourse.SkillsLevels)
		if err != nil {
			tx.Rollback()
			return err
		}
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
