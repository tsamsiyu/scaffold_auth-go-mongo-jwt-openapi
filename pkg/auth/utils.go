package auth

import "net/http"

func ExtractTokenFromHttpHeaders(header http.Header) string {
	auth := header.Get("Authorization")

	return auth
}
