//go:build unit

package expense

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	SetDB(db)

	mock.ExpectQuery("INSERT INTO expenses").
		WithArgs("strawberry smoothie", 79.0, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))

	body := strings.NewReader(`
		{"title": "strawberry smoothie", "amount": 79, "note": "night market promotion discount 10 bath", "tags": ["food","beverage"]}
	`)
	req := httptest.NewRequest(http.MethodPost, "/expenses", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	expected := `{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}
`

	// Act
	err = CreateExpense(c)
	actual := rec.Body.String()

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, expected, actual)
}
