package auth

import (
  "testing"
  "net/http"
)

type someHandler struct {}

func(h *someHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestAuthHandler(t *testing.T) {
  MustAuth(&someHandler{})
}