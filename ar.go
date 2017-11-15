package main

import (
  "log"
  "fmt"
  "database/sql"
)

type relation struct {
  db *sql.DB
  query string
  values []int
}

func (rel *relation) Sel(table, column string) *relation {
  rel.query = fmt.Sprintf("SELECT %v FROM %v ", column, table)
  return rel
}

func (rel *relation) Where(column string, value int) *relation {
  // brute force insert "WHERE column = value"
  rel_count := len(rel.values) + 1
  rel.query = rel.query + fmt.Sprintf("WHERE %v = $%v", column, rel_count)
  rel.values = append(rel.values, value)
  return rel
}

func (rel *relation) Exec() (string, error) {
  // TODO: dynamically pass the values
  var value string
  err := rel.db.QueryRow(rel.query, rel.values[0]).Scan(&value)
  if err != nil {
    log.Fatal(err)
    return "", err
  }
  return value, nil
}
