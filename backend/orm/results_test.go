package orm_test

// import (
// 	"testing"

// 	"github.com/ZakkBob/AskDave/backend/orm"
// 	"github.com/ZakkBob/AskDave/gocommon/hash"
// 	"github.com/ZakkBob/AskDave/gocommon/robots"
// 	"github.com/ZakkBob/AskDave/gocommon/tasks"
// 	"github.com/ZakkBob/AskDave/gocommon/url"
// 	"github.com/pashagolub/pgxmock/v4"
// )

// // func TestValidatorByUrl(t *testing.T) {
// // 	mock, err := pgxmock.NewPool()
// // 	if err != nil {
// // 		t.Fatal(err)
// // 	}
// // 	defer mock.Close()

// // 	allowed := []string{"allowed1", "allowed2"}
// // 	disallowed := []string{"disallowed1", "disallowed2", "disallowed3"}

// // 	rows := pgxmock.NewRows([]string{"allowed_patterns", "disallowed_patterns"}).
// // 		AddRow(allowed, disallowed)

// // 	mock.ExpectQuery(`SELECT`).
// // 		WithArgs("https://mateishome.page").
// // 		WillReturnRows(rows)

// // 	v, err := robots.ValidatorByUrl(mock, "https://mateishome.page")

// // 	assert.Equal(t, v.AllowedStrings(), allowed, "allowed patterns should match")
// // 	assert.Equal(t, v.DisallowedStrings(), disallowed, "disallowed patterns should match")

// // 	if err != nil {
// // 		t.Errorf("error was not expected while getting validator: %s", err)
// // 	}

// // 	if err := mock.ExpectationsWereMet(); err != nil {
// // 		t.Errorf("there were unfulfilled expectations: %s", err)
// // 	}
// // }

// func TestSaveResultsRobots(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	orm.MockConnect(mock)
// 	defer orm.Close()

// 	u, _ := url.ParseAbs("test.com")
// 	u2, _ := url.ParseAbs("test2.com")
// 	u3, _ := url.ParseAbs("test3.com")

// 	allowedStrings := []string{
// 		"123",
// 		"123",
// 		"123",
// 	}

// 	disallowedStrings := []string{
// 		"123",
// 		"123",
// 		"123",
// 	}

// 	validator, _ := robots.FromStrings(allowedStrings, disallowedStrings)

// 	results := tasks.Results{
// 		Robots: map[string]*tasks.RobotsResult{
// 			u.String(): {
// 				Url:           &u,
// 				Success:       true,
// 				FailureReason: tasks.NoFailure,
// 				Hash:          hash.Hashs(""),
// 				Changed:       true,
// 				Validator:     validator,
// 			}, u2.String(): {
// 				Url:           &u2,
// 				Success:       false,
// 				FailureReason: tasks.NoFailure,
// 				Hash:          hash.Hashs(""),
// 				Changed:       true,
// 				Validator:     validator,
// 			}, u3.String(): {
// 				Url:           &u3,
// 				Success:       true,
// 				FailureReason: tasks.NoFailure,
// 				Hash:          hash.Hashs(""),
// 				Changed:       false,
// 				Validator:     validator,
// 			}},
// 		Pages: map[string]*tasks.PageResult{},
// 	}

// 	mock.ExpectExec("UPDATE").
// 		WithArgs(allowedStrings, disallowedStrings, u.String()).
// 		WillReturnResult(pgxmock.NewResult("", 1))

// 	err = orm.SaveResults(&results)
// 	if err != nil {
// 		t.Errorf("error was not expected while saving results: %s", err)
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }
