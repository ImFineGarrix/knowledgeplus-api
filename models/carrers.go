package models

import (
	"gorm.io/gorm"
)

// all function career not finish
type Career struct {
	CareerID     int64                   `gorm:"column:career_id;primaryKey;autoIncrement;" json:"career_id"`
	Name         string                  `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
	Description  string                  `gorm:"column:description; default:NULL; type:LONGTEXT;"  json:"description" binding:"max=1500"`
	Groups       []GroupsInCareers       `gorm:"many2many:groups_careers;foreignKey:CareerID;joinForeignKey:CareerID;References:GroupID;joinReferences:GroupID" json:"groups"`
	SkillsLevels []SkillsLevelsInCareers `gorm:"foreignKey:CareerID; References:CareerID;" json:"skills_levels"`
}

type GroupsInCareers struct {
	GroupID int64  `gorm:"column:group_id;primaryKey;autoIncrement;" json:"group_id"`
	Name    string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
	// Sections []SectionsInGroup `gorm:"many2many:sections_groups;foreignKey:GroupID;joinForeignKey:GroupID;References:SectionID;joinReferences:SectionID" json:"sections"`
}

type SectionInCareers struct {
	SectionID   int64  `gorm:"column:section_id; primaryKey; autoIncrement" json:"section_id"`
	Name        string `gorm:"column:name; not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description string `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
}

type SkillsLevelsInCareers struct {
	SkillsLevelsID int              `gorm:"column:skills_levels_id; primaryKey; autoIncrement;" json:"skills_levels_id"`
	SkillID        *int             `gorm:"column:skill_id;" json:"skill_id"`
	KnowledgeDesc  string           `gorm:"column:knowledge_desc;" json:"knowledge_desc"`
	AbilityDesc    string           `gorm:"column:ability_desc;" json:"ability_desc"`
	LevelID        int              `gorm:"column:level_id; not null" json:"level_id"`
	CourseID       *int             `gorm:"column:course_id;" json:"-"`
	CareerID       *int             `gorm:"column:career_id;" json:"-"`
	Skill          SkillInCareers   `gorm:"foreignKey:SkillID;references:SkillID" json:"skill"`
	Course         *CourseInCareers `gorm:"foreignKey:CourseID;references:CourseID" json:"courses"`
}

type SkillInCareers struct {
	SkillID      *int                   `gorm:"column:skill_id;primaryKey" json:"-"`
	Name         string                 `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description  string                 `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl     string                 `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
	Type         string                 `gorm:"column:type; default:NULL; type:ENUM('SOFT','HARD');" json:"type" binding:"max=100"`
	SkillsLevels []SkillsLevelsInSkills `gorm:"foreignKey:SkillID; References:SkillID;" json:"-"`
}

type CourseInCareers struct {
	CourseID        int                     `gorm:"column:course_id;primaryKey" json:"-"`
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

// type UpdateCareerModels struct {
// 	CareerID    int64        `gorm:"column:career_id;primaryKey;autoIncrement;" json:"career_id"`
// 	Name        string       `gorm:"column:name; not null; type:VARCHAR(45)" json:"name" binding:"max=45"`
// 	Description string       `gorm:"column:description; default:NULL; type:LONGTEXT;"  json:"description" binding:"max=1500"`
// 	Categories  []Categories `gorm:"many2many:categories_careers;foreignKey:CareerID;joinForeignKey:CareerID;References:CategoryID;joinReferences:CategoryID" json:"categories"`
// 	Skills      []Skills     `gorm:"many2many:careers_skills;foreignKey:CareerID;joinForeignKey:CareerID;References:SkillID;joinReferences:SkillID" json:"skills"`
// }

func (Career) TableName() string {
	return "careers"
}

func (SkillsLevelsInCareers) TableName() string {
	return "skills_levels"
}

func (CourseInCareers) TableName() string {
	return "courses"
}

func (SkillInCareers) TableName() string {
	return "skills"
}

func (GroupsInCareers) TableName() string {
	return "groups"
}

// GetCareers retrieves all Career records from the database with pagination.
func GetCareers(db *gorm.DB, page, limit int) (careers []Career, pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Preload("SkillsLevels.Skill").Preload("SkillsLevels.Course.Organization").Preload("Groups").
		Offset(offset).Limit(limit).
		Find(&careers).Error
	if err != nil {
		return nil, Pagination{}, err
	}

	// Calculate total pages
	var totalCount int64
	if err := db.Model(&Career{}).Count(&totalCount).Error; err != nil {
		return nil, Pagination{}, err
	}

	totalPages := int(totalCount)

	pagination = Pagination{
		Page:  page,
		Total: totalPages,
		Limit: limit,
	}

	return careers, pagination, nil
}

// GetCareersByCourseId retrieves careers based on the provided CourseID with pagination.
func GetCareersByCourseId(db *gorm.DB, courseID, page, limit int) (careers []Career, pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.
		Preload("SkillsLevels.Skill").
		Preload("Groups").
		Preload("SkillsLevels.Course.Organization").
		Joins("JOIN skills_levels ON careers.career_id = skills_levels.career_id").
		Where("skills_levels.course_id = ?", courseID).
		Group("careers.career_id").
		Distinct().
		Offset(offset).Limit(limit).
		Find(&careers).Error
	if err != nil {
		return nil, Pagination{}, err
	}

	// Count total careers for pagination
	var totalCount int64
	if err := db.Model(&Career{}).
		Joins("JOIN skills_levels ON careers.career_id = skills_levels.career_id").
		Where("skills_levels.course_id = ?", courseID).
		Group("careers.career_id").
		Group("careers.career_id").
		Distinct().
		Count(&totalCount).Error; err != nil {
		return nil, Pagination{}, err
	}

	totalPages := int(totalCount)

	pagination = Pagination{
		Page:  page,
		Total: totalPages,
		Limit: limit,
	}

	return careers, pagination, nil
}

// GetCareersWithFilters retrieves Career records from the database with filtering and pagination. for use frontend
func GetCareersWithFilters(db *gorm.DB, page, limit int, search string, groupID int64) (careers []Career, pagination Pagination, err error) {
	offset := (page - 1) * limit

	// Create a query builder with preloads and filters
	query := db.Preload("SkillsLevels.Skill").Preload("SkillsLevels.Course.Organization").Preload("Groups").
		Offset(offset).Limit(limit)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if groupID != 0 {
		query = query.Joins("JOIN groups_careers ON careers.career_id = groups_careers.career_id").
			Where("groups_careers.group_id = ?", groupID)
	}

	err = query.Find(&careers).Error
	if err != nil {
		return nil, Pagination{}, err
	}

	// Calculate total pages
	var totalCount int64
	if err := query.Model(&Career{}).Count(&totalCount).Error; err != nil {
		return nil, Pagination{}, err
	}

	totalPages := int(totalCount)

	pagination = Pagination{
		Page:  page,
		Total: totalPages,
		Limit: limit,
	}

	return careers, pagination, nil
}

// GetCareersBySkillId retrieves careers based on the provided SkillID with pagination.
func GetCareersBySkillId(db *gorm.DB, skillID, page, limit int) (careers []Career, pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.
		Preload("SkillsLevels.Skill").
		Preload("Groups").
		Preload("SkillsLevels.Course.Organization").
		Joins("JOIN skills_levels ON careers.career_id = skills_levels.career_id").
		Where("skills_levels.skill_id = ?", skillID).
		Distinct().
		Offset(offset).Limit(limit).
		Find(&careers).Error
	if err != nil {
		return nil, Pagination{}, err
	}

	// Count total careers for pagination
	var totalCount int64
	if err := db.Model(&Career{}).
		Joins("JOIN skills_levels ON careers.career_id = skills_levels.career_id").
		Where("skills_levels.skill_id = ?", skillID).
		Count(&totalCount).Error; err != nil {
		return nil, Pagination{}, err
	}

	totalPages := int(totalCount)

	pagination = Pagination{
		Page:  page,
		Total: totalPages,
		Limit: limit,
	}

	return careers, pagination, nil
}

// GetCareerById retrieves a Career by its ID from the database.
func GetCareerById(db *gorm.DB, Career *Career, id int) (err error) {
	err = db.Where("career_id = ?", id).Preload("SkillsLevels.Skill").Preload("SkillsLevels.Course.Organization").Preload("Groups").First(Career).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateCareer creates a new Career record in the database.
func CreateCareer(db *gorm.DB, career *Career) (err error) {
	// err = db.Omit("Categories").Create(career).Error
	err = db.Create(career).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateCareer updates an existing Career record in the database. ** แก้รายละเอียดใน group ไม่ได้(ระบบตอนแก้ดู id เป็นหลัก,ถ้าไม่ได้ส่ง id แล้วทีรายละเอียดอื่นๆส่งมามันจะส้าง id ใหม่ใน group ให้)
func UpdateCareer(db *gorm.DB, updatedCareer *Career) (err error) {
	existingCareer := &Career{}

	// Begin a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the career record exists
	err = tx.Where("career_id = ?", updatedCareer.CareerID).Preload("SkillsLevels.Skill").Preload("SkillsLevels.Course.Organization").First(existingCareer).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete all existing SkillsLevels records for the specified course
	err = tx.Where("career_id = ?", existingCareer.CareerID).Delete(&SkillsLevels{}).Error
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
	err = tx.Save(updatedCareer).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Clear existing associations within the transaction
	err = tx.Model(existingCareer).Association("Groups").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update existing groups with the new one (if provided)
	if len(updatedCareer.Groups) > 0 {
		err = tx.Model(existingCareer).Association("Groups").Append(updatedCareer.Groups)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteCareerById deletes a Career by its ID from the database.
func DeleteCareerById(db *gorm.DB, id int) (err error) {
	err = db.Where("career_id = ?", id).Delete(&Career{}).Error
	if err != nil {
		return err
	}
	return nil
}

// Pagination struct
type Pagination struct {
	Page  int `json:"page"`
	Total int `json:"total"`
	Limit int `json:"limit"`
}

// get by id ที่เอาแค่ course อย่างเดียว (ทำเป็นเส้นใหม่)(ทำหลัง course เสร็จ)
