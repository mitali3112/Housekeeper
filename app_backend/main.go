package main

import (
	s "app_backend/controllers"
	m "app_backend/model"
	"fmt"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

var db *gorm.DB
var err error

func init() {
	// Log error if .env file does not exist
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {

	//creating Database using gorm(an ORM which simplifies the mapping and persistance of the models to the database)
	db, err = gorm.Open("sqlite3", "./dbase.db")
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
	r.POST("/seeker_login", s.Login_auth((db)))
	r.POST("/provider_login", s.Login_auth(db))
	r.GET("/seeker_home", nil)
	r.GET("/provider_home", nil)
	r.GET("/services", s.Listing_services(db))
	r.GET("/services/:ServiceId", s.List_service(db))
	//When the seeker tries to book a service, the data has to be updated in the bookings table
	r.POST("/services/:ServiceId/book", s.Book(db))

	// var store = cookie.NewStore([]byte(controllers.RandToken(64)))
	// //Using middleware, store is the storage engine created before and can be replaced by other engines
	// //mysession is the name that will be stored in the cookie on the browser. The server uses this name to find the corresponding session
	// store.Options(sessions.Options{
	// 	Path:   "/",
	// 	MaxAge: 86400 * 7,
	// })
	// r.Use(sessions.Sessions("mysession", store))

	// auth := r.Group("/auth")
	// auth.Use(middleware.Authentication())
	// {
	// 	auth.GET("/test", func(c *gin.Context) {
	// 		c.JSON(200, gin.H{
	// 			"message": "Everything is ok",
	// 		})
	// 	})
	// }
	r.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})

	r.Run(":8080")
}
func Login_auth(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var auth m.Login
		var storedAuth m.Login
		c.BindJSON(&auth)
		err := db.Where("Email = ?", auth.Email).First(&storedAuth).Error
		if err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			match := strings.Compare(auth.Password, storedAuth.Password)
			if match == 0 {
				fmt.Println("match")
				session := sessions.Default(c)
				session.Set("id", 12090292)
				session.Set("email", "test@gmail.com")
				session.Save()
				c.JSON(200, gin.H{"message": "login successful!"})
			} else {
				fmt.Println("No match")
				c.JSON(401, gin.H{"message": "Login Failed!"})
			}

		}
	}

	return gin.HandlerFunc(fn)

}
