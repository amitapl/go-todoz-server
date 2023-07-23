package main

type TodoItem struct {
    Text string	`json:"text"`
    Done  bool	`json:"done"`
	Id string	`json:"id"`
}
