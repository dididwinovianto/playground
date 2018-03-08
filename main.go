package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"net/http"
)

var db *gorm.DB

func init()  {
	// open database connection
	var err error
	db, err = gorm.Open("mysql",
		"playground:playground@tcp(192.169.2.5:3306)/playground?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	// migrate the schema
	db.AutoMigrate(&userModel{})
	
}

func main() {
	router := gin.Default()

	v1 := router.Group("api/v1/users")
	{
		v1.POST("/create", createUser)
		v1.GET("/fetch", fetchAllUser)
		v1.GET("/fetch/:id", fetchSingleUser)
		v1.PUT("/update/:id", updateUser)
//		v1.DELETE("/:id", deleteUser)
	}
	router.Run()
}

type(
	userModel struct {
		gorm.Model
		Username	string `json:"username"`
		Completed 	int `	json:"completed"`
	}
	transformedUser struct {
		ID			uint 	`json:"id"`
		Username 	string 	`json:"username"`
		Completed 	bool 	`json:"completed"`
	}
)


// create a new user
func createUser(c *gin.Context)  {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	user := userModel{Username:c.PostForm("username"), Completed: completed}
	db.Save(&user)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Successfully create user", "resourceId": user.ID})
}

// fetch all user
func fetchAllUser(c *gin.Context){
	var user []userModel
	var _user []transformedUser

	db.Find(&user)

	if len(user) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}

	// transforms user for building a good reponse
	for _, item := range user {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_user = append(_user, transformedUser{ID: item.ID, Username: item.Username, Completed: completed})
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _user})
	}
}

// fetch single user
func fetchSingleUser(c *gin.Context)  {
	var user userModel
	userID := c.Param("id")

	db.First(&user, userID)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found"})
		return
	}

	completed := false
	if user.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_user := transformedUser{ID: user.ID, Username: user.Username, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _user})
}

// update user
func updateUser(c *gin.Context)  {
	var user userModel
	userID := c.Param("id")

	db.First(&user, userID)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}

	db.Model(&user).Update("username", c.PostForm("username"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&user).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "User updated successfully"})
}