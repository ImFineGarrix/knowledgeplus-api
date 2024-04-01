package controllers

import (
	"errors"
	"knowledgeplus/go-api/models"
	"knowledgeplus/go-api/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type SectionRepo struct {
	Db      *gorm.DB
	UserDb  *gorm.DB
	AdminDb *gorm.DB
}

// ทำถึงแค่นี้
func NewSectionRepo(db *gorm.DB, dbUser *gorm.DB, adminDb *gorm.DB) *SectionRepo {
	db.AutoMigrate(&models.Section{}, &models.Groups{})
	dbUser.AutoMigrate(&models.Section{}, &models.Groups{})
	adminDb.AutoMigrate(&models.Section{}, &models.Groups{})
	return &SectionRepo{Db: db, UserDb: dbUser, AdminDb: adminDb}
}

// GetSections retrieves all Section records from the database.
func (repository *SectionRepo) GetSections(c *gin.Context) {
	var sections []models.Section
	// Check if the key "userRole" exists in the context
	var repoByRole = repository.Db
	if userRole, ok := c.Get("userRole"); ok {
		switch userRole {
		case "user":
			log.Default().Println("It's user")
			repoByRole = repository.UserDb
		case "admin":
			log.Default().Println("It's admin")
			repoByRole = repository.AdminDb
		case "owner":
			log.Default().Println("It's owner")
			repoByRole = repository.Db
		default:
			log.Default().Println("It's default")
		}
	} else {
		log.Default().Println("userRole key does not exist in context")
		repoByRole = repository.UserDb
	}

	err := models.GetSections(repoByRole, &sections)
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
	// Check if the key "userRole" exists in the context
	var repoByRole = repository.Db
	if userRole, ok := c.Get("userRole"); ok {
		switch userRole {
		case "user":
			repoByRole = repository.UserDb
		case "admin":
			repoByRole = repository.AdminDb
		case "owner":
			repoByRole = repository.Db
		default:
		}
	} else {
		repoByRole = repository.UserDb
	}

	err := models.GetsectionById(repoByRole, &section, id)
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

	// Check if the key "userRole" exists in the context
	var repoByRole = repository.Db
	if userRole, ok := c.Get("userRole"); ok {
		switch userRole {
		case "user":
			repoByRole = repository.UserDb
		case "admin":
			repoByRole = repository.AdminDb
		case "owner":
			repoByRole = repository.Db
		default:
		}
	} else {
		repoByRole = repository.UserDb
	}

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
			c.JSON(http.StatusBadRequest, out)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}

	// Check if the name already exists in the database
	var existingSection models.Section
	if err := repoByRole.Where("name = ?", section.Name).First(&existingSection).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}
	err := models.Createsection(repoByRole, &section)
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

	// Check if the key "userRole" exists in the context
	var repoByRole = repository.Db
	if userRole, ok := c.Get("userRole"); ok {
		switch userRole {
		case "user":
			repoByRole = repository.UserDb
		case "admin":
			repoByRole = repository.AdminDb
		case "owner":
			repoByRole = repository.Db
		default:
		}
	} else {
		repoByRole = repository.UserDb
	}

	// Fetch the existing section by ID
	err := repoByRole.First(&existingSection, id).Error
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
			c.JSON(http.StatusBadRequest, out)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}

	// Check if the name already exists in the database
	var existingSections models.Section
	if err := repoByRole.Where("name = ? AND section_id != ?", updatedSection.Name, id).First(&existingSections).Error; err == nil {
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
	err = repoByRole.Save(&existingSection).Error
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

	// Check if the key "userRole" exists in the context
	var repoByRole = repository.Db
	if userRole, ok := c.Get("userRole"); ok {
		switch userRole {
		case "user":
			repoByRole = repository.UserDb
		case "admin":
			repoByRole = repository.AdminDb
		case "owner":
			repoByRole = repository.Db
		default:
		}
	} else {
		repoByRole = repository.UserDb
	}

	err := repoByRole.Where("section_id = ?", id).First(&section).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	// Delete associated records in sections_groups table
	err = repoByRole.Exec("DELETE FROM sections_groups WHERE section_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete the section record
	err = repoByRole.Delete(&section).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Section and associated records deleted successfully"})
}
