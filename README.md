# Membangun Aplikasi Berbasis Restful API menggunakan Go-Lang

# Agenda

- Praktik Membuat Aplikasi RESTful API
- Membuat OpenAPI
- Setup Database
- Membuat Model
- Membuat Repository
- Membuat Service
- Membuat Validation
- Membuat Controller
- dan lain-lain

# Aplikasi CRUD Sederhana

- Kita akan membuat aplikasi CRUD sederhana
- Tujuannya untuk belajar RESTful API, bukan untuk membuat aplikasi
- Kita akan membuat aplikasi CRUD untuk data category
- Dimana data category memiliki atribut `id(number)` dan `name(string)`
- Kita akan buat API untuk semua operasinya, Create Category, Get Category, List Category, Update Category, dan Delete Category
- Semua API akan kita tambahkan Authentication berupa API-Key

# Menambah Dependency

- Driver MySQL : https://github.com/go-sql-driver/mysql
- HTTP Router : https://github.com/julienschmidt/httprouter
- Validation : https://github.com/go-playground/validator

# Database

- Database Name : golang_restful_api
- Database user : root
- Database password : root

# Category Domain

- adalah model untuk database
- atau representasi semua data yang ada di database
- kadang dibilang adalah data untuk request dan response
- disebut juga sebagai entity

# Category Repository

- adalah data access layer nya ke domain category
- kita menggunakan Repository Pattern

# Category Repository

- kita implementasi kontrak dari Repository nya
