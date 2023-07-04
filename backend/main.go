package main

import (
	"context"
	"log"

	"todoapp/controllers"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Create a new Fiber app
	app := fiber.New()

	// Set up the MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new MongoDB collection
	todoCollection := client.Database("todoapp").Collection("todos")

	// Create a new TodoController instance
	todoController := controllers.NewTodoController(todoCollection)

	// Define the routes
	app.Get("/api/todos", todoController.GetTodos)
	app.Post("/api/todos", todoController.CreateTodo)
	app.Put("/api/todos/:id", todoController.UpdateTodo)
	app.Delete("/api/todos/:id", todoController.DeleteTodo)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
