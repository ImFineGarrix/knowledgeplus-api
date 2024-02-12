package models

import (
	"gorm.io/gorm"
)

type Level struct {
	LevelID int    `gorm:"column:level_id; primaryKey;" json:"level_id"`
	Level   string `gorm:"column:level; not null;" json:"level"`
}

func (Level) TableName() string {
	return "levels"
}

// GetLevels retrieves all Level records from the database.
func GetLevels(db *gorm.DB, Levels *[]Level) (err error) {
	err = db.Find(Levels).Error
	if err != nil {
		return err
	}
	return nil
}

// GetLevelById retrieves a Level by its ID from the database.
func GetLevelById(db *gorm.DB, Level *Level, id int) (err error) {
	err = db.Where("Level_id = ?", id).First(Level).Error
	if err != nil {
		return err
	}
	return nil
}
