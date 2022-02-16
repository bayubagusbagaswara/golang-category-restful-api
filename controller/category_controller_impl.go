package controller

import (
	"golangrestfulapi/helper"
	"golangrestfulapi/model/web"
	"golangrestfulapi/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// bikin struct yang isinya adalah data/field yang dibutuhkan function implementasi dari CategoryController
// kita butuh CategoryService, karena CategoryService adalah interface, maka tidak perlu diset sebagai pointer
type CategoryControllerImpl struct {
	CategoryService service.CategoryService
}

func (controller *CategoryControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	// buat categoryCreateRequest, dimana hanya data name
	categoryCreateRequest := web.CategoryCreateRequest{}
	helper.ReadFromRequestBody(request, &categoryCreateRequest)

	// lalu panggil service nya, masukkan context dan datanya, yakni categoryCreateRequest
	categoryResponse := controller.CategoryService.Create(request.Context(), categoryCreateRequest)

	// buat response, sesuai dengan response yang baku web response
	webResponse := web.WebResponse{
		Code:   201,
		Status: "CREATED",
		Data:   categoryResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CategoryControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	categoryUpdateRequest := web.CategoryUpdateRequest{}
	helper.ReadFromRequestBody(request, &categoryUpdateRequest)

	// karena kita mengirim parameter berupa id, maka kita harus tangkap id nya
	categoryId := params.ByName("categoryId")
	// konversi dulu menjadi integer
	id, err2 := strconv.Atoi(categoryId)
	helper.PanicIfError(err2)

	// masukkan id kedalam categoryUpdateRequest
	categoryUpdateRequest.Id = id

	categoryResponse := controller.CategoryService.Update(request.Context(), categoryUpdateRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CategoryControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	categoryId := params.ByName("categoryId")
	id, err2 := strconv.Atoi(categoryId)
	helper.PanicIfError(err2)

	controller.CategoryService.Delete(request.Context(), id)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CategoryControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryId := params.ByName("categoryId")
	id, err2 := strconv.Atoi(categoryId)
	helper.PanicIfError(err2)

	categoryResponse := controller.CategoryService.FindById(request.Context(), id)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CategoryControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	categoryResponses := controller.CategoryService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
