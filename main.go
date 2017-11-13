package main

import (
  "time"
  "strconv"
  "os"
  "fmt"
)

type User struct {
  Id int
  Name string
}

func (u *User) String() string {
  return fmt.Sprintf("Id: %v, Name: %v", u.Id, u.Name)
}

func printError(err error) {
  dateTimeFormat := "2006/01/02 15:04:05"
  fmt.Printf("%v [SQL GO] Error: %v\n", time.Now().Format(dateTimeFormat), err)
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