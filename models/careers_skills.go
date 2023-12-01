package models

type CareerSkill struct {
	CareersSkillID int `gorm:"primaryKey;autoIncrement"`
	CareerID       int `gorm:"foreignKey:CareerID ;references:CareerID "`
	SkillID        int `gorm:"foreignKey:CategoryID ;references:CategoryID "`
}

func (CareerSkill) TableName() string {
	return "careers_skills"
}
