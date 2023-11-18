package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Courses struct {
	CourseID       int    `gorm:"column:course_id;primaryKey"`
	Label          string `gorm:"column:label;not null"`
	Value          string `gorm:"column:value;not null"`
	Description    string `gorm:"column:description;default:NULL"`
	CourseLevel    string `gorm:"column:course_level;default:NULL"`
	LearnHours     string `gorm:"column:learn_hours;default:NULL"`
	AcademicYear   string `gorm:"column:academic_year;default:NULL"`
	CourseLink     string `gorm:"column:course_link;default:NULL"`
	OrganizationID int
	Organizations  Organization `gorm:"foreignKey:OrganizationID;references:OrganizationID"` // optional you can use organization_id in references
}

func (Courses) TableName() string {
	return "courses"
}

type Organization struct {
	OrganizationID int    `gorm:"column:organization_id;primaryKey"`
	Name           string `gorm:"column:name;not null"`
	Description    string `gorm:"column:description;default:NULL"`
	ImageUrl       string `gorm:"column:image_url;default:NULL"`
}

func (Organization) TableName() string {
	return "organizations"
}

func GetCourses(db *gorm.DB, Courses *[]Courses) (err error) {
	fmt.Println(db.Preload("Organizations").Find(Courses).Statement)
	err = db.Preload("Organizations").Find(Courses).Error
	if err != nil {
		return err
	}
	return nil
}

func GetCourseById(db *gorm.DB, Course *Courses, id int) (err error) {
	err = db.Where("Course_id = ?", id).Preload("Organizations").First(Course).Error
	if err != nil {
		return err
	}
	return nil
}
