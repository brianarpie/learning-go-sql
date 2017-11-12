package main

import (
  "strconv"
  "os"
  "fmt"
  "log"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type User struct {
  Id int
  Name string
}

func fetchUserById(db *sql.DB, id int) (*User, error) {
  var name string
  err := db.QueryRow("SELECT name from users WHERE id = ?", id).Scan(&name)
  if err != nil {
    log.Fatal(err)
    return nil, err
  }
  return &User{Id: id, Name: name}, nil
}

func getUserNameById(db *sql.DB, id int) (string, error) {
  var name string

  stmt, err := db.Prepare("SELECT name from users WHERE id = ?")
  if err != nil {
    log.Fatal(err)
    return "", err
  }
  defer stmt.Close()

  rows, err := stmt.Query(id)
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

func addUser(db *sql.DB, desired_name string) *User {
  var (
    id int
    name string
  )
  stmt, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
  if err != nil {
    log.Fatal(err)
    return nil
  }
  res, err := stmt.Exec(desired_name)
  if err != nil {
    log.Fatal(err)
    return nil
  }
  lastId, err := res.LastInsertId()
  if err != nil {
    log.Fatal(err)
    return nil
  }
  err = db.QueryRow("SELECT id, name FROM users WHERE id = ?", lastId).Scan(&id, &name)
  if err != nil {
    log.Fatal(err)
    return nil
  }

  return &User{Id: id, Name: name}
}

func (u *User) String() string {
  return fmt.Sprintf("Id: %v, Name: %v", u.Id, u.Name)
}

func printError(err error) {
  fmt.Println("SQL GO Error:", err)
}

func main() {
  db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/hello")
  if err != nil {
    log.Fatal(err)
  }

  args := os.Args[1:]
  switch args[0] {
    case "add":
      switch args[1] {
        case "user":
          fmt.Println("User Added", addUser(db, args[2]))
      }
    case "fetch":
      switch args[1] {
        case "user":
          intId, err := strconv.Atoi(args[2])
          if err != nil {
            printError(err)
          }
          user, err := fetchUserById(db, intId)
          if err != nil {
            printError(err)
          }
          fmt.Println("User Fetched", user)
      }
  }
  
  defer db.Close()
}
