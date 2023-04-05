package cors

import (
	"net/http"
	"regexp"
	"userService/services/api-gateway/config"
)

func Cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func allowedOrigin(origin string) bool {
	conf := config.PublicConfig()
	if conf.Cors == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(conf.Cors, origin); matched {
		return true
	}
	return false
}
