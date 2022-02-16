package middleware

import (
	"golangrestfulapi/helper"
	"golangrestfulapi/model/web"
	"net/http"
)

// buat struct dari method handler, karena middleware harus dalam bentuk handler

type AuthMiddleware struct {

	// middleware perlu meneruskan requestnya ke handle berikutnya
	Handler http.Handler
}

// function untuk membuat new middleware ini
func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	if "RAHASIA" == request.Header.Get("X-API-Key") {
		// ok
		// teruskan ke handler selanjutnya
		middleware.Handler.ServeHTTP(writer, request)

	} else {
		// error
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}
		helper.WriteToResponseBody(writer, webResponse)
	}

}
