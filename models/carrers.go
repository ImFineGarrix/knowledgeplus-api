package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Career struct {
	CareerID     int64                   `gorm:"column:career_id;primaryKey;autoIncrement;" json:"career_id"`
	Name         string                  `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
	Description  string                  `gorm:"column:description; default:NULL; type:LONGTEXT;"  json:"description" binding:"max=1500"`
	Groups       []GroupsInCareers       `gorm:"many2many:groups_careers;foreignKey:CareerID;joinForeignKey:CareerID;References:GroupID;joinReferences:GroupID" json:"groups"`
	SkillsLevels []SkillsLevelsInCareers `gorm:"many2many:careers_skills_levels;foreignKey:CareerID;joinForeignKey:CareerID;References:SkillsLevelsID;joinReferences:SkillsLevelsID" json:"skills_levels"`
}

type GroupsInCareers struct {
	GroupID int64  `gorm:"column:group_id;primaryKey;autoIncrement;" json:"group_id"`
	Name    string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
	// Sections []SectionsInGroup `gorm:"many2many:sections_groups;foreignKey:GroupID;joinForeignKey:GroupID;References:SectionID;joinReferences:SectionID" json:"sections"`
}

type SkillsLevelsInCareers struct {
	SkillsLevelsID int               `gorm:"column:skills_levels_id; primaryKey; autoIncrement;" json:"skills_levels_id"`
	SkillID        *int              `gorm:"column:skill_id;" json:"skill_id"`
	KnowledgeDesc  string            `gorm:"column:knowledge_desc;" json:"knowledge_desc"`
	AbilityDesc    string            `gorm:"column:ability_desc;" json:"ability_desc"`
	LevelID        int               `gorm:"column:level_id; not null" json:"level_id"`
	Skill          SkillInCareers    `gorm:"foreignKey:SkillID;references:SkillID" json:"skill"`
	Courses        []CourseInCareers `gorm:"many2many:courses_skills_levels;foreignKey:SkillsLevelsID;joinForeignKey:SkillsLevelsID;References:CourseID;joinReferences:CourseID" json:"courses"`
}

type SkillInCareers struct {
	SkillID      *int                   `gorm:"column:skill_id;primaryKey" json:"-"`
	Name         string                 `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description  string                 `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl     string                 `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=2000"`
	Type         string                 `gorm:"column:type; default:NULL; type:ENUM('SOFT','HARD');" json:"type" binding:"max=100"`
	SkillsLevels []SkillsLevelsInSkills `gorm:"foreignKey:SkillID; References:SkillID;" json:"-"`
}

type CourseInCareers struct {
	CourseID        int                  `gorm:"column:course_id;primaryKey" json:"course_id"`
	Name            string               `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Description     string               `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	LearnHours      string               `gorm:"column:learn_hours; default:NULL; type:VARCHAR(45);" json:"learn_hours"`
	AcademicYear    string               `gorm:"column:academic_year; default:NULL; type:VARCHAR(45);" json:"academic_year"`
	CourseLink      string               `gorm:"column:course_link; default:NULL; type:LONGTEXT;" json:"course_link"`
	CourseType      string               `gorm:"column:course_type; default:NULL; type:VARCHAR(45);" json:"course_type"`
	LearningOutcome string               `gorm:"column:learning_outcome; default:NULL; type:LONGTEXT;" json:"learning_outcome"`
	OrganizationID  int                  `gorm:"column:organization_id" json:"organization_id"`
	Organization    OrganizationInCareer `gorm:"foreignKey:OrganizationID;references:OrganizationID" json:"organization"`
	// SkillsLevels    []SkillsLevelsInCourses `gorm:"foreignKey:CourseID; References:CourseID;" json:"-"`
}

type OrganizationInCareer struct {
	OrganizationID int    `gorm:"column:organization_id; primaryKey;" json:"-"`
	Name           string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"max=255"`
	Description    string `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl       string `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=2000"`
}

type CareerForRecommendSkillsLevels struct {
	CurrentCareerID  int64 `json:"current_career_id"`
	UserSkillsLevels []int `json:"user_skills_levels"`
}

type ReturnRecommendSkillsLevels struct {
	DifferenceSkillsLevels []SkillsLevelsInCareers `json:"difference_skills_levels"`
}

// Error implements error.
func (ReturnRecommendSkillsLevels) Error() string {
	panic("unimplemented")
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

func (OrganizationInCareer) TableName() string {
	return "organizations"
}

// GetCareers retrieves all Career records from the database with pagination.
func GetCareers(db *gorm.DB, page, limit int) (careers []Career, pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Preload("SkillsLevels.Skill").Preload("SkillsLevels.Courses.Organization").Preload("Groups").
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

	// Calculate total pages
	totalPages := int(totalCount) / limit
	if int(totalCount)%limit != 0 {
		totalPages++
	}

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
	err = db.Preload("SkillsLevels.Skill").Preload("SkillsLevels.Courses.Organization").Preload("Groups").
		Joins("JOIN careers_skills_levels ON careers.career_id = careers_skills_levels.career_id").
		Joins("JOIN skills_levels ON careers_skills_levels.skills_levels_id = skills_levels.skills_levels_id").
		Joins("JOIN courses_skills_levels ON skills_levels.skills_levels_id = courses_skills_levels.skills_levels_id").
		Where("courses_skills_levels.course_id = ?", courseID).
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
		Joins("JOIN careers_skills_levels ON careers.career_id = careers_skills_levels.career_id").
		Joins("JOIN skills_levels ON careers_skills_levels.skills_levels_id = skills_levels.skills_levels_id").
		Joins("JOIN courses_skills_levels ON skills_levels.skills_levels_id = courses_skills_levels.skills_levels_id").
		Where("courses_skills_levels.course_id = ?", courseID).
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
	query := db.Preload("SkillsLevels.Skill").Preload("SkillsLevels.Courses.Organization").Preload("Groups").
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
	err = db.Preload("SkillsLevels.Skill").Preload("SkillsLevels.Courses.Organization").Preload("Groups").
		Joins("JOIN careers_skills_levels ON careers.career_id = careers_skills_levels.career_id").
		Joins("JOIN skills_levels ON careers_skills_levels.skills_levels_id = skills_levels.skills_levels_id").
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
		Joins("JOIN careers_skills_levels ON careers.career_id = careers_skills_levels.career_id").
		Joins("JOIN skills_levels ON careers_skills_levels.skills_levels_id = skills_levels.skills_levels_id").
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
	err = db.Where("career_id = ?", id).Preload("SkillsLevels.Skill").Preload("SkillsLevels.Courses.Organization").Preload("Groups").First(Career).Error
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
	err = tx.Where("career_id = ?", updatedCareer.CareerID).Preload("SkillsLevels.Skill").Preload("SkillsLevels.Courses.Organization").Preload("Groups").First(existingCareer).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// // Delete all existing SkillsLevels records for the specified course
	// err = tx.Where("career_id = ?", existingCareer.CareerID).Delete(&SkillsLevels{}).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// // Commit the delete operation
	// if err := tx.Commit().Error; err != nil {
	// 	return err
	// }

	// // Begin a new transaction for the update operation
	// tx = db.Begin()

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

	// Clear existing associations within the transaction
	err = tx.Model(existingCareer).Association("SkillsLevels").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update existing skillslevels with the new one (if provided)
	if len(updatedCareer.SkillsLevels) > 0 {
		err = tx.Model(existingCareer).Association("SkillsLevels").Append(updatedCareer.SkillsLevels)
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

// func RecommendSkillsLevelsByCareer(db *gorm.DB, currentUserSkills *CareerForRecommendSkillsLevels) (Error error) {
// 	currentSkills := currentUserSkills.UserSkillsLevels

// 	var career Career
// 	if err := db.Where("career_id = ?", currentUserSkills.CurrentCareerID).Preload("SkillsLevels.Skill").Preload("SkillsLevels.Courses.Organization").Preload("SkillsLevels").Preload("Groups").First(&career).Error; err != nil {
// 		return nil
// 	}

// 	var skillLevelsIDs []int
// 	for _, skillLevel := range career.SkillsLevels {
// 		skillLevelsIDs = append(skillLevelsIDs, skillLevel.SkillsLevelsID)
// 	}

// 	// Find the differences between currentSkills and skillLevelsIDs
// 	differences := diffArray(currentSkills, skillLevelsIDs)

// 	// Fetch full SkillsLevels data for the differences
// 	var returnSkillsLevels []SkillsLevelsInCareers
// 	if len(differences) > 0 {
// 		var skillsLevelsFromDB []SkillsLevelsInCareers
// 		if err := db.Where("skills_levels_id IN (?)", differences).Preload("Skill").Preload("Courses.Organization").Find(&skillsLevelsFromDB).Error; err != nil {
// 			return nil
// 		}
// 		returnSkillsLevels = skillsLevelsFromDB
// 	} else {
// 		returnSkillsLevels = make([]SkillsLevelsInCareers, 0)
// 	}

// 	var returnResult ReturnRecommendSkillsLevels
// 	returnResult.DifferenceSkillsLevels = returnSkillsLevels

// 	// Return the full SkillsLevels data for the differences
// 	return returnResult
// }

// func diffArray(arr1, arr2 []int) []int {
// 	diff := make([]int, 0)

// 	// Create a map to store values of arr1 for quick lookup
// 	arr1Map := make(map[int]bool)
// 	for _, num := range arr1 {
// 		arr1Map[num] = true
// 	}

// 	// Check each element of arr2 if it exists in arr1
// 	for _, num := range arr2 {
// 		if !arr1Map[num] {
// 			diff = append(diff, num)
// 		}
// 	}

// 	return diff
// }

func RecommendSkillsLevelsByCareer(db *gorm.DB, currentUserSkills *CareerForRecommendSkillsLevels) (result ReturnRecommendSkillsLevels, ReturnError error) {
	currentSkills := currentUserSkills.UserSkillsLevels

	var career Career
	if err := db.Where("career_id = ?", currentUserSkills.CurrentCareerID).Preload("SkillsLevels.Skill").Preload("SkillsLevels.Courses.Organization").Preload("SkillsLevels").Preload("Groups").First(&career).Error; err != nil {
		return ReturnRecommendSkillsLevels{}, err
	}

	var skillLevelsIDs []int
	for _, skillLevel := range career.SkillsLevels {
		skillLevelsIDs = append(skillLevelsIDs, skillLevel.SkillsLevelsID)
	}

	// Find the differences between currentSkills and skillLevelsIDs
	differences := diffArray(currentSkills, skillLevelsIDs)

	// Fetch full SkillsLevels data for the differences
	returnSkillsLevels := make([]SkillsLevelsInCareers, 0)
	if len(differences) > 0 {
		var skillsLevelsFromDB []SkillsLevelsInCareers
		if err := db.Where("skills_levels_id IN (?)", differences).Preload("Skill").Preload("Courses.Organization").Find(&skillsLevelsFromDB).Error; err != nil {
			return ReturnRecommendSkillsLevels{}, err
		}

		if len(skillsLevelsFromDB) == 0 {
			return ReturnRecommendSkillsLevels{returnSkillsLevels}, nil
		}
		// Filter skills_levels based on level_id
		for _, skillLevel := range skillsLevelsFromDB {
			for _, userSkillID := range currentUserSkills.UserSkillsLevels {
				if skillLevel.LevelID > userSkillID && skillLevel.LevelID <= 7 {
					returnSkillsLevels = append(returnSkillsLevels, skillLevel)
					break // Move to the next skillLevel
				}
			}
		}
	}

	fmt.Println(returnSkillsLevels)

	var returnResult ReturnRecommendSkillsLevels
	returnResult.DifferenceSkillsLevels = returnSkillsLevels
	fmt.Println(returnResult.DifferenceSkillsLevels)

	// Return the full SkillsLevels data for the differences
	return returnResult, nil
}

func diffArray(arr1, arr2 []int) []int {
	diff := make([]int, 0)

	// Create a map to store values of arr1 for quick lookup
	arr1Map := make(map[int]bool)
	for _, num := range arr1 {
		arr1Map[num] = true
	}

	// Check each element of arr2 if it exists in arr1
	for _, num := range arr2 {
		if !arr1Map[num] {
			diff = append(diff, num)
		}
	}

	return diff
}

// Pagination struct
type Pagination struct {
	Page  int `json:"page"`
	Total int `json:"total"`
	Limit int `json:"limit"`
}

// get by id ที่เอาแค่ course อย่างเดียว (ทำเป็นเส้นใหม่)(ทำหลัง course เสร็จ)
