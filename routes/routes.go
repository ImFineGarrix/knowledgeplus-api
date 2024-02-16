package routes

import (
	"knowledgeplus/go-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(defaultPath *gin.RouterGroup) {
	// Initialize middleware
	// authMiddleware := middleware.AuthMiddleware()

	defaultPath.POST("/auth/login", controllers.NewAuthRepo().LoginHandler)
	// defaultPath.Use(authMiddleware).POST("/auth/register", controllers.NewAuthRepo().CreateUserHandler)
	// defaultPath.POST("/admins", controllers.NewAuthRepo().CreateUserHandler)

	SectionRepo := controllers.NewSectionRepo()
	GroupRepo := controllers.NewGroupRepo()
	CareerRepo := controllers.NewCareerRepo()
	// CategoriesRepo := controllers.NewCategoriesRepo()
	SkillRepo := controllers.NewSkillRepo()
	LevelsRepo := controllers.NewLevelsRepo()
	CourseRepo := controllers.NewCourseRepo()
	OrganizationRepo := controllers.NewOrganizationsRepo()
	UserRepo := controllers.NewUserRepo()

	/** careers models */

	// backoffice
	defaultPath.GET("/backoffice/careers", CareerRepo.GetAllCareersWithFilters)
	defaultPath.GET("/backoffice/careers/:id", CareerRepo.GetCareer)
	defaultPath.POST("/backoffice/careers", CareerRepo.CreateCareer)
	defaultPath.PUT("/backoffice/careers/:id", CareerRepo.UpdateCareer)
	defaultPath.DELETE("/backoffice/careers/:id", CareerRepo.DeleteCareer)

	//frontend
	defaultPath.GET("/careers", CareerRepo.GetCareers)
	defaultPath.GET("/careers/:id", CareerRepo.GetCareer)
	defaultPath.GET("/careers-by-course/:course_id", CareerRepo.GetCareersByCourseId)
	defaultPath.GET("/careers-by-skill/:skill_id", CareerRepo.GetCareersBySkillId)

	/** skills **/
	//backoffice

	defaultPath.GET("/backoffice/skills/:id", SkillRepo.GetSkillById)
	defaultPath.POST("/backoffice/skills", SkillRepo.CreateSkill)
	defaultPath.PUT("/backoffice/skills/:id", SkillRepo.UpdateSkill)
	defaultPath.DELETE("/backoffice/skills/:id", SkillRepo.DeleteSkillById)
	//frontend
	defaultPath.GET("/skills", SkillRepo.GetSkills)
	defaultPath.GET("/skills-by-course/:course_id", SkillRepo.GetSkillsByCourseId)
	defaultPath.GET("/skills-by-career/:career_id", SkillRepo.GetSkillsByCareerId)
	defaultPath.GET("/skills/:id", SkillRepo.GetSkillById)

	/** levels **/
	//backoffice
	defaultPath.GET("/backoffice/levels", LevelsRepo.GetLevels)
	defaultPath.GET("/backoffice/levels/:id", LevelsRepo.GetLevelById)
	//frontend
	defaultPath.GET("/levels", LevelsRepo.GetLevels)
	defaultPath.GET("/levels/:id", LevelsRepo.GetLevelById)

	/** organizations **/
	//backoffice
	defaultPath.GET("/backoffice/organizations", OrganizationRepo.GetOrganizations)
	defaultPath.GET("/backoffice/organizations/:id", OrganizationRepo.GetOrganizationById)
	defaultPath.POST("/backoffice/organizations", OrganizationRepo.CreateOrganization)
	defaultPath.PUT("/backoffice/organizations/:id", OrganizationRepo.UpdateOrganization)
	defaultPath.DELETE("/backoffice/organizations/:id", OrganizationRepo.DeleteOrganization)
	//frontend
	defaultPath.GET("/organizations", OrganizationRepo.GetOrganizations)
	defaultPath.GET("/organizations/:id", OrganizationRepo.GetOrganizationById)

	/** sections **/
	//backoffice
	defaultPath.GET("/backoffice/sections", SectionRepo.GetSections)
	defaultPath.GET("/backoffice/sections/:id", SectionRepo.GetSectionById)
	defaultPath.POST("/backoffice/sections", SectionRepo.CreateSection)
	defaultPath.PUT("/backoffice/sections/:id", SectionRepo.UpdateSection)
	defaultPath.DELETE("/backoffice/sections/:id", SectionRepo.DeleteSectionById)
	//frontend
	defaultPath.GET("/sections", SectionRepo.GetSections)
	defaultPath.GET("/sections/:id", SectionRepo.GetSectionById)

	/** groups **/
	//backoffice
	defaultPath.GET("/backoffice/groups", GroupRepo.GetGroups)
	defaultPath.GET("/backoffice/groups/:id", GroupRepo.GetGroupById)
	defaultPath.POST("/backoffice/groups", GroupRepo.CreateGroup)
	defaultPath.PUT("/backoffice/groups/:id", GroupRepo.UpdateGroup)
	defaultPath.DELETE("/backoffice/groups/:id", GroupRepo.DeleteGroupById)
	//frontend
	defaultPath.GET("/groups", GroupRepo.GetGroups)
	defaultPath.GET("/groups-with-section", GroupRepo.GetAllGroupsHaveSection)
	defaultPath.GET("/groups/:id", GroupRepo.GetGroupById)

	/** courses **/
	//backoffice
	defaultPath.GET("/backoffice/courses", CourseRepo.GetCourses)
	defaultPath.GET("/backoffice/courses/:id", CourseRepo.GetCourseById)
	defaultPath.POST("/backoffice/courses", CourseRepo.CreateCourse)
	defaultPath.PUT("/backoffice/courses/:id", CourseRepo.UpdateCourse)
	defaultPath.DELETE("/backoffice/courses/:id", CourseRepo.DeleteCourseById)
	//frontend
	defaultPath.GET("/courses", CourseRepo.GetCourses)
	defaultPath.GET("/courses/:id", CourseRepo.GetCourseById)
	defaultPath.GET("/courses-by-skill/:skill_id", CourseRepo.GetCoursesBySkillId)
	defaultPath.GET("/courses-by-career/:career_id", CourseRepo.GetCoursesByCareerId)

	/** admins **/
	//backoffice
	defaultPath.POST("/backoffice/admins", UserRepo.CreateUser)
	defaultPath.GET("/backoffice/admins", UserRepo.GetUsers)
	defaultPath.GET("/backoffice/admins/:id", UserRepo.GetUserById)
	defaultPath.PUT("/backoffice/admins/:id", UserRepo.UpdateUser)
	defaultPath.DELETE("/backoffice/admins/:id", UserRepo.DeleteUserById)

}
