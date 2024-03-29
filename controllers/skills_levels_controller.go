package controllers

import (
	"knowledgeplus/go-api/database"
	"knowledgeplus/go-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SkillsLevelsRepo struct {
	Db *gorm.DB
}

func NewSkillsLevelsRepo() *SkillsLevelsRepo {
	db := database.InitDb()
	return &SkillsLevelsRepo{Db: db}
}

// GetAllSkillsLevels retrieves all SkillsLevels records from the database with pagination and optional skill name search.
func (repository *SkillsLevelsRepo) GetAllSkillsLevels(c *gin.Context) {
	var (
		skillsLevels []models.SkillsLevels
		pagination   models.Pagination
		err          error
	)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // set a default limit
	}

	skillName := c.Query("search") // Retrieve skill name search parameter

	pagination, err = models.GetAllSkillsLevels(repository.Db, page, limit, skillName, &skillsLevels) // Pass skillName to the function
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skills_levels": skillsLevels,
		"pagination":    pagination,
	})
}
