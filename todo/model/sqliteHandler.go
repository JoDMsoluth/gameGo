package model

import (
	// sql에 대한 유틸함수 포함
	"database/sql"
	// sqlite pakcage 앞에 '_'는 암시적으로 임시로 사용한다는 뜻
	_ "github.com/mattn/go-sqlite3"

	"time"
)

type sqliteHandler struct {
	db *sql.DB
}

func (s *sqliteHandler) GetTodo(sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionId=?", sessionId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// True가 될떄까지 레코드를 읽어온다.
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
		// todos에 데이터 저장한다.
		todos = append(todos, &todo)
	}
	return todos
}

func (s *sqliteHandler) AddTodo(name string, sessionId string) *Todo {
	// prepare로 이용
	stmt, err := s.db.Prepare("INSERT INTO todos (sessionId, name, completed, createdAt) VALUES (?, ?, ?, datetime('now'))")
	if err != nil {
		panic(err)
	}
	// ?에 들어갈 값들을 넣는다.
	rst, err := stmt.Exec(sessionId, name, false)
	if err != nil {
		panic(err)
	}

	// 마지막으로 추가된 레코드 값의 id를 반환
	id, _ := rst.LastInsertId()
	var todo Todo
	todo.ID = int(id)
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()

	return &todo
}

func (s *sqliteHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=?")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}
	// 쿼리문으로 영향받은 레코드 개수를 반환
	cnt, _ := rst.RowsAffected()
	return cnt > 0
}

func (s *sqliteHandler) CompleteTodo(id int, complete bool) bool{
	stmt, err := s.db.Prepare("UPDATE todos SET completed=? WHERE id=?")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(complete, id)
	if err != nil {
		panic(err)
	}

	// 쿼리문으로 영향받은 레코드 개수를 반환
	cnt, _ := rst.RowsAffected()
	return cnt > 0
}

func (s *sqliteHandler) Close()  {
	s.db.Close()
}

func newSqliteHandler(filepath string) DBHandler{
	database, err := sql.Open("sqlite3", filepath)

	if err != nil {
		panic(err)
	}

	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id	INTEGER	PRIMARY KEY	AUTOINCREMENT,
			sessionId	STRING,
			name	TEXT,
			completed	BOOLEAN,
			createdAt	DATETIME
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (
			sessionId ASC
		);`)
	
	statement.Exec()
	return &sqliteHandler{db : database}
}