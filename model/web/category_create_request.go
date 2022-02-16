package web

// untuk melakukan Create, kita butuh request apa aja sih, sebenarnya kita cuma butuh name
// kita tambahkan juga validate untuk data name yagn dikirim
// valdate required artinya harus diisi
type CategoryCreateRequest struct {
	Name string `validate:"required,max:200,min=1"`
}
