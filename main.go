package main

import (
	"knowledgeplus/go-api/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// func init() {
// 	initializers.LoadEnvVariables()
// }

func main() {
	r := setupRouter()
	_ = r.Run(":8081")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	defaultPath := r.Group("/api")

	// CareerRepo := controllers.NewCareerRepo()
	// CategoriesRepo := controllers.NewCategoriesRepo()
	// SkillRepo := controllers.NewSkillRepo()
	// LevelsRepo := controllers.NewLevelsRepo()
	// CourseRepo := controllers.NewCourseRepo()
	OrganizationRepo := controllers.NewOrganizationsRepo()

	// defaultPath.POST("/careers", CareerRepo.CreateCareer)
	// defaultPath.GET("/careers", CareerRepo.GetCareers)
	// defaultPath.GET("/careers_have_categories", CareerRepo.GetCareersWithHaveCategories)
	// defaultPath.GET("/careers/:id", CareerRepo.GetCareer)
	// defaultPath.PUT("/careers/:id", CareerRepo.UpdateCareer)
	// defaultPath.DELETE("/careers/:id", CareerRepo.DeleteCareer)

	// defaultPath.GET("/categories", CategoriesRepo.GetCategories)
	// defaultPath.GET("/categories/:id", CategoriesRepo.GetCategoryById)
	// defaultPath.POST("/categories", CategoriesRepo.CreateCategory)
	// defaultPath.PUT("/categories/:id", CategoriesRepo.UpdateCategory)
	// defaultPath.DELETE("/categories/:id", CategoriesRepo.DeleteCategoryById)

	// defaultPath.GET("/skills", SkillRepo.GetSkills)
	// defaultPath.GET("/skills/:id", SkillRepo.GetSkillById)
	// defaultPath.POST("/skills", SkillRepo.CreateSkill)
	// defaultPath.PUT("/skills/:id", SkillRepo.UpdateSkill)
	// defaultPath.DELETE("/skills/:id", SkillRepo.DeleteSkillById)

	// defaultPath.GET("/levels", LevelsRepo.GetLevels)
	// defaultPath.GET("/levels/:id", LevelsRepo.GetLevelById)

	// defaultPath.GET("/courses", CourseRepo.GetCourses)
	// defaultPath.GET("/courses/:id", CourseRepo.GetCourse)

	defaultPath.POST("/organizations", OrganizationRepo.CreateOrganization)
	defaultPath.GET("/organizations", OrganizationRepo.GetOrganizations)
	defaultPath.GET("/organizations/:id", OrganizationRepo.GetOrganizationById)
	defaultPath.PUT("/organizations/:id", OrganizationRepo.UpdateOrganization)
	defaultPath.DELETE("/organizations/:id", OrganizationRepo.DeleteOrganization)

	return r
}
