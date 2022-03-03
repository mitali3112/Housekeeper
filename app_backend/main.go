package main

import (
	s "app_backend/controllers"
	m "app_backend/model"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

func main() {

	//creating Database using gorm(an ORM which simplifies the mapping and persistance of the models to the database)
	db, err = gorm.Open("sqlite3", "./backend.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&m.Seeker{})
	db.AutoMigrate(&m.Login{})
	db.AutoMigrate(&m.Booking{})
	db.AutoMigrate(&m.ServiceAndProvider{})

	//creating variable using gin Web Framework to handle routing and serving HTTP requests
	//r :=gin.Default() does not work, it gives a huge runtime error
	r := gin.New()

	r.GET("/", s.Home(db))
	r.POST("/seeker_registration", s.Create_seeker(db))
	r.POST("/service_registration", s.Create_service(db))
	r.POST("/seeker_login", s.Login_auth(db))
	r.POST("/provider_login", s.Login_auth(db))
	r.GET("/seeker_home", nil)
	r.GET("/provider_home", nil)
	r.GET("/services", s.Listing_services(db))
	r.GET("/services/:ServiceId", s.List_service(db))
	//When the seeker tries to book a service, the data has to be updated in the bookings table
	r.POST("/services/:ServiceId/book", s.Book(db))

	var store = cookie.NewStore([]byte("something-very-secret"))
	//Using middleware, store is the storage engine created before and can be replaced by other engines
	//mysession is the name that will be stored in the cookie on the browser. The server uses this name to find the corresponding session
	r.Use(sessions.Sessions("mysession", store))
	fmt.Println(store)

	r.GET("/hello", func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("hello") != "world" {
			session.Set("hello", "world")
			session.Save()
		}

		c.JSON(200, gin.H{"hello": session.Get("hello")})
	})

	r.Run(":8080")
}
