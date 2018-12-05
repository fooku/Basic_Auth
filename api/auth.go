package api

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fooku/authBasic/models"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
)

type Ur struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

func Test(c echo.Context) error {
	u := &Ur{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.JSON(http.StatusOK, u)
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	err, user := models.FindUser(username)

	if err == mgo.ErrNotFound {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid email"}
	}

	if user.ComparePassword(password) {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Jon Snow"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid password"}
}

type UserRequest struct {
	Username string `json:"username" form:"username" query:"username"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}

func Register(c echo.Context) error {
	var user models.User

	u := new(UserRequest)
	err := c.Bind(u)

	fmt.Println(u)
	err = models.AddUser(user, u.Username, u.Email, u.Password)

	return err
}
