package main

import (
	"golangrestfulapi/app"
	"golangrestfulapi/controller"
	"golangrestfulapi/exception"
	"golangrestfulapi/helper"
	"golangrestfulapi/repository"
	"golangrestfulapi/service"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

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

	// bikin untuk router jika terjadi error
	router.PanicHandler = exception.ErrorHandler

	// bikin server
	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}
	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
