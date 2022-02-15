package domain

// setiap table di database perlu kita representasikan dalam data struct
// di table category hanya ada kolom id dan name

type Category struct {
	Id   int
	Name string
}
