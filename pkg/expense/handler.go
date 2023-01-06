package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func CreateExpense(c echo.Context) error {
	var exp Expense
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow(`
		INSERT INTO expenses (title, amount, note, tags)
		values ($1, $2, $3, $4)
		RETURNING id
		`, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags),
	)
	err = row.Scan(&exp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, exp)
}

func GetExpense(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT * FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	row := stmt.QueryRow(id)
	var exp Expense
	err = row.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, exp)
}

func GetAllExpenses(c echo.Context) error {
	stmt, err := db.Prepare("SELECT * FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	var exps []Expense
	for rows.Next() {
		var exp Expense
		err = rows.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
		exps = append(exps, exp)
	}

	return c.JSON(http.StatusOK, exps)
}

func UpdateExpense(c echo.Context) error {
	var exp Expense
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := db.Prepare(`
		UPDATE expenses
		SET title = $2, amount = $3, note = $4, tags = $5
		WHERE id = $1
		RETURNING id
		`,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	id := c.Param("id")
	row := stmt.QueryRow(id, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))
	err = row.Scan(&exp.ID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, exp)
}
