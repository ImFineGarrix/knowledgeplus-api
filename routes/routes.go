package routes

import (
	"knowledgeplus/go-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(defaultPath *gin.RouterGroup) {
	// Initialize middleware
	// authMiddleware := middleware.AuthMiddleware()

	defaultPath.GET("/auth/login", controllers.NewAuthRepo().LoginHandler)
	// defaultPath.Use(authMiddleware).POST("/auth/register", controllers.NewAuthRepo().CreateUserHandler)
	defaultPath.POST("/auth/register", controllers.NewAuthRepo().CreateUserHandler)

	SectionRepo := controllers.NewSectionRepo()
	GroupRepo := controllers.NewGroupRepo()
	CareerRepo := controllers.NewCareerRepo()
	// CategoriesRepo := controllers.NewCategoriesRepo()
	SkillRepo := controllers.NewSkillRepo()
	LevelsRepo := controllers.NewLevelsRepo()
	CourseRepo := controllers.NewCourseRepo()
	OrganizationRepo := controllers.NewOrganizationsRepo()

	defaultPath.POST("/careers", CareerRepo.CreateCareer)
	defaultPath.GET("/careers", CareerRepo.GetCareers)
	// defaultPath.GET("/careers_have_categories", CareerRepo.GetCareersWithHaveCategories)
	defaultPath.GET("/careers/:id", CareerRepo.GetCareer)
	defaultPath.PUT("/careers/:id", CareerRepo.UpdateCareer)
	defaultPath.DELETE("/careers/:id", CareerRepo.DeleteCareer)

	defaultPath.GET("/skills", SkillRepo.GetSkills)
	defaultPath.GET("/skills/:id", SkillRepo.GetSkillById)
	defaultPath.POST("/skills", SkillRepo.CreateSkill)
	defaultPath.PUT("/skills/:id", SkillRepo.UpdateSkill)
	defaultPath.DELETE("/skills/:id", SkillRepo.DeleteSkillById)

	defaultPath.GET("/levels", LevelsRepo.GetLevels)
	defaultPath.GET("/levels/:id", LevelsRepo.GetLevelById)

	// defaultPath.GET("/courses", CourseRepo.GetCourses)
	// defaultPath.GET("/courses/:id", CourseRepo.GetCourse)

	defaultPath.POST("/organizations", OrganizationRepo.CreateOrganization)
	defaultPath.GET("/organizations", OrganizationRepo.GetOrganizations)
	defaultPath.GET("/organizations/:id", OrganizationRepo.GetOrganizationById)
	defaultPath.PUT("/organizations/:id", OrganizationRepo.UpdateOrganization)
	defaultPath.DELETE("/organizations/:id", OrganizationRepo.DeleteOrganization)

	defaultPath.GET("/sections", SectionRepo.GetSections)
	defaultPath.GET("/sections/:id", SectionRepo.GetSectionById)
	defaultPath.POST("/sections", SectionRepo.CreateSection)
	defaultPath.PUT("/sections/:id", SectionRepo.UpdateSection)
	defaultPath.DELETE("/sections/:id", SectionRepo.DeleteSectionById)

	defaultPath.GET("/groups", GroupRepo.GetGroups)
	defaultPath.GET("/groups/:id", GroupRepo.GetGroupById)
	defaultPath.POST("/groups", GroupRepo.CreateGroup)
	defaultPath.PUT("/groups/:id", GroupRepo.UpdateGroup)
	defaultPath.DELETE("/groups/:id", GroupRepo.DeleteGroupById)

	defaultPath.GET("/courses", CourseRepo.GetCourses)
	defaultPath.GET("/courses/:id", CourseRepo.GetCourseById)
	defaultPath.POST("/courses", CourseRepo.CreateCourse)
	defaultPath.PUT("/courses/:id", CourseRepo.UpdateCourse)
	defaultPath.DELETE("/courses/:id", CourseRepo.DeleteCourseById)

}
