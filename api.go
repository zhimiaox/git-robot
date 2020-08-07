package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartApi() {
	r := gin.Default()
	r.Static("/html", "html")
	r.StaticFile("/", "html/index.html")
	gr := r.Group("/api")
	gr.POST("/user", Set)
	gr.GET("/user/search", Search)
	gr.DELETE("/user/:id", Del)
	gin.SetMode(gin.DebugMode)
	httpServer := &http.Server{
		Addr:           Config.Server.APIListen,
		Handler:        r,
		ReadTimeout:    time.Duration(Config.Server.ReadTimeOut) * time.Second,
		WriteTimeout:   time.Duration(Config.Server.WriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logrus.Infof("Start HTTP Service Listening %s", Config.Server.APIListen)
	httpServer.ListenAndServe()
}

func Search(c *gin.Context) {
	email, ok := c.GetQuery("email")
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}
	result := make([]*userResp, 0)
	for _, v := range userStorage.List() {
		if v.Email == email {
			result = append(result, v.TOVO())
		}
	}
	c.JSON(http.StatusOK, result)
}

func Set(c *gin.Context) {
	param := &userParam{}
	if err := c.ShouldBind(param); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	user := param.TOUser()
	check := userStorage.Get(user.Sign())
	if check != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}
	err := userStorage.Set(user)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	GitWork(user)
	c.JSON(http.StatusOK, user.Sign())
}

func Del(c *gin.Context) {
	id := c.Param("id")
	if err := userStorage.Del(id); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
