package main

import (
  "fmt"
  "log"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

func getUserNameById(db *sql.DB, id int) (string, error) {
  var name string

  rows, err := db.Query("SELECT name from users WHERE id = ?", id)
  if err != nil {
    log.Fatal(err)
    return "", err
  }

  defer rows.Close()

  for rows.Next() {
    err := rows.Scan(&name)
    if err != nil {
      log.Fatal(err)
      return "", err
    }
    log.Println("users", id, name)
  }

  err = rows.Err()
  if err != nil {
    log.Fatal(err)
    return "", err
  }

  return name, nil
}

func main() {
  db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/hello")

  if err != nil {
    log.Fatal(err)
  }

  err = db.Ping()
  if err != nil {
    log.Fatal(err)
  }

  name, err := getUserNameById(db, 1)
  fmt.Printf("User Name Fetched: %s\n", name)

  defer db.Close()
}
