package main

import (
  "time"
  "strconv"
  "os"
  "fmt"
  "log"
  "database/sql"
  _ "github.com/lib/pq"
)

type User struct {
  Id int
  Name string
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

func (u *User) String() string {
  return fmt.Sprintf("Id: %v, Name: %v", u.Id, u.Name)
}

func printError(err error) {
  dateTimeFormat := "2006/01/02 15:04:05"
  fmt.Printf("%v [SQL GO] Error: %v\n", time.Now().Format(dateTimeFormat), err)
}

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

func main() {
  db := initDatabase()

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
