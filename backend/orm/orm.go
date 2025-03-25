package orm

import (
	"context"
	"fmt"
)

func ClearDB() error {
	_, err := dbpool.Exec(context.Background(), "TRUNCATE TABLE site, page, link RESTART IDENTITY CASCADE;")
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}
