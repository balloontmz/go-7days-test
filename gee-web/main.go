package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/balloontmz/gee-web/gee"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Fail(500, "some error")
		log.Printf("[%d] %s in %v only for v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	e := gee.New()
	e.Use(gee.Logger())
	e.Use(gee.Recovery())
	e.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	e.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello world")
	})

	e.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s , you're at %s\n", c.Param("name"), c.Path)
	})

	e.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	e.GET("/error", func(c *gee.Context) {
		name := []string{}
		_ = name[100]
		c.String(http.StatusOK, "can't be here")
	})

	v2 := e.Group("/v2")
	// v2.Use(onlyForV2())
	v2.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello v2")
	})
	v2.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, fmt.Sprintf("hello %s", c.Param("name")))
	})

	e.Run(":9999")
}
