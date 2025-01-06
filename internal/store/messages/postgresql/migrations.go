package postgresql

const queryCreateTableMessages = `CREATE TABLE IF NOT EXISTS messages (
	id INT NOT NULL,
    date TIMESTAMPTZ NOT NULL, 
    text TEXT,
    chat_id INT NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    user_name VARCHAR(255),
    from_id INT,
    reply_message_id INT,
    bot_message_id INT,
    PRIMARY KEY (id)
);`

func (s Store) RunMigrations() error {
	if _, err := s.db.Exec(queryCreateTableMessages); err != nil {
		return err
	}

	return nil
}
