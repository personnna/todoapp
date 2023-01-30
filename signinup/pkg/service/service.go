package service

import (
	"signinup"
	"signinup/pkg/repository"
)

type Authorization interface {
	CreateUser(user signinup.User) (int, error)
	GenerateToken(user, password string) (string, error)
	ParseToken(token string) (int, error)
}

type ToDoList interface {
	Create(userId int, list signinup.TodoList) (int, error)
	GetAll(userId int) ([]signinup.TodoList, error)
	GetById(userId, listId int) (signinup.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input signinup.UpdateListInput) error
}

type ToDoItem interface {
	Create(userId, listId int, item signinup.TodoItem) (int, error)
	GetAll(userId, listId int) ([]signinup.TodoItem, error)
	GetById(userId, itemId int) (signinup.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input signinup.UpdateItemInput) error
}

type Service struct {
	Authorization
	ToDoList
	ToDoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		ToDoList:      NewTodoListService(repos.ToDoList),
		ToDoItem:      NewTodoItemService(repos.ToDoItem, repos.ToDoList),
	}
}
