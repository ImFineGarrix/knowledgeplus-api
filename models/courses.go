package models

import (
	"gorm.io/gorm"
)

type Course struct {
	CourseID       int                   `gorm:"column:course_id;primaryKey" json:"course_id"`
	Name           string                `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description    string                `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	LearnHours     string                `gorm:"column:learn_hours; default:NULL; type:VARCHAR(45);" json:"learn_hours"`
	AcademicYear   string                `gorm:"column:academic_year; default:NULL; type:VARCHAR(45);" json:"academic_year"`
	CourseLink     string                `gorm:"column:course_link; default:NULL; type:LONGTEXT;" json:"course_link"`
	Organization   OrganizationInCourses `gorm:"foreignKey:OrganizationID;references:OrganizationID" json:"organization"`
	OrganizationID int                   `gorm:"column:organization_id" json:"-"`
	SkillsLevels   []SkillsLevels        `gorm:"foreignKey:CourseID; References:CourseID;" json:"skills_levels"`
}

type OrganizationInCourses struct {
	OrganizationID int    `gorm:"column:organization_id; primaryKey;" json:"organization_id"`
	Name           string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"max=255"`
	Description    string `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl       string `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
}

func (Course) TableName() string {
	return "courses"
}

func (OrganizationInCourses) TableName() string {
	return "organizations"
}

// GetCourses retrieves all Course records from the database with pagination.
func GetCourses(db *gorm.DB, page, limit int, courses *[]Course) (pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Preload("Organization").Preload("SkillsLevels").Model(&Course{}).
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
	err = db.Where("course_id = ?", id).
		Preload("Organization").
		Preload("SkillsLevels").
		First(course).
		Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateCourse updates an existing Course record in the database.
func UpdateCourse(db *gorm.DB, updatedCourse *Course) (err error) {
	existingCourse := &Course{}

	// Begin a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the course record exists
	err = tx.Preload("Organization").Preload("SkillsLevels").First(existingCourse, "course_id = ?", updatedCourse.CourseID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update only the specified fields if they are not empty
	existingCourse.Name = updatedCourse.Name
	existingCourse.Description = updatedCourse.Description
	existingCourse.LearnHours = updatedCourse.LearnHours
	existingCourse.AcademicYear = updatedCourse.AcademicYear
	existingCourse.CourseLink = updatedCourse.CourseLink

	// Update organization_id only
	existingCourse.OrganizationID = updatedCourse.Organization.OrganizationID

	db.Save(existingCourse)

	// Clear existing associations within the transaction
	err = tx.Model(existingCourse).Association("SkillsLevels").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update existing skills_levels with the new ones (if provided)
	if len(updatedCourse.SkillsLevels) > 0 {
		err = tx.Model(existingCourse).Association("SkillsLevels").Append(updatedCourse.SkillsLevels)
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

// DeleteCourseById deletes a Course record by its ID from the database.
func DeleteCourseById(db *gorm.DB, id int) (err error) {
	err = db.Where("course_id = ?", id).Delete(&Course{}).Error
	if err != nil {
		return err
	}
	return nil
}
