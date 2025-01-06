package models

import "time"

type TgMessage struct {
	ID               int       `db:"id"`
	Date             time.Time `db:"date"`
	Text             string    `db:"text"`
	ChatID           int       `db:"chat_id"`
	FirstName        string    `db:"first_name"`
	LastName         string    `db:"last_name"`
	UserName         string    `db:"user_name"`
	FromID           int       `db:"from_id"`
	ReplyToMessageID int       `db:"reply_message_id"`
	BotMessageID     *int      `db:"bot_message_id"`
}
