package main

import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
	)
	
func defaultHandler(c *gin.Context){
	c.HTML(http.StatusOK,"default.html",gin.H{})
}

func setupRouter(r *gin.Engine){
	r.LoadHTMLGlob("templates/**/*.html");
	r.GET("/",defaultHandler)
}

func main(){
	r := gin.Default()
	setupRouter(r)
	err := r.Run(":8000")
	if err != nil{
		log.Fatalf("gin Run error: %s", err)
	}
}