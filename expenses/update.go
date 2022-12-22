package expenses

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/lib/pq"
)

func UpdateExpense(c echo.Context) error {
	id := c.Param("id")
	expense := new(Expense)
	if err := c.Bind(expense); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`
			UPDATE expenses
			SET title = $1, amount = $2, note = $3, tags = $4
			WHERE id = $5;
		`, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	intVar, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, Expense{
		ID:     intVar,
		Title:  expense.Title,
		Amount: expense.Amount,
		Note:   expense.Note,
		Tags:   expense.Tags,
	})
}
