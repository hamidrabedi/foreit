package main

import (
    "fmt"
    "github.com/forgego/forge/db"
    "context"
)

func main() {
    fmt.Println("Testing db tx...")
    var dbNil *db.DB = nil

    // Test if calling BeginTx on a nil DB crashes
    _, err := dbNil.BeginTx(context.Background(), nil)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("No error on BeginTx")
    }

    // Test WithTx
    err = dbNil.WithTx(context.Background(), func(tx *db.Tx) error { return nil })
    if err != nil {
        fmt.Println("Error WithTx:", err)
    } else {
        fmt.Println("No error WithTx")
    }
}
