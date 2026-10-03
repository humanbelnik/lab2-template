package postgres

import "database/sql"

const schema = `
CREATE TABLE IF NOT EXISTS payment
(
    id          SERIAL PRIMARY KEY,
    payment_uid uuid        NOT NULL,
    status      VARCHAR(20) NOT NULL
        CHECK (status IN ('PAID', 'CANCELED')),
    price       INT         NOT NULL
);
`

func Migrate(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
