package main

import (
  "os"
  "fmt"
  "log"
  "database/sql"
  _ "github.com/lib/pq"
)

func initDatabase() *sql.DB {
  username := os.Getenv("DB_USERNAME")
  password := os.Getenv("DB_PASSWORD")
  database := os.Getenv("DB_NAME")

  database_url := fmt.Sprintf("postgres://%v:%v@127.0.0.1:5432/%v?sslmode=disable", username, password, database)

  db, err := sql.Open("postgres", database_url)
  if err != nil {
    log.Fatal(err)
  }

  return db
}

func fetchUserById(db *sql.DB, id int) (*User, error) {
  var name string
  err := db.QueryRow("SELECT name from users WHERE id = $1", id).Scan(&name)
  if err != nil {
    log.Fatal(err)
    return nil, err
  }
  return &User{Id: id, Name: name}, nil
}

func fetchUserByName(db *sql.DB, name string) (*User, error) {
  var id int
  err := db.QueryRow("SELECT id FROM users WHERE name = $1", name).Scan(&id)
  if err != nil {
    log.Fatal(err)
    return nil, err
  }
  return &User{Id: id, Name: name}, nil
}

func addUser(db *sql.DB, desired_name string) *User {
  var (
    id int
    name string
  )
  stmt, err := db.Prepare("INSERT INTO users(name) VALUES($1) RETURNING id, name")
  if err != nil {
    log.Fatal(err)
    return nil
  }
  err = stmt.QueryRow(desired_name).Scan(&id, &name)
  if err != nil {
    log.Fatal(err)
    return nil
  }
  return &User{Id: id, Name: name}
}
