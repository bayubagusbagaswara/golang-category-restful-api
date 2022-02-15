package repository

import (
	"context"
	"database/sql"
	"errors"
	"golangrestfulapi/helper"
	"golangrestfulapi/model/domain"
)

type CategoryRepositoryImpl struct {
}

func (repository *CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	// buat dulu SQL nya
	SQL := "insert into category(name) values(?)"

	// masukan ctx, SQLnya, dan argument yang dimasukkan, hanya name
	r, err := tx.ExecContext(ctx, SQL, category.Name)
	helper.PanicIfError(err)

	// dapatkan id terakhirnya
	id, err2 := r.LastInsertId()
	helper.PanicIfError(err2)

	// id nya kita ambil dari auto generate, sehingga akan ditambahkan dari id yang terakhir
	category.Id = int(id)
	return category
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	SQL := "update category set name = ? where id = ?"
	_, err := tx.ExecContext(ctx, SQL, category.Name, category.Id)
	helper.PanicIfError(err)

	// disini kita tidak perlu ubah lagi id nya
	return category
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, category domain.Category) {
	SQL := "delete friom category where id = ?"
	_, err := tx.ExecContext(ctx, SQL, category.Id)
	helper.PanicIfError(err)

	// tidak perlu return
}

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) (domain.Category, error) {
	SQL := "select id, name from category where id = ?"
	r, err := tx.QueryContext(ctx, SQL, categoryId)
	helper.PanicIfError(err)

	category := domain.Category{}

	if r.Next() {
		// jika datanya ada, maka kita ambil
		err2 := r.Scan(&categoryId, &category.Name)
		helper.PanicIfError(err2)
		return category, nil

	} else {
		// jika tidak ada datanya
		return category, errors.New("category is not found")
	}
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Category {
	SQL := "select id, name form category"
	r, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)

	var categories []domain.Category
	// selama datanya ada, maka kita balikkan datanya
	for r.Next() {
		category := domain.Category{}
		err2 := r.Scan(&category.Id, &category.Name)
		helper.PanicIfError(err2)
		categories = append(categories, category)
	}
	return categories
}
