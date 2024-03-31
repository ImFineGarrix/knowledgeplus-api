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

type OrganizationsRepo struct {
	Db *gorm.DB
}

func NewOrganizationsRepo() *OrganizationsRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Organizations{}, &models.Organizations{})
	return &OrganizationsRepo{Db: db}
}

// get Organizationss
func (repository *OrganizationsRepo) GetOrganizations(c *gin.Context) {
	var Organizations []models.Organizations
	print(Organizations)
	err := models.GetOrganizations(repository.Db, &Organizations)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Organizations)
}

// get Organizations by id
func (repository *OrganizationsRepo) GetOrganizationById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var Organizations models.Organizations
	err := models.GetOrganizationById(repository.Db, &Organizations, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Organizations)
}

func (repository *OrganizationsRepo) CreateOrganization(c *gin.Context) {
	var input models.Organizations
	if err := c.ShouldBindJSON(&input); err != nil {
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
	var existingOrganizations models.Organizations
	if err := repository.Db.Where("name = ?", input.Name).First(&existingOrganizations).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}

	err := models.CreateOrganization(repository.Db, &input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (repository *OrganizationsRepo) UpdateOrganization(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var existingOrganization models.Organizations
	var updatedOrganization models.UpdateOrganizationModels
	var currentOrganization models.UpdateOrganizationModels
	err := repository.Db.First(&existingOrganization, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&updatedOrganization); err != nil {
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
	var existingOrganizations models.Organizations
	if err := repository.Db.Where("name = ? AND organization_id != ?", updatedOrganization.Name, id).First(&existingOrganizations).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}

	currentOrganization.Name = updatedOrganization.Name
	currentOrganization.Description = updatedOrganization.Description
	currentOrganization.ImageUrl = updatedOrganization.ImageUrl

	// existingOrganization.Name = updatedOrganization.Name
	// Update other fields as needed

	err = models.UpdateOrganization(repository.Db, id, &currentOrganization)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, currentOrganization)
}

func (repository *OrganizationsRepo) DeleteOrganization(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var organization models.Organizations

	err := models.GetOrganizationById(repository.Db, &organization, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = models.DeleteOrganization(repository.Db, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}
