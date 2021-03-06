package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User model
type User struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	HasPassword
	Email     string `bson:"email"`
	Username  string `bson:"username"`
	FirstName string `bson:"firsname"`
	LasttName string `bson:"lastname"`
	Timestamp time.Time
}

// AddUser > insert data to mongoDB
func AddUser(user User, u string, e string, p string) error {
	s := mongoSession.Copy()
	defer s.Close()

	c := s.DB(database).C("users")

	index := mgo.Index{
		Key:    []string{"username", "email"},
		Unique: true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	user.SetPassword(p)
	user.Username = u
	user.Email = e
	user.Timestamp = time.Now()
	err = c.Insert(&user)
	// err := s.DB(database).C("users").Update("Username", &user)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: err}
	}
	return nil
}

// FindUser > ค้นหาUser จาก User ID
func FindUser(username string) (User, error) {
	s := mongoSession.Copy()
	defer s.Close()

	var result User
	err := s.DB(database).C("users").Find(bson.M{"username": username}).One(&result)
	// Find(bson.M{"username": "sdfasdfasdf"}).Sort("-timestamp").All(&results)

	if err != nil {
		return result, err
	}

	fmt.Println("Results All: ", result.HasPassword)

	return result, nil
}
