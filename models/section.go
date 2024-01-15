package models

import (
	"gorm.io/gorm"
)

type Section struct {
	SectionID int64    `gorm:"column:section_id; primaryKey;autoIncrement" json:"section_id"`
	Name      string   `gorm:"column:name; not null; type:VARCHAR(255);" json:"name" binding:"required,max=255"`
	ImageUrl  string   `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
	Groups    []Groups `gorm:"many2many:sections_groups;foreignKey:SectionID;joinForeignKey:SectionID;References:GroupID;joinReferences:GroupID" json:"-"`
}

type Groups struct {
	GroupID int64  `gorm:"column:group_id;primaryKey;autoIncrement;" json:"group_id"`
	Name    string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
}

func (Section) TableName() string {
	return "sections"
}

type UpdateSectionModels struct {
	SectionID int64  `gorm:"column:section_id; primaryKey;autoIncrement" json:"section_id"`
	Name      string `gorm:"column:name; not null; type:VARCHAR(255);" json:"name" binding:"max=255"`
	ImageUrl  string `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
}

// Getsections retrieves all section records from the database.
func GetSections(db *gorm.DB, section *[]Section) (err error) {
	err = db.Preload("groups").Find(section).Error
	if err != nil {
		return err
	}
	return nil
}

// GetsectionById retrieves a section by its ID from the database.
func GetsectionById(db *gorm.DB, section *Section, id int) (err error) {
	err = db.Where("section_id = ?", id).Preload("groups").First(section).Error
	if err != nil {
		return err
	}
	return nil
}

// Creategroup creates a new group record in the database.
func Createsection(db *gorm.DB, section *Section) (err error) {
	err = db.Create(section).Error
	if err != nil {
		return err
	}

	return nil
}

// Updatesection updates a section record by ID.
func Updatesection(db *gorm.DB, section *Section) (err error) {
	err = db.Save(section).Error
	if err != nil {
		return err
	}
	return nil
}

// DeletesectionById deletes a section record by its ID from the database.
func DeletesectionById(db *gorm.DB, id int) (err error) {
	err = db.Where("section_id = ?", id).Delete(&Section{}).Error
	if err != nil {
		return err
	}
	return nil
}
