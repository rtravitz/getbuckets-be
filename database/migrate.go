package database

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for a db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add buckets",
		Script: `
		CREATE TABLE buckets (
			id SERIAL NOT NULL PRIMARY KEY,
			lat DOUBLE PRECISION NOT NULL,
			lng DOUBLE PRECISION NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE OR REPLACE FUNCTION trigger_set_timestamp()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER set_timestamp_buckets
		BEFORE UPDATE ON buckets
		FOR EACH ROW
		EXECUTE PROCEDURE trigger_set_timestamp();
		`,
	},
}
