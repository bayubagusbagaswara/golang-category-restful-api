package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// bkin interface dulu untuk kontrak function yang akan digunakan semua API
// functionnya harus standar, karena harus mengikuti handler dari si HTTP
// karena kita menggunakan httprouter, maka parameter untuk functionnya, selain ada writer(response) dan request
type CategoryController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
