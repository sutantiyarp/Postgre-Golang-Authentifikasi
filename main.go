package main

import (
    "log"

    "github.com/joho/godotenv"

    "hello-fiber/config"
    "hello-fiber/database"
)

func main() {
    // load .env
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env not loaded:", err)
    }

    // NewApp will call ConnectMongoDB internally
    app := config.NewApp()

    // disconnect saat program keluar (DisconnectMongoDB harus aman dipanggil jika belum terhubung)
    defer func() {
        if err := database.DisconnectMongoDB(); err != nil {
            log.Println("Error disconnecting from MongoDB:", err)
        }
    }()

    log.Fatal(app.Listen(":3000"))
}