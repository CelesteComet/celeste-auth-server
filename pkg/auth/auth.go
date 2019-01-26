package auth

import (
  "fmt"
  "net/http"
  "github.com/dgrijalva/jwt-go"
)

type AuthHandler struct {
  next http.Handler
}

type AdminHandler struct {
  next http.Handler  
}

type LogOutHandler struct {}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("JWT")
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }


  // If Cookie exists, check the JWT
  tokenString := cookie.Value

  // Get parsed token 
  claims := jwt.MapClaims{}
  vToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
    _, ok := token.Method.(*jwt.SigningMethodHMAC)
    if !ok {
      return nil, nil
    }
    return []byte("secret"), nil
  })
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }

  // Validate token
  if !vToken.Valid {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  } 

  // do something with decoded claims
  for key, val := range claims {
    fmt.Printf("Key: %v, value: %v\n", key, val)
  }

  // Set HttpOnly To Prevent Future Tampering
  http.SetCookie(w, &http.Cookie{
    Name:   "JWT",
    Value: tokenString,
    HttpOnly: true,
    Path: "/",
  })

  h.next.ServeHTTP(w, r)
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (h *LogOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  http.SetCookie(w, &http.Cookie{
    Name:   "JWT",
    Value:  "",
    Path:   "/",
    MaxAge: -1,
  })
}

func MustAuth(handler http.Handler) http.Handler {
  return &AuthHandler{next: handler}
}

func IsAdmin(handler http.Handler) http.Handler {
  return &AdminHandler{next: handler} 
}

