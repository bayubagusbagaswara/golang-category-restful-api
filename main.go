package main

import (
	"golangrestfulapi/app"
	"golangrestfulapi/controller"
	"golangrestfulapi/repository"
	"golangrestfulapi/service"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

func main() {

	// bikin validator
	validator := validator.New()

	// bikin db nya
	db := app.NewDB()

	// bikin categoryRepository
	categoryRepository := repository.NewCategoryRepository()

	// bikin service
	categoryService := service.NewCategoryService(categoryRepository, db, validator)

	// bikin controller
	categoryController := controller.NewCategoryController(categoryService)

	// bikin router nya
	router := httprouter.New()

	router.GET("/api/categories", categoryController.FindAll)
	router.GET("/api/categories/:categoryId", categoryController.FindById)
	router.POST("/api/categories", categoryController.Create)
	router.PUT("/api/categories/:categorId", categoryController.Update)
	router.DELETE("/api/categories/:categoryId", categoryController.Delete)

}
