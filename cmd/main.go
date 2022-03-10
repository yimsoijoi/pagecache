package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yimsoijoi/pagecache/handler"
	"github.com/yimsoijoi/pagecache/redis"
)

func main() {
	rdb := redis.New()
	_handler := handler.New(rdb)
	app := fiber.New()

	app.Post("/", _handler.Handle)
	log.Fatal(app.Listen(":8000"))
}
