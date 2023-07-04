package controllers

import (
	"context"
	"log"
	"net/http"
	m "todoapp/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type TodoController struct {
	todoCollection *mongo.Collection
}

func NewTodoController(collection *mongo.Collection) *TodoController {
	return &TodoController{
		todoCollection: collection,
	}
}

func (tc *TodoController) GetTodos(c *fiber.Ctx) error {
	// Fetch all todos from the collection
	ctx := context.Background()
	cursor, err := tc.todoCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error fetching todos:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Error fetching todos"})
	}
	defer cursor.Close(ctx)

	// Decode the cursor to a slice of Todo objects
	var todos []m.Todo
	if err := cursor.All(ctx, &todos); err != nil {
		log.Println("Error decoding todos:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Error decoding todos"})
	}

	return c.JSON(todos)
}

func (tc *TodoController) CreateTodo(c *fiber.Ctx) error {
	// Parse request body to Todo object
	todo := new(m.Todo)
	if err := c.BodyParser(todo); err != nil {
		log.Println("Error parsing request:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Insert the todo into the collection
	ctx := context.Background()
	result, err := tc.todoCollection.InsertOne(ctx, todo)
	if err != nil {
		log.Println("Error creating todo:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Error creating todo"})
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(http.StatusCreated).JSON(todo)
}

func (tc *TodoController) UpdateTodo(c *fiber.Ctx) error {
	// Parse request body to Todo object
	todo := new(m.Todo)
	if err := c.BodyParser(todo); err != nil {
		log.Println("Error parsing request:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Get the ID from the request URL params
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		log.Println("Error parsing ID:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid ID"})
	}

	// Update the todo in the collection
	ctx := context.Background()
	update := bson.M{"$set": todo}
	_, err = tc.todoCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		log.Println("Error updating todo:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating todo"})
	}

	return c.JSON(todo)
}

func (tc *TodoController) DeleteTodo(c *fiber.Ctx) error {
	// Get the ID from the request URL params
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		log.Println("Error parsing ID:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid ID"})
	}

	// Delete the todo from the collection
	ctx := context.Background()
	_, err = tc.todoCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Println("Error deleting todo:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Error deleting todo"})
	}

	return c.SendStatus(http.StatusNoContent)
}
