package postgresql

import (
	"context"

	"scraper_nike/internal/models"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetMessageByID(ctx context.Context, msgID int) (models.TgMessage, error) {
	var msg models.TgMessage

	query := `
		SELECT id, date, text, chat_id, first_name, last_name, user_name, from_id, reply_message_id
		FROM messages
		WHERE id = ?
	`

	err := s.db.GetContext(ctx, &msg, query, msgID)

	if err != nil {
		return models.TgMessage{}, err
	}

	return msg, nil
}

func (s *Store) SaveMessage(ctx context.Context, msg models.TgMessage) error {
	query := `
		INSERT INTO messages (id, date, text, chat_id, first_name, last_name, user_name, from_id, reply_message_id)
		VALUES (:id, :date, :text, :chat_id, :first_name, :last_name, :user_name, :from_id, :reply_message_id)
		ON CONFLICT (id) DO UPDATE SET
			date = EXCLUDED.date,
			text = EXCLUDED.text,
			chat_id = EXCLUDED.chat_id,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			user_name = EXCLUDED.user_name,
			from_id = EXCLUDED.from_id,
			reply_message_id = EXCLUDED.reply_message_id
	`

	_, err := s.db.NamedExecContext(ctx, query, msg)
	return err
}
