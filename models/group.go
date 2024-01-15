package models

import (
	"gorm.io/gorm"
)

type Group struct {
	GroupID int64            `gorm:"column:group_id; primaryKey;autoIncrement" json:"group_id"`
	Name    string           `gorm:"column:name; not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	Careers []CareersInGroup `gorm:"many2many:groups_careers;foreignKey:GroupID;joinForeignKey:GroupID;References:CareerID;joinReferences:CareerID" json:"careers"`
}

type CareersInGroup struct {
	CareerID int64  `gorm:"column:career_id;primaryKey;autoIncrement;" json:"career_id"`
	Name     string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
}

func (Group) TableName() string {
	return "groups"
}

type UpdateGroupModels struct {
	GroupID int64  `gorm:"column:group_id; primaryKey;autoIncrement" json:"group_id"`
	Name    string `gorm:"column:name; not null; type:VARCHAR(255);" json:"name" binding:"max=255"`
}

// Getgroups retrieves all group records from the database.
func Getgroups(db *gorm.DB, group *[]Group) (err error) {
	err = db.Preload("groups").Preload("careers").Find(group).Error
	if err != nil {
		return err
	}
	return nil
}

// GetgroupById retrieves a group by its ID from the database.
func GetgroupById(db *gorm.DB, group *Group, id int) (err error) {
	err = db.Where("group_id = ?", id).Preload("careers").First(group).Error
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

// Updategroup updates a group record by ID.
func Updategroup(db *gorm.DB, group *Group) (err error) {
	err = db.Save(group).Error
	if err != nil {
		return err
	}
	return nil
}

// DeletegroupById deletes a group record by its ID from the database.
func DeletegroupById(db *gorm.DB, id int) (err error) {
	err = db.Where("group_id = ?", id).Delete(&Group{}).Error
	if err != nil {
		return err
	}
	return nil
}
