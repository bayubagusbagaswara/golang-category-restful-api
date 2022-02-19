package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"golangrestfulapi/app"
	"golangrestfulapi/controller"
	"golangrestfulapi/helper"
	"golangrestfulapi/middleware"
	"golangrestfulapi/model/domain"
	"golangrestfulapi/repository"
	"golangrestfulapi/service"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

// kita ngehit langusng endpoint nya
// atau langsung integration test

func setupTestDB() *sql.DB {
	// db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/golang_restful_api_test")
	db, err := sql.Open("mysql", "root:root@/golang_restful_api_test")
	helper.PanicIfError(err)

	// set connection pooling
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {

	validator := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validator)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	// masukkan middlewarenya sekalian
	return middleware.NewAuthMiddleware(router)
}

// create function
func truncateCategory(db *sql.DB) {
	// setiap menjalankan test, kita akan hapus dulu data di database
	db.Exec("TRUNCATE category")
}

// skenario pertama unutk Test Create Category
func TestCreateCategorySuccess(t *testing.T) {

	db := setupTestDB()
	// hapus atau truncate dulu data category di database
	truncateCategory(db)

	// bikin routernya
	router := setupRouter(db)

	// bikin request body
	requestBody := strings.NewReader(`{"name" : "Gadget"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)

	// buat header di requestnya
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	// panggil routernya
	router.ServeHTTP(recorder, request)

	// dapetin responsenya
	response := recorder.Result()

	// lalu bisa kita baca body responsenya, atau kita cek apakah sesuai dengan request kita
	assert.Equal(t, 201, response.StatusCode)

	// kita lihat isi body nya, baca semua isi body nya
	body, _ := io.ReadAll(response.Body)

	// bikin variable map dulu, untuk menyimpan data dari body
	var responseBody map[string]interface{}

	// decode data body
	json.Unmarshal(body, &responseBody)

	// web response
	assert.Equal(t, 200, responseBody["code"].(float64))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
}
func TestCreateCategoryFailed(t *testing.T) {
	// test untuk create failed, misalnya data yang kita kirim tidak lolos validasi
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name" : "G"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 201, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, responseBody["code"].(float64))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
	// dan response nya tidak ada datanya, karena gagal
}
func TestUpdateCategorySuccess(t *testing.T) {

	db := setupTestDB()
	truncateCategory(db)

	// karena data di database sudah terhapus, maka kita coba insert data baru dulu
	// dengan menggunakan repository, kita coba membuat 1 category gadget
	// kita bungkus didalam transaction
	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	// disini kita akan mendapatkan data balikan berupa id
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	// bikin routernya
	router := setupRouter(db)

	// bikin request
	requestBody := strings.NewReader(`{"name" : "Gadget"}`)

	// saat kirim request, kita mengirim data id (dari data yang sudah kita create sebelumnya) data body yang berisi data baru
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	// bikin response dari recorder
	recorder := httptest.NewRecorder()

	// jalankan routernya
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	// bikin variable untuk menyimpan data response
	var responseBody map[string]interface{}

	// encode
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	// id dari body nya harus sama dengan id yang kita kirim
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
}

func TestUpdateCategoryFailed(t *testing.T) {
	// update category failed, misal data baru yang dikirim tidak lolos validasi, jadi BAD REQUEST
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setupRouter(db)

	// misal data yang kirim ini gagal validasi, alias data name kita kosongkan
	requestBody := strings.NewReader(`{"name" : ""}`)

	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
	// jadi tidak ada data balikannya
}
func TestGetCategorySuccess(t *testing.T) {

	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setupRouter(db)
	// tidak membutuhkan body, karena kita hanya mengambil data nya, jadi kita isi nil
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, category.Name, responseBody["data"].(map[string]interface{})["name"])
}
func TestGetCategoryFailed(t *testing.T) {

	// kita coba yang category nya not found
	// artinya id yang kita kirim, tidak ditemukan datanya
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/404", nil)
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}
func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	// kita create category, dan kita mendapatkan id category
	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	// bikin router
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	// jalankan router
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}
func TestDeleteCategoryFailed(t *testing.T) {
	// delete failed, artinya tidak ada data yang dimaksud di database
	// jadi kita tidak perlu create data terlebih dahulu
	// dan kita masukkan id sembarang (id tidak ada di database)
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/404", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}
func TestListCategoryiesSuccess(t *testing.T) {

	// kita cukup insert beberapa category
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category1 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	category2 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Computer",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	// balikan datanya berupa slice
	var categories = responseBody["data"].([]interface{})

	categoryResponse1 := categories[0].(map[string]interface{})
	categoryResponse2 := categories[1].(map[string]interface{})

	assert.Equal(t, category1.Id, int(categoryResponse1["id"].(float64)))
	assert.Equal(t, category1.Name, categoryResponse1["name"])
	assert.Equal(t, category2.Id, int(categoryResponse2["id"].(float64)))
	assert.Equal(t, category2.Name, categoryResponse2["name"])
}

// test Unauthorized, misalnya API Key nya disalahkan
func TestUnauthorized(t *testing.T) {

	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)

	request.Header.Add("Content-Type", "application/json")

	// tinggal gunakan API KEY yang salah, agar terkena Authorized
	request.Header.Add("X-API-Key", "SALAH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "UNAUTHORIZED", responseBody["status"])
}
