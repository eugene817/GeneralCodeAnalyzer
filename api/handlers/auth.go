package handlers

import (
  "os"
  "time"
  "net/http"
	"github.com/eugene817/GeneralCodeAnalyzer/database"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
  "github.com/golang-jwt/jwt/v5"
)

func (h *Handler) Register (c echo.Context) error {
  type req struct{ Username, Password string }
  var r req
  if err := c.Bind(&r); err != nil {
    return c.JSON(400, echo.Map{"error": "invalid payload"})
  }
  hash, _ := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
  user := database.User{Username: r.Username, PasswordHash: string(hash)}
  if err := h.db.Create(&user).Error; err != nil {
    return c.Redirect(http.StatusSeeOther, "/register?error=User+exists")  
  }
  return c.Render(http.StatusOK, "login_partial", echo.Map{
    "Error": "Now please log in",
  })
}


func (h *Handler) Login (c echo.Context) error {
  type req struct{ Username, Password string }
  var r req
  if err := c.Bind(&r); err != nil {
    return c.JSON(400, echo.Map{"error": "invalid payload"})
  }
  user, err := h.Authenticate(r.Username, r.Password)
  if err != nil {
    return c.Redirect(http.StatusSeeOther, "/login?error=Invalid+credentials")  
  }

  t, err := h.CreateJWT(user); if err != nil {
    return c.Redirect(http.StatusInternalServerError, "/login?error=Internal+Server+Error")
  }
 
  c.SetCookie(&http.Cookie{Name: "jwt", Value: t, Path: "/"})
  return c.Render(http.StatusOK, "index", nil)
}


func (h *Handler) Authenticate(userName, pass string) (database.User, error) {
  var user database.User
  if err := h.db.Where("username = ?", userName).First(&user).Error; err != nil {
    return database.User{}, err
  }
  if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pass)); err != nil {
    return database.User{}, err
  }

  return user, nil
}


func (h *Handler) CreateJWT(user database.User) (string, error) {
   // creating JWT
  token := jwt.New(jwt.SigningMethodHS256)
  claims := token.Claims.(jwt.MapClaims)
  claims["sub"] = user.ID
  claims["name"] = user.Username
  claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

  t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
  if err != nil {
    return "", err
  }
  return t, nil
}
