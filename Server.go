package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func rootHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hi there, I love %s!", c.Request.URL)
	log.Printf("Go %s", c.Request.URL)
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func postTodoList(store *Store) func(c *gin.Context) {
	return func(c *gin.Context) {
		list := TodoList{}

		if err := c.BindJSON(&list); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if list.Id == "" || len(list.Id) != 8 {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid id, should have 8 characters"))
			return
		}

		err := store.putTodoList(list)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.JSON(http.StatusOK, &list)

		log.Printf("List posted %s", list.Id)
	}
}

func getTodoList(store *Store) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" || len(id) != 8 {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid id, should have 8 characters"))
			return
		}

		list, err := store.getTodoList(id)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.JSON(http.StatusOK, &list)

		log.Printf("List get %s", list.Id)
	}
}

func startServer(store *Store) {
	r := gin.Default()

	r.GET("/", rootHandler)
	r.GET("/hello", helloHandler)
	r.POST("/api/todoList", postTodoList(store))
	r.GET("/api/todoList/:id", getTodoList(store))

	fmt.Printf("Starting server at port 8080\n")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func main() {
	store := Store{}
	err := store.setupAws()
	if err != nil {
		log.Fatal(err)
		return
	}

	startServer(&store)
}
