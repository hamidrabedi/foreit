package main

import (
    "fmt"
    "github.com/forgego/forge/db"
)

func main() {
    var dbNil *db.DB = nil

    // Check which ones panic
    fmt.Println("Testing DB methods on nil receiver")

    func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Dialect Panic caught:", r)
            }
        }()
        dbNil.Dialect()
    }()

    func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("PoolStats Panic caught:", r)
            }
        }()
        dbNil.PoolStats()
    }()

    func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("RebindPlaceholders Panic caught:", r)
            }
        }()
        dbNil.RebindPlaceholders("SELECT * FROM users")
    }()

    func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Ping Panic caught:", r)
            }
        }()
        dbNil.Ping(nil)
    }()

    func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("IsConnected Panic caught:", r)
            }
        }()
        dbNil.IsConnected()
    }()
}
