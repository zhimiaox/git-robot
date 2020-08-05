package main

import "github.com/gin-gonic/gin"

func StartApi() {
	r := gin.Default()
	gr := r.Group("/api")
	gr.POST("/user", Set)
	gr.GET("/user/search", Search)
	gr.DELETE("/user", Del)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func Search(c *gin.Context) {

}

func Set(c *gin.Context) {

}

func Del(c *gin.Context) {

}
