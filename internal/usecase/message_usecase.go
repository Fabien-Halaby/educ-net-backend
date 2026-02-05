package usecase

import (
	"context"
	"educnet/internal/domain"
	"educnet/internal/repository"
	"fmt"
)

type MessageUseCase interface {
	SendMessage(ctx context.Context, userID, classID int, content string) (domain.Message, error)
	GetClassMessages(ctx context.Context, classID, limit int) ([]domain.Message, error)
	CanAccessClass(ctx context.Context, userID, classID int) (bool, error)
}

type messageUseCase struct {
	repo repository.MessageRepository
}

func NewMessageUseCase(repo repository.MessageRepository) MessageUseCase {
	return &messageUseCase{repo}
}

func (uc *messageUseCase) SendMessage(ctx context.Context, userID, classID int, content string) (domain.Message, error) {
	//! Vérif autorisation
	canAccess, err := uc.CanAccessClass(ctx, userID, classID)
	if err != nil || !canAccess {
		return domain.Message{}, fmt.Errorf("access denied to class")
	}

	return uc.repo.CreateMessage(ctx, userID, classID, content)
}

func (uc *messageUseCase) GetClassMessages(ctx context.Context, classID, limit int) ([]domain.Message, error) {
	return uc.repo.GetRecentMessages(ctx, classID, limit)
}

func (uc *messageUseCase) CanAccessClass(ctx context.Context, userID, classID int) (bool, error) {
	return uc.repo.UserInClass(ctx, userID, classID)
}
