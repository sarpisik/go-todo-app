package todo

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type TodoController struct {
	repository *TodoRepository
}

func (controller *TodoController) GetAll(c *fiber.Ctx) error {
	var todos []Todo = controller.repository.FindAll()

	return c.JSON(todos)
}

func (controller *TodoController) GetOne(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	todo, err := controller.repository.FindOne(id)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(todo)
}

func (controller *TodoController) CreateOne(c *fiber.Ctx) error {
	var data = new(Todo)

	if err := c.BodyParser(data); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid input",
			"error":   err.Error(),
		})
	}

	todo, err := controller.repository.CreateOne(*data)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to create a new todo",
			"error":   err,
		})
	}

	return c.JSON(todo)
}

func (controller *TodoController) UpdateOne(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	todo, err := controller.repository.FindOne(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": err,
		})
	}

	todoData := new(Todo)
	if err := c.BodyParser(todoData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to create a new todo",
			"error":   err,
		})
	}

	todo.Name = todoData.Name
	todo.Description = todoData.Description
	todo.Status = todoData.Status

	updatedTodo, err := controller.repository.UpdateOne(todo)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update the todo",
			"error":   err,
		})
	}

	return c.JSON(updatedTodo)
}

func (controller *TodoController) DeleteOne(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	rowsAffected := controller.repository.DeleteOne(id)
	statusCode := http.StatusNoContent
	if rowsAffected == 0 {
		statusCode = http.StatusBadRequest
	}

	return c.Status(statusCode).JSON(nil)
}

func NewTodoController(repository *TodoRepository) *TodoController {
	return &TodoController{
		repository: repository,
	}
}

func Register(router fiber.Router, db *gorm.DB) {
	db.AutoMigrate(&Todo{})
	todoRepository := NewTodoRepository(db)
	todoController := NewTodoController(todoRepository)

	todoRouter := router.Group("/todos")
	todoRouter.Get("/", todoController.GetAll)
	todoRouter.Get("/:id", todoController.GetOne)
	todoRouter.Post("/", todoController.CreateOne)
	todoRouter.Put("/:id", todoController.UpdateOne)
	todoRouter.Delete("/:id", todoController.DeleteOne)
}
