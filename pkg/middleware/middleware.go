package middleware

import (
  "net/http"
  "log"
)

func Cors(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Println("PASSING THROUGH CORS MIDDLEWARE")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Expose-Headers", "JWT")
    w.Header().Set("Access-Control-Expose-Headers", "Jwt")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
    if (r.Method == "OPTIONS") {
      w.WriteHeader(http.StatusOK)
      return
    } else {
      next.ServeHTTP(w, r)
    }
  })
}