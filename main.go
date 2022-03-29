package main

import (
	"context"
	"strconv"

	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func toNumberFromString(start string) int {
    val, err := strconv.Atoi(start)
    if err != nil {
			panic(err)
	}
    return val
}

func main() {
    app := fiber.New()
    app.Use(cors.New())
    godotenv.Load()
	rdb := redis.NewClient(&redis.Options{
        Addr: os.Getenv("DB_URL"),
        Password: os.Getenv("DB_PASSWORD"), // no password set
        DB:       0,  // use default DB
    })

    app.Get("/clap/:article/:claps", func(c *fiber.Ctx) error {
        val, err := rdb.Get(ctx, c.Params("article")).Result()
        if err != nil {
            if (err.Error() == "redis: nil") {
                val = "0";
            } else {
                panic(err)
            }
        }
		err2 := rdb.Set(ctx, c.Params("article"), toNumberFromString(val) + toNumberFromString(c.Params("claps")), 0).Err()
		if err2 != nil {
			panic(err2)
		}
        return c.SendString("👏")
    })

    app.Get("/get-claps/:article", func(c *fiber.Ctx) error {
        val, err := rdb.Get(ctx, c.Params("article")).Result()
        if err != nil {
            if (err.Error() == "redis: nil") {
                return c.SendString("0")
            } else {
                panic(err)
            }
        }
        return c.SendString(val)
    })

    app.Listen(":3000")
}