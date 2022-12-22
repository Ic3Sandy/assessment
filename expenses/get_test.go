package expenses

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetExpenses(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// mock db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "test", 100, "test", pq.Array([]string{"test"}))
	mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses").WillReturnRows(rows)
	SetDB(db)

	if assert.NoError(t, GetExpenses(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `[{"id":1,"title":"test","amount":100,"note":"test","tags":["test"]}]`+"\n", rec.Body.String())
	}
}

func TestGetExpensesById(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "test", 100, "test", pq.Array([]string{"test"}))
	mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses").WillReturnRows(rows)
	SetDB(db)

	if assert.NoError(t, GetExpensesById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"id":1,"title":"test","amount":100,"note":"test","tags":["test"]}`+"\n", rec.Body.String())
	}
}
