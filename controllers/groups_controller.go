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

type GroupRepo struct {
	Db *gorm.DB
}

func NewGroupRepo() *GroupRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Group{}, &models.CareersInGroup{})
	return &GroupRepo{Db: db}
}

// GetGroups retrieves all Group records from the database.
func (repository *GroupRepo) GetGroups(c *gin.Context) {
	var groups []models.Group

	err := models.Getgroups(repository.Db, &groups)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetAllGroupsHaveSection retrieves all groups from the database where each group must have at least one section.
func (repository *GroupRepo) GetAllGroupsHaveSection(c *gin.Context) {
	var groups []models.Group

	err := models.Getgroups(repository.Db, &groups) // Set to true to filter groups with sections
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetGroupById retrieves a Group by its ID from the database.
func (repository *GroupRepo) GetGroupById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var group models.Group

	err := models.GetgroupById(repository.Db, &group, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, group)
}

// CreateGroup creates a new Group record.
func (repository *GroupRepo) CreateGroup(c *gin.Context) {
	var group models.Group

	if err := c.ShouldBindJSON(&group); err != nil {
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
	var existingGroup models.Group
	if err := repository.Db.Where("name = ?", group.Name).First(&existingGroup).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}

	err := models.Creategroup(repository.Db, &group)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// UpdateGroup updates a Group record by ID.
func (repository *GroupRepo) UpdateGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedGroup models.Group

	err := repository.Db.First(&updatedGroup, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&updatedGroup); err != nil {
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
	var existingGroup models.Group
	if err := repository.Db.Where("name = ? AND group_id != ?", updatedGroup.Name, updatedGroup.GroupID).First(&existingGroup).Error; err == nil {

		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}

	var currentUpdatedGroup models.UpdateGroupModels

	currentUpdatedGroup.GroupID = updatedGroup.GroupID
	currentUpdatedGroup.Name = updatedGroup.Name
	currentUpdatedGroup.Sections = updatedGroup.Sections
	currentUpdatedGroup.Careers = updatedGroup.Careers

	err = models.Updategroup(repository.Db, &currentUpdatedGroup)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, updatedGroup)
}

// DeleteGroupById deletes a Group record by ID.
func (repository *GroupRepo) DeleteGroupById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var group models.Group

	err := repository.Db.Where("group_id = ?", id).First(&group).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Delete associated records in sections_groups table
	err = repository.Db.Exec("DELETE FROM sections_groups WHERE group_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete associated records in groups_careers table
	err = repository.Db.Exec("DELETE FROM groups_careers WHERE group_id = ?", id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Delete the group record
	err = repository.Db.Delete(&group).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group and associated records deleted successfully"})
}
