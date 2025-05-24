package main

import (
	"ai_chat/internal/chat/service"
	chat "ai_chat/internal/chat/store"
	"ai_chat/internal/config"
	"ai_chat/internal/llm"
	messages "ai_chat/internal/messages/store"
	"ai_chat/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MessageRequest represents the structure of the request for sending a new message
type MessageRequest struct {
	Content string `json:"content"`
}

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Add CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Configure server to handle JSON properly
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.Printf("Request Body: %s", reqBody)
	}))

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&pkg.Chat{}, &pkg.Message{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	client := llm.NewAnthropicClient(config.Anthropic)
	chatRepo := chat.NewPostgresChatRepository(db)
	messageRepo := messages.NewPostgresMessageRepository(db)
	chatService := service.NewService(chatRepo, messageRepo, client)

	// Create a new chat
	e.POST("/chats", func(c echo.Context) error {
		newChat, err := chatService.CreateChat(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusCreated, newChat)
	})

	// Get list of all chats
	e.GET("/chats", func(c echo.Context) error {
		chats, err := chatService.ListChats(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, chats)
	})

	// Get messages for a specific chat
	e.GET("/chats/:chat_id/messages", func(c echo.Context) error {
		chatIDStr := c.Param("chat_id")
		chatID, err := uuid.Parse(chatIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid chat ID format",
			})
		}

		messages, err := chatService.GetMessages(c.Request().Context(), chatID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, messages)
	})

	// Send a new message in a specific chat - improved version
	e.POST("/chats/:chat_id/new_message", func(c echo.Context) error {
		chatIDStr := c.Param("chat_id")
		chatID, err := uuid.Parse(chatIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid chat ID format",
			})
		}

		// Parse request body with explicit error handling
		var req MessageRequest
		if err := c.Bind(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("Invalid request format: %v", err),
			})
		}

		log.Printf("Received message request: %+v", req)

		// Validate message content
		if req.Content == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Message content cannot be empty",
			})
		}

		// Send message to the chat service
		assistantMessage, err := chatService.SendMessage(c.Request().Context(), chatID, req.Content)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, assistantMessage)
	})

	log.Fatal(e.Start(fmt.Sprintf(":%d", config.Server.Port)))
}
