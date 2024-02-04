package controllers

import (
	"errors"
	"fmt"
	"knowledgeplus/go-api/database"
	"knowledgeplus/go-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SkillRepo struct {
	Db *gorm.DB
}

func NewSkillRepo() *SkillRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Skill{}, &models.Levels{})
	return &SkillRepo{Db: db}
}

// GetSkills retrieves all Skill records from the database.
func (repository *SkillRepo) GetSkills(c *gin.Context) {
	var (
		skills     []models.Skill
		pagination models.Pagination
	)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // set a default limit
	}

	pagination, err = models.GetSkills(repository.Db, page, limit, &skills)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skills":     skills,
		"pagination": pagination,
	})
}

// get Skills by id
func (repository *SkillRepo) GetSkillById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var Skills models.Skill
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
func (repository *SkillRepo) CreateSkill(c *gin.Context) {
	var Skill models.Skill
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

// UpdateSkill updates a Skill record by ID.
func (repository *SkillRepo) UpdateSkill(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedSkill models.Skill

	// Check if the skill record exists
	err := models.GetSkillById(repository.Db, &updatedSkill, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Bind the updated data from the request
	if err := c.ShouldBindJSON(&updatedSkill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the skill record
	err = models.UpdateSkill(repository.Db, &updatedSkill)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, updatedSkill)
}

// DeleteSkillById deletes a Skill record by ID.
func (repository *SkillRepo) DeleteSkillById(c *gin.Context) {
	fmt.Println("Hello, This is delete skill function")
	fmt.Println(c.Errors)
	id, _ := strconv.Atoi(c.Param("id"))
	var skill models.Skill

	err := repository.Db.Where("skill_id = ?", id).Preload("Careers").Preload("SkillsLevels").First(&skill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete associated records in careers_skills table
	err = repository.Db.Exec("DELETE FROM careers_skills WHERE skill_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete associated records in skills_levels table
	err = repository.Db.Exec("DELETE FROM skills_levels WHERE skill_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete the skill record
	err = repository.Db.Delete(&skill).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill and associated records deleted successfully"})
}
