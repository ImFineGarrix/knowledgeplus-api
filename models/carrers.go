package models

import (
	"gorm.io/gorm"
)

type Career struct {
	CareerID    int64             `gorm:"column:career_id;primaryKey;autoIncrement;" json:"career_id"`
	Name        string            `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
	Description string            `gorm:"column:description; default:NULL; type:LONGTEXT;"  json:"description" binding:"max=1500"`
	Groups      []GroupsInCareers `gorm:"many2many:groups_careers;foreignKey:CareerID;joinForeignKey:CareerID;References:GroupID;joinReferences:GroupID" json:"groups"`
	Skills      []Skills          `gorm:"many2many:careers_skills;foreignKey:CareerID;joinForeignKey:CareerID;References:SkillID;joinReferences:SkillID" json:"-"`
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

type Skills struct {
	SkillID     int    `gorm:"column:skill_id;primaryKey" json:"skill_id"`
	Name        string `gorm:"column:name;not null; type:VARCHAR(255);" json:"name" binding:"max=255"`
	Description string `gorm:"column:description;default:NULL type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl    string `gorm:"column:image_url;default:NULL type:LONGTEXT;" json:"image_url" binding:"max=5000"`
	Type        string `gorm:"column:type;type:ENUM('SOFT','HARD');" json:"type"`
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

func (Skills) Tablename() string {
	return "skills"
}

func (GroupsInCareers) TableName() string {
	return "groups"
}

// เอาแค่ group ไม่เอา skill
// GetCareers retrieves all Career records from the database with pagination.
func GetCareers(db *gorm.DB, page, limit int) (careers []Career, pagination Pagination, err error) {
	offset := (page - 1) * limit
	err = db.Preload("Groups").Preload("Skills").
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

//ยังไม่ใช้
// // GetCareersWithHaveCategories retrieves careers where category_id in many-to-many is not null with pagination.
// func GetCareersWithHaveCategories(db *gorm.DB, page, limit int) (careers []Career, pagination Pagination, err error) {
// 	offset := (page - 1) * limit
// 	err = db.Joins("JOIN categories_careers ON careers.career_id = categories_careers.career_id").
// 		Where("categories_careers.category_id IS NOT NULL").
// 		Preload("Categories").Preload("Skills").
// 		Offset(offset).Limit(limit).
// 		Find(&careers).Error
// 	if err != nil {
// 		return nil, Pagination{}, err
// 	}

// 	// Calculate total pages
// 	var totalCount int64
// 	if err := db.Joins("JOIN categories_careers ON careers.career_id = categories_careers.career_id").
// 		Where("categories_careers.category_id IS NOT NULL").
// 		Model(&Career{}).Count(&totalCount).Error; err != nil {
// 		return nil, Pagination{}, err
// 	}

// 	totalPages := int(totalCount)

// 	pagination = Pagination{
// 		Page:  page,
// 		Total: totalPages,
// 		Limit: limit,
// 	}

// 	return careers, pagination, nil
// }

// Get By id ให้โชว์ Skill อย่างเดียว(ทำหลัง Course เสร็จ)
// GetCareerById retrieves a Career by its ID from the database.
func GetCareerById(db *gorm.DB, Career *Career, id int) (err error) {
	err = db.Where("career_id = ?", id).Preload("Groups").Preload("Skills").First(Career).Error
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
	err = tx.Where("career_id = ?", updatedCareer.CareerID).Preload("Groups").Preload("Skills").First(existingCareer).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update only the specified fields if they are not empty
	if updatedCareer.Name != "" {
		existingCareer.Name = updatedCareer.Name
	}

	if updatedCareer.Description != "" {
		existingCareer.Description = updatedCareer.Description
	}

	db.Save(existingCareer)

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
