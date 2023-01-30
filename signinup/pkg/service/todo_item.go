package service

import (
	"signinup"
	"signinup/pkg/repository"
)

type TodoItemService struct {
	repo     repository.ToDoItem
	listRepo repository.ToDoList
}

func NewTodoItemService(repo repository.ToDoItem, listRepo repository.ToDoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userId, listId int, item signinup.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(listId, item)
}
func (s *TodoItemService) GetAll(userId, listId int) ([]signinup.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, itemId int) (signinup.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId, itemId int, input signinup.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}
