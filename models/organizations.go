package models

import (
	"gorm.io/gorm"
)

type Organizations struct {
	OrganizationID int    `gorm:"column:organization_id; primaryKey;" json:"organization_id"`
	Name           string `gorm:"column:name; not null; type:VARCHAR(255)" json:"name" binding:"required,max=255"`
	Description    string `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl       string `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
}

type UpdateOrganizationModels struct {
	Name        string `gorm:"column:name; type:VARCHAR(255);" json:"name" binding:"max=255"`
	Description string `gorm:"column:description; default:NULL; type:LONGTEXT;" json:"description" binding:"max=1500"`
	ImageUrl    string `gorm:"column:image_url; default:NULL; type:LONGTEXT;" json:"image_url" binding:"max=5000"`
}

func (Organizations) TableName() string {
	return "organizations"
}

// GetOrganizations retrieves all Organization records from the database.
func GetOrganizations(db *gorm.DB, Organizations *[]Organizations) (err error) {
	err = db.Find(Organizations).Error
	if err != nil {
		return err
	}
	return nil
}

// GetOrganizationById retrieves an Organization by its ID from the database.
func GetOrganizationById(db *gorm.DB, Organization *Organizations, id int) (err error) {
	err = db.Where("organization_id = ?", id).First(Organization).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateOrganization creates a new Organization record in the database.
func CreateOrganization(db *gorm.DB, organization *Organizations) (err error) {
	err = db.Create(organization).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrganization updates an existing Organization record in the database.
func UpdateOrganization(db *gorm.DB, id int, organization *UpdateOrganizationModels) (err error) {
	err = db.Model(&Organizations{}).Where("organization_id = ?", id).Updates(organization).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteOrganization deletes an Organization record from the database by its ID.
func DeleteOrganization(db *gorm.DB, id int) (err error) {
	err = db.Where("organization_id = ?", id).Delete(&Organizations{}).Error
	if err != nil {
		return err
	}
	return nil
}
