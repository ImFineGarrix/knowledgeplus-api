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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	err := models.GetOrganizationById(repository.Db, &existingOrganization, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&updatedOrganization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentOrganization.Name = updatedOrganization.Name
	currentOrganization.Description = updatedOrganization.Description
	currentOrganization.Description = updatedOrganization.Description

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
