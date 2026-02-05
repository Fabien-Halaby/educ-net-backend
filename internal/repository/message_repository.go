package repository

import (
	"context"
	"database/sql"
	"educnet/internal/domain"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, userID, classID int, content string) (domain.Message, error)
	GetRecentMessages(ctx context.Context, classID int, limit int) ([]domain.Message, error)
	UserInClass(ctx context.Context, userID, classID int) (bool, error)
	DeleteMessage(ctx context.Context, messageID int) error
}

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) CreateMessage(ctx context.Context, userID, classID int, content string) (domain.Message, error) {
	var messageID int64
	err := r.db.QueryRowContext(ctx, `
        INSERT INTO messages (content, user_id, class_id) 
        VALUES ($1, $2, $3) 
        RETURNING id
    `, content, userID, classID).Scan(&messageID)
	if err != nil {
		return domain.Message{}, err
	}

	var msg domain.Message
	err = r.db.QueryRowContext(ctx, `
        SELECT * FROM messages_view WHERE id = $1
    `, messageID).Scan(
		&msg.ID, &msg.Content, &msg.MessageType, &msg.FileURL, &msg.IsPinned,
		&msg.CreatedAt, &msg.UpdatedAt, &msg.User.ID, &msg.User.FirstName,
		&msg.User.LastName, &msg.User.FullName, &msg.User.Role, &msg.User.AvatarURL,
		&msg.ClassID, &msg.ClassName,
	)
	return msg, err
}

func (r *messageRepository) GetRecentMessages(ctx context.Context, classID int, limit int) ([]domain.Message, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT * FROM get_recent_messages($1, $2)
    `, classID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message

		rows.Scan(
			&msg.ID,
			&msg.Content,
			&msg.User.ID,
			&msg.User.FullName,
			&msg.User.Role,
			&msg.User.AvatarURL,
			&msg.CreatedAt,
		)

		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (r *messageRepository) UserInClass(ctx context.Context, userID, classID int) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM student_classes 
            WHERE student_id = $1 AND class_id = $2
        )
    `, userID, classID).Scan(&exists)

	return exists, err
}

func (r *messageRepository) DeleteMessage(ctx context.Context, messageID int) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM messages WHERE id = $1", messageID)
	return err
}
