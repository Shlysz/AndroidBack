package main

import (
	"github.com/gin-gonic/gin"
	"src/androidBackground/handler"
	"src/androidBackground/respo"
)

func main() {

	Server := gin.Default()
	Server.Use(gin.Logger())
	//初始化数据库

	Server.POST("/login", loginHandler)
	Server.POST("/chat", handler.HandleChat)
	err := Server.Run(":8080")
	if err != nil {
		panic(err)
	}

}

func loginHandler(mygin *gin.Context) {
	//获取post的内容
	username := mygin.PostForm("username")
	password := mygin.PostForm("password")
	account := respo.Account{Username: username, Password: password}
	if account.Login() {
		mygin.JSON(200, gin.H{
			"message": "success",
			"key":     account.Key,
		})
	} else {
		mygin.JSON(200, gin.H{
			"message": "fail",
		})
	}

}
