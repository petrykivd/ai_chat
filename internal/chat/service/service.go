package service

import (
	chat "ai_chat/internal/chat/store"
	"ai_chat/internal/llm"
	messages "ai_chat/internal/messages/store"
	"ai_chat/pkg"
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateChat(ctx context.Context) (pkg.Chat, error)

	GetChat(ctx context.Context, id uuid.UUID) (pkg.Chat, error)

	ListChats(ctx context.Context) ([]pkg.Chat, error)

	SendMessage(ctx context.Context, chatID uuid.UUID, content string) (llm.Message, error)

	GetMessages(ctx context.Context, chatID uuid.UUID) ([]pkg.Message, error)
}

type ServiceImpl struct {
	chatRepo    chat.Repository
	messageRepo messages.Repository
	llmClient   llm.Client
}

func NewService(chatRepo chat.Repository, messageRepo messages.Repository, llmClient llm.Client) Service {
	return &ServiceImpl{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
		llmClient:   llmClient,
	}
}

func (s *ServiceImpl) CreateChat(ctx context.Context) (pkg.Chat, error) {
	newChat := pkg.Chat{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}

	id, err := s.chatRepo.CreateChat(ctx, newChat)
	if err != nil {
		return pkg.Chat{}, err
	}

	return s.chatRepo.GetChat(ctx, id)
}

func (s *ServiceImpl) GetChat(ctx context.Context, id uuid.UUID) (pkg.Chat, error) {
	return s.chatRepo.GetChat(ctx, id)
}

func (s *ServiceImpl) ListChats(ctx context.Context) ([]pkg.Chat, error) {
	return s.chatRepo.ListChats(ctx)
}

func (s *ServiceImpl) SendMessage(ctx context.Context, chatID uuid.UUID, content string) (llm.Message, error) {
	_, err := s.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		return llm.Message{}, err
	}

	userMessage := pkg.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		Role:      "user",
		Content:   content,
		CreatedAt: time.Now(),
	}

	_, err = s.messageRepo.CreateMessage(ctx, userMessage)
	if err != nil {
		return llm.Message{}, err
	}

	messageHistory, err := s.messageRepo.ListMessages(ctx, chatID)

	convertedMessages := make([]llm.Message, 0)

	for _, m := range messageHistory {
		convertedMessages = append(convertedMessages, s.llmClient.MessageFromModel(m))
	}

	assistantMessage, err := s.llmClient.ChatComplete(ctx, convertedMessages)

	newAssistantMessage := pkg.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		Role:      assistantMessage.Role,
		Content:   assistantMessage.Content,
		CreatedAt: time.Now(),
	}

	_, err = s.messageRepo.CreateMessage(ctx, newAssistantMessage)
	if err != nil {
		return llm.Message{}, err
	}

	return assistantMessage, nil
}

func (s *ServiceImpl) GetMessages(ctx context.Context, chatID uuid.UUID) ([]pkg.Message, error) {
	_, err := s.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return s.messageRepo.ListMessages(ctx, chatID)
}
