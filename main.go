package main

import (
	"golangrestfulapi/app"
	"golangrestfulapi/controller"
	"golangrestfulapi/helper"
	"golangrestfulapi/middleware"
	"golangrestfulapi/repository"
	"golangrestfulapi/service"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-playground/validator/v10"
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
	router := app.NewRouter(categoryController)

	// bikin server
	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}
	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
