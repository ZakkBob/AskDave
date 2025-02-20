package main

import (
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

type AnyTime struct{}

func TestShouldGetValidatorByUrl(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"allowed_patterns, disallowed_patterns"}).
		AddRow([]string{"", ""})

	mock.ExpectQuery(`SELECT`).
		WithArgs("https://mateishome.page").
		WillReturnRows(rows)

	_, err = ValidatorByUrl(mock, "https://mateishome.page")
	if err != nil {
		t.Errorf("error was not expected while getting validator: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
