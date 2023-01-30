package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"signinup"
	"strings"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item signinup.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RE", listsItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2) RETURNING id", todoItemTable)
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]signinup.TodoItem, error) {
	var items []signinup.TodoItem
	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id WHERE li.list_id = $1 AND ul.user_id = $2",
		todoItemTable, listsItemsTable, usersListsTable)
	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (signinup.TodoItem, error) {
	var item signinup.TodoItem
	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id WHERE ti.item_id = $1 AND ul.user_id = $2",
		todoItemTable, listsItemsTable, usersListsTable)
	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}
	return item, nil
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s ti USING %s li, %s ul"+
		"WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2",
		todoItemTable, listsItemsTable, usersListsTable)

	_, err := r.db.Exec(query, userId, itemId)
	return err
}

func (r *TodoItemPostgres) Update(userId, itemId int, input signinup.UpdateItemInput) error {
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
	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$id", argsId))
		argsId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d",
		todoListsTable, setQuery, listsItemsTable, usersListsTable, argsId, argsId+1)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}
