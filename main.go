package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fooku/authBasic/api"
	"github.com/fooku/authBasic/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/bcrypt"
)

const (
	mongoURL = "mongodb://banana:banana1234@ds123834.mlab.com:23834/video_online"
)

type User struct {
	Name string `json:"name" xml:"name"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash is xx
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	u := &User{
		Name: name,
	}
	fmt.Println(user, claims["admin"])
	return c.JSON(http.StatusOK, u)
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	err := models.Init(mongoURL)
	if err != nil {
		log.Fatalf("can not init model; %v", err)
	}

	// password := "secret"
	// hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	// fmt.Println("Password:", password)
	// fmt.Println("Hash:    ", hash)

	// match := CheckPasswordHash(password, hash)
	// fmt.Println("Match:   ", match)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Login route
	e.POST("/login", api.Login)
	e.POST("/register", api.Register)
	e.GET("/test", api.Test)
	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":" + port))
}
