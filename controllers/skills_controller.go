package controllers

import (
	"errors"
	"knowledgeplus/go-api/database"
	"knowledgeplus/go-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SkillsRepo struct {
	Db *gorm.DB
}

func NewSkillsRepo() *SkillsRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Skills{},&models.Levels{})
	return &SkillsRepo{Db: db}
}

//get Skillss
func (repository *SkillsRepo) GetSkills(c *gin.Context) {
	var Skills []models.Skills
	print(Skills)
	err := models.GetSkills(repository.Db, &Skills)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Skills)
}

//get Skills by id
func (repository *SkillsRepo) GetSkillById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var Skills models.Skills
	err := models.GetSkillById(repository.Db, &Skills, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Skills)
}

// CreateSkill creates a new Skill record.
func (repository *SkillsRepo) CreateSkill(c *gin.Context) {
	var Skill models.Skills
	// var CategoriesID = Skill.Categories[0].CategoryID
	if err := c.ShouldBindJSON(&Skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println(&Skill)
	err := models.CreateSkill(repository.Db, &Skill)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// fmt.Println(Skill.SkillID)

	c.JSON(http.StatusCreated, Skill)
}
