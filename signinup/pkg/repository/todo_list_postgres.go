package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"signinup"
	"strings"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list signinup.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]signinup.TodoList, error) {
	var lists []signinup.TodoList

	query := fmt.Sprintf("SELECT * FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)
	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (signinup.TodoList, error) {
	var list signinup.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl "+
		"INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = list",
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId)
	return list, err
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2", todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input signinup.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$id", argsId))
		argsId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$id", argsId))
		argsId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argsId, argsId+1)
	args = append(args, listId, userId)

	logrus.Debug("updateQuery: %s", query)
	logrus.Debug("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
