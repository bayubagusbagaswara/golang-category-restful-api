package app

import (
	"golangrestfulapi/controller"
	"golangrestfulapi/exception"

	"github.com/julienschmidt/httprouter"
)

// bikin function untuk create router
func NewRouter(categoryController controller.CategoryController) *httprouter.Router {

	router := httprouter.New()

	router.GET("/api/categories", categoryController.FindAll)
	router.GET("/api/categories/:categoryId", categoryController.FindById)
	router.POST("/api/categories", categoryController.Create)
	router.PUT("/api/categories/:categorId", categoryController.Update)
	router.DELETE("/api/categories/:categoryId", categoryController.Delete)

	router.PanicHandler = exception.ErrorHandler
	return router
}
