package crawlerapi

import (
	"fmt"
	"net/http"

	"github.com/ZakkBob/AskDave/backend/orm"
	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Init() {
	router = gin.Default()
	router.GET("/tasks", GetTasks)
	router.POST("/results", PostResults)
	router.Run() // listen and serve on 0.0.0.0:8080
}

func GetTasks(c *gin.Context) {
	t, _ := orm.NextTasks(10)

	c.IndentedJSON(http.StatusOK, t)
}

func PostResults(c *gin.Context) {
	var results tasks.Results

	if err := c.BindJSON(&results); err != nil {
		return
	}

	err := orm.SaveResults(&results)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v")
}
