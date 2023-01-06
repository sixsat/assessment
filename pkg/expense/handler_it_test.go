//go:build integration

package expense

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const serverPort = 2565

func TestITCreateExpense(t *testing.T) {
	// Arrange
	eh := setup(t)
	defer teardown(t, eh)

	_, err := db.Exec("TRUNCATE TABLE expenses RESTART IDENTITY")
	if err != nil {
		t.Log(err)
	}

	body := strings.NewReader(`
		{"title": "strawberry smoothie", "amount": 79, "note": "night market promotion discount 10 bath", "tags": ["food","beverage"]}
	`)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), body)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	expected := `{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}
`

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, expected, string(byteBody))
	}
}

func TestITGetExpense(t *testing.T) {
	// Arrange
	eh := setup(t)
	defer teardown(t, eh)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/1", serverPort), strings.NewReader(""))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	expected := `{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}
`

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, string(byteBody))
	}
}

func setup(t *testing.T) *echo.Echo {
	InitDB("host=localhost port=5432 user=root password=secret dbname=expense sslmode=disable")

	eh := echo.New()
	go func(e *echo.Echo) {
		e.POST("/expenses", CreateExpense)
		e.GET("/expenses/:id", GetExpense)
		e.GET("/expenses", GetAllExpenses)
		e.PUT("/expenses/:id", UpdateExpense)

		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)

	for {
		conn, _ := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if conn != nil {
			conn.Close()
			break
		}
	}

	return eh
}

func teardown(t *testing.T, eh *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := eh.Shutdown(ctx)
	assert.NoError(t, err)
}
