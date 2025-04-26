package config


import (
  "fmt"
  "os"
  "github.com/joho/godotenv"
)


func init() {
  err := godotenv.Load()
  if err != nil {
    fmt.Println("No .env file found, using default settings")
  }
}


func GetPort() string{
  port := os.Getenv("PORT")
  if port == "" {
    port = "8081"
  }

  return ":" + port
}
