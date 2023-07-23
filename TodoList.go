package main

import (
	"time"
)

type TodoList struct {
    TodoItems []TodoItem	`json:"todoItems"`
	Name string	            `json:"name"`
    Id string           	`json:"id"`
    Notes string        	`json:"notes"`
    CreationDate time.Time	`json:"creationDate"`
}
