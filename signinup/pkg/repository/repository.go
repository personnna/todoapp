package repository

import (
	"github.com/jmoiron/sqlx"
	"signinup"
)

type Authorization interface {
	CreateUser(user signinup.User) (int, error)
	GetUser(username, password string) (signinup.User, error)
}

type ToDoList interface {
	Create(userId int, list signinup.TodoList) (int, error)
	GetAll(userId int) ([]signinup.TodoList, error)
	GetById(userId, listId int) (signinup.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input signinup.UpdateListInput) error
}

type ToDoItem interface {
	Create(listId int, item signinup.TodoItem) (int, error)
	GetAll(userId, listId int) ([]signinup.TodoItem, error)
	GetById(userId, itemId int) (signinup.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input signinup.UpdateItemInput) error
}

type Repository struct {
	Authorization
	ToDoList
	ToDoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		ToDoList:      NewTodoListPostgres(db),
		ToDoItem:      NewTodoItemPostgres(db),
	}
}
