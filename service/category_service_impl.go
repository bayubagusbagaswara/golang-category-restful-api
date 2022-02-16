package service

import (
	"context"
	"database/sql"
	"golangrestfulapi/helper"
	"golangrestfulapi/model/domain"
	"golangrestfulapi/model/web"
	"golangrestfulapi/repository"

	"github.com/go-playground/validator/v10"
)

// struct disini berguna untuk menyimpan semua data/field/property yang dibutuhkan oleh CategoryServiceImpl
type CategoryServiceImpl struct {

	// disini kita butuh repository, karena untuk manipulasi datanya melalui repository
	CategoryRepository repository.CategoryRepository
	// kita buth koneksi ke databasenya
	DB *sql.DB
	// kita butuh vaidate juga
	Validate *validator.Validate
}

// bikin function untuk membuat CategoryService, ini seperti Constructor
func NewCategoryService(categoryRepository repository.CategoryRepository, DB *sql.DB, validate *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
		DB:                 DB,
		Validate:           validate,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse {
	// sebelum mengirimkan data name, maka kita lakukan validasi
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	// karena kita menggunakan database transactional, maka untuk request kita butuh transactional
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)

	// pertama kita method diatas diwrap trancational, tapi juga saat terjadi error maka kita rollback
	defer helper.CommitOrRollback(tx)

	// create category
	category := domain.Category{
		Name: request.Name,
	}

	c := service.CategoryRepository.Save(ctx, tx, category)
	return helper.ToCategoryResponse(c)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {
	// lalukan validasi dulu datanya
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// cek apakah ada category di database
	category, err := service.CategoryRepository.FindById(ctx, tx, request.Id)
	helper.PanicIfError(err)

	// jika ada data category, maka kita ubah data name category nya dengan data baru dari request
	category.Name = request.Name

	c := service.CategoryRepository.Update(ctx, tx, category)
	return helper.ToCategoryResponse(c)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)
	// cek apakah ada category di database
	category, err := service.CategoryRepository.FindById(ctx, tx, categoryId)
	helper.PanicIfError(err)

	// jika category nya ada di database, maka bisa kita hapus
	service.CategoryRepository.Delete(ctx, tx, category)
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	c, err2 := service.CategoryRepository.FindById(ctx, tx, categoryId)
	helper.PanicIfError(err2)

	return helper.ToCategoryResponse(c)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	categoryList := service.CategoryRepository.FindAll(ctx, tx)

	return helper.ToCategoryResponses(categoryList)
}
