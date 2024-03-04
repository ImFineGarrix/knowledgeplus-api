package controllers

import (
	"errors"
	"knowledgeplus/go-api/database"
	"knowledgeplus/go-api/models"
	"knowledgeplus/go-api/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type SectionRepo struct {
	Db *gorm.DB
}

func NewSectionRepo() *SectionRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Section{}, &models.Groups{})
	return &SectionRepo{Db: db}
}

// GetSections retrieves all Section records from the database.
func (repository *SectionRepo) GetSections(c *gin.Context) {
	var sections []models.Section

	err := models.GetSections(repository.Db, &sections)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, sections)

}

// GetSectionById retrieves a Section by its ID from the database.
func (repository *SectionRepo) GetSectionById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var section models.Section

	err := models.GetsectionById(repository.Db, &section, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, section)
}

// CreateSection creates a new Section record.
func (repository *SectionRepo) CreateSection(c *gin.Context) {
	var section models.Section

	if err := c.ShouldBindJSON(&section); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]response.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = response.ErrorMsg{
					Code:    http.StatusBadRequest,
					Field:   fe.Field(),
					Message: response.GetErrorMsg(fe),
				}
			}
			c.JSON(http.StatusCreated, out)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}

	// Check if the name already exists in the database
	var existingSection models.Section
	if err := repository.Db.Where("name = ?", section.Name).First(&existingSection).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}
	err := models.Createsection(repository.Db, &section)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, section)
}

// UpdateSection updates a Section record by ID.
func (repository *SectionRepo) UpdateSection(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var existingSection models.Section

	// Fetch the existing section by ID
	err := repository.Db.First(&existingSection, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var updatedSection models.UpdateSectionModels
	if err := c.ShouldBindJSON(&updatedSection); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]response.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = response.ErrorMsg{
					Code:    http.StatusBadRequest,
					Field:   fe.Field(),
					Message: response.GetErrorMsg(fe),
				}
			}
			c.JSON(http.StatusCreated, out)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}

	// Check if the name already exists in the database
	var existingSections models.Section
	if err := repository.Db.Where("name = ? AND section_id != ?", updatedSection.Name, updatedSection.SectionID).First(&existingSections).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// Update only the fields you want to change
	existingSection.Name = updatedSection.Name
	existingSection.ImageUrl = updatedSection.ImageUrl

	// Save the changes
	err = repository.Db.Save(&existingSection).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, existingSection)
}

// DeleteSectionById deletes a Section record by ID.
func (repository *SectionRepo) DeleteSectionById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var section models.Section

	err := repository.Db.Where("section_id = ?", id).First(&section).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	// Delete associated records in sections_groups table
	err = repository.Db.Exec("DELETE FROM sections_groups WHERE section_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete the section record
	err = repository.Db.Delete(&section).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Section and associated records deleted successfully"})
}
