package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User DB
type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func main() {
	app := fiber.New()
	app.Use(middleware.Logger(middleware.LoggerConfig{
		Format:     " ${time} ${method} ${path}",
		TimeFormat: " 15:04:05",
		TimeZone:   " Asia/Bangkok",
		Output:     os.Stdout,
	}))

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Hello, World! Fiber")
	})

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	fmt.Println("Fiber served on http://localhost:3000/api/")

	collection := client.Database("go_mongodb").Collection("users")

	app.Get("/api/user", func(c *fiber.Ctx) {

		findOptions := options.Find()
		// findOptions.SetLimit(2)
		var results []*User

		cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
		if err != nil {
			log.Fatal(err)
		}

		for cur.Next(context.TODO()) {
			var elem User
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, &elem)

			if err := cur.Err(); err != nil {
				log.Println(err)
			}

			cur.Close(context.TODO())
			c.JSON(results)
		}
	})
	app.Post("/api/user/add", func(c *fiber.Ctx) {
		// body := c.Body()
		// fmt.Println(body)

		p := new(User)
		if err := c.BodyParser(p); err != nil {
			log.Fatal(err)
		}

		dummy := User{p.Username, p.Password}
		insertResult, err := collection.InsertOne(context.TODO(), dummy)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
		c.JSON(insertResult)

	})
	app.Get("/api/user/:username", func(c *fiber.Ctx) {
		log.Println("Get request with value: " + c.Params("username"))
		var result User
		filter := bson.D{{"username", c.Params("username")}}

		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			log.Println(err)
		}
		c.JSON(result)
	})
	app.Put("/api/user", func(c *fiber.Ctx) {
		p := new(User)

		if err := c.BodyParser(p); err != nil {
			log.Println(err)
		}

		filter := bson.D{{"username", p.Username}}

		update := bson.D{
			{"$set", bson.D{
				{"Password", p.Password},
			}},
		}
		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println(err)
		}

		c.JSON(updateResult)
	})
	app.Delete("api/user/:username", func(c *fiber.Ctx) {
		log.Println("Get request with value: " + c.Params("username"))

		filter := bson.D{{"username", c.Params("username")}}
		deleteResult, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(deleteResult)

	})

	app.Listen(":3000")

}
