package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-contrib/multitemplate"

	"gopkg.in/gin-gonic/gin.v1"
	redis "gopkg.in/redis.v5"
)

type Register struct {
	UserName             string `form:"username" binding:"required"`
	Email                string `form:"email"    binding:"required"`
	Password             string `form:"password" binding:"required"`
	PasswordConfirmation string `form:"password_confirmation" binding:"required"`
}

// Setup global variable for redis connection
// Needs a better abstraction later, maybe
var db *redis.Client

// Wrapper functions to create secure random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// Load nested templates
func createMyRender() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("index", "templates/layout.tmpl", "templates/index.tmpl")
	r.AddFromFiles("register", "templates/layout.tmpl", "templates/register.tmpl")
	r.AddFromFiles("register_success", "templates/layout.tmpl", "templates/register_success.tmpl")

	return r
}

// Handler definitions
func index_handler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{})
}

func register_handler(c *gin.Context) {
	cookie, _ := c.Cookie("id")

	val, err := db.Get(cookie).Result()
	if err == redis.Nil {
		fmt.Println("key does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key", val)
	}

	c.HTML(http.StatusOK, "register", gin.H{"cookie": cookie})
}

func register_submit_handler(c *gin.Context) {
	var form Register
	if c.Bind(&form) == nil {
		random_bytes, _ := GenerateRandomString(128)

		err := db.Set(random_bytes, "hukl", 0).Err()
		if err != nil {
			panic(err)
		}

		c.SetCookie("id", random_bytes, 3600, "/", "", false, false)

		c.HTML(http.StatusOK, "register_success", gin.H{
			"username": form.UserName,
			"email":    form.Email,
			"password": form.Password,
		})
	}
}

func main() {
	db = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := db.Ping().Result()
	fmt.Println(pong, err)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.HTMLRender = createMyRender()

	router.Static("/assets", "./assets")

	router.GET("/", index_handler)
	router.GET("/register", register_handler)
	router.POST("/register", register_submit_handler)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
