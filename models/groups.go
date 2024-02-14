package models

import (
	"gorm.io/gorm"
)

type Group struct {
	GroupID  int64             `gorm:"column:group_id; primaryKey;autoIncrement" json:"group_id"`
	Name     string            `gorm:"column:name; not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Sections []SectionsInGroup `gorm:"many2many:sections_groups;foreignKey:GroupID;joinForeignKey:GroupID;References:SectionID;joinReferences:SectionID" json:"sections"`
	Careers  []CareersInGroup  `gorm:"many2many:groups_careers;foreignKey:GroupID;joinForeignKey:GroupID;References:CareerID;joinReferences:CareerID" json:"-"`
}

type CareersInGroup struct {
	CareerID int64  `gorm:"column:career_id;primaryKey;autoIncrement;" json:"career_id"`
	Name     string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
}

type SectionsInGroup struct {
	SectionID int64  `gorm:"column:section_id;primaryKey;autoIncrement;" json:"section_id"`
	Name      string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
}

func (Group) TableName() string {
	return "groups"
}

func (CareersInGroup) TableName() string {
	return "careers"
}

func (SectionsInGroup) TableName() string {
	return "sections"
}

type UpdateGroupModels struct {
	GroupID  int64             `gorm:"column:group_id; primaryKey;autoIncrement" json:"group_id"`
	Name     string            `gorm:"column:name; not null; type:VARCHAR(255);" json:"name" binding:"max=255"`
	Sections []SectionsInGroup `gorm:"many2many:sections_groups;foreignKey:GroupID;joinForeignKey:GroupID;References:SectionID;joinReferences:SectionID" json:"sections"`
	Careers  []CareersInGroup  `gorm:"many2many:groups_careers;foreignKey:GroupID;joinForeignKey:GroupID;References:CareerID;joinReferences:CareerID" json:"careers"`
}

// Getgroups retrieves all group records from the database.
func Getgroups(db *gorm.DB, group *[]Group) (err error) {
	err = db.Preload("Sections").Preload("Careers").Find(group).Error
	if err != nil {
		return err
	}
	return nil
}

// Getgroups retrieves all group records from the database with optional filtering based on sections.
func GetgroupsHaveSection(db *gorm.DB, group *[]Group, filterGroupsBySections bool) (err error) {
	query := db.Preload("Sections").Preload("Careers")
	query = query.Joins("JOIN sections_groups ON sections_groups.group_id = groups.group_id").
		Group("groups.group_id").
		Having("COUNT(sections_groups.section_id) > 0")
	err = query.Find(group).Error
	if err != nil {
		return err
	}
	return nil
}

// GetgroupById retrieves a group by its ID from the database.
func GetgroupById(db *gorm.DB, group *Group, id int) (err error) {
	err = db.Where("group_id = ?", id).Preload("Sections").Preload("Careers").First(group).Error
	if err != nil {
		return err
	}
	return nil
}

// Creategroup creates a new group record in the database.
func Creategroup(db *gorm.DB, group *Group) (err error) {
	err = db.Create(group).Error
	if err != nil {
		return err
	}

	return nil
}

// // Updategroup updates a group record by ID.
// func Updategroup(db *gorm.DB, updatedGroup *UpdateGroupModels) (err error) {
// 	fmt.Println("1")
// 	// Find the existing group by ID
// 	existingGroup := Group{}
// 	if err := db.First(&existingGroup, updatedGroup.GroupID).Error; err != nil {
// 		return err
// 	}
// 	fmt.Println("2")
// 	// Begin a transaction
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()
// 	fmt.Println("2.1")
// 	// Update only the specified fields if they are not empty
// 	if updatedGroup.Name != "" {
// 		existingGroup.Name = updatedGroup.Name
// 	}
// 	fmt.Println("2.2")
// 	// // Save the changes
// 	db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&existingGroup)
// 	fmt.Println("2.3")
// 	// Clear existing associations within the transaction
// 	err = tx.Model(existingGroup).Association("Sections").Clear()
// 	if err != nil {
// 		fmt.Println("this is error")
// 		tx.Rollback()
// 		return err
// 	}
// 	fmt.Println("3")
// 	err = tx.Model(existingGroup).Association("Careers").Clear()
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	fmt.Println("4")
// 	// Update existing categories with the new one (if provided)
// 	if len(updatedGroup.Sections) > 0 {
// 		err = tx.Model(existingGroup).Association("Sections").Append(updatedGroup.Sections)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}
// 	fmt.Println("5")
// 	// Update existing categories with the new one (if provided)
// 	if len(updatedGroup.Careers) > 0 {
// 		err = tx.Model(existingGroup).Association("Careers").Append(updatedGroup.Careers)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}
// 	fmt.Println("6")
// 	// Commit the transaction
// 	err = tx.Commit().Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// Updategroup updates a group record by ID.
func Updategroup(db *gorm.DB, updatedGroup *UpdateGroupModels) (err error) {
	// Find the existing group by ID
	existingGroup := Group{}
	if err := db.First(&existingGroup, updatedGroup.GroupID).Error; err != nil {
		return err
	}

	// Begin a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update only the specified fields if they are not empty
	if updatedGroup.Name != "" {
		existingGroup.Name = updatedGroup.Name
	}

	// Save the changes
	db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&existingGroup)

	// Clear existing associations within the transaction
	err = tx.Model(&existingGroup).Association("Sections").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&existingGroup).Association("Careers").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update existing sections with the new ones (if provided)
	if len(updatedGroup.Sections) > 0 {
		err = tx.Model(&existingGroup).Association("Sections").Replace(updatedGroup.Sections)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update existing careers with the new ones (if provided)
	if len(updatedGroup.Careers) > 0 {
		err = tx.Model(&existingGroup).Association("Careers").Replace(updatedGroup.Careers)
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

// // Updatesection updates a section record by ID.
// func Updategroup(db *gorm.DB, updatedGroup *UpdateGroupModels) (err error) {
// 	// Find the existing group by ID
// 	existingGroup := Group{}
// 	if err := db.First(&existingGroup, updatedGroup.GroupID).Error; err != nil {
// 		return err
// 	}

// 	// Update only the specified fields if they are not empty
// 	if updatedGroup.Name != "" {
// 		existingGroup.Name = updatedGroup.Name
// 	}

// 	// // Save the changes
// 	// if err := db.Save(&existingGroup).Error; err != nil {
// 	// 	return err
// 	// }
// 	db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&existingGroup)

// 	return nil
// }

// DeletegroupById deletes a group record by its ID from the database.
func DeletegroupById(db *gorm.DB, id int) (err error) {
	err = db.Where("group_id = ?", id).Delete(&Group{}).Error
	if err != nil {
		return err
	}
	return nil
}
