package repository

import (
	"context"
	"database/sql"
	"golangrestfulapi/model/domain"
)

// sebelumnya kita harus membuat kontraknya dulu
// kontraknya sendiri menggunakan interface
type CategoryRepository interface {

	// kita buat operasinya, yakni function untuk CRUD
	// untuk parameternya harus diawali dengan Context
	// jadi parameter pertama adalah context

	// saat menggunakan database relasional, ada baiknya pada function yang dibuat di repository itu harus mendukung transactional
	// di golang untuk melakukan eksekusi menggunakan database transactional itu harus menggunakan method berbeda dengan biasanya
	// biasaya mungkin langsung menggunakan dari DB, seperti db.Execute
	// kalau transactional harus menggunakan db.Begin, baru nanti jadi object transactional, baru gunakan object trasactional itu
	// akhirnya nanti kita bikin 2 function yang berbeda, yakni yang transactional dan tidak transactional
	// tapi ini kita menggunakan Transactional saja
	// nanti kita bisa setting, apakah transactional itu ReadOnly atau tidak
	// untuk parameter kedua adalah transactional, kita gunakan sql.Tx

	// parameter ketiga kita pake data aslinya

	Save(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category
	Update(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category
	Delete(ctx context.Context, tx *sql.Tx, category domain.Category)
	FindById(ctx context.Context, tx *sql.Tx, categoryId int) (domain.Category, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Category
}
