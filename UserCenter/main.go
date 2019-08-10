package main

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"log"

	"github.com/gin-gonic/gin"

	"./controllers"
	_ "./models"
	"./utils"
)

func init() {
	log.SetPrefix("AuthCenter ")

}

func main() {
	r := gin.Default()
	auth := r.Group("/auth")
	auth.POST("/user", controllers.SignUp)
	auth.POST("profile", controllers.SignIn)

	info := r.Group("/info", controllers.JWTAuthMiddleware())
	info.GET("/profile", controllers.GetProfile)
	info.PUT("/profile", controllers.PutProfile)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("nameValidator", utils.NameValidator)
		_ = v.RegisterValidation("emailValidator", utils.EmailValidator)
		_ = v.RegisterValidation("mobileValidator", utils.MobileValidator)
		_ = v.RegisterValidation("passwordValidator", utils.PasswordValidator)
		_ = v.RegisterValidation("genderValidator", utils.GenderValidator)
	}

	err := r.Run(":8080")
	if nil != err {
		log.Fatal(err)
	}
}
