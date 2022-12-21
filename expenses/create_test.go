package expenses

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {
	dataJSON := `{"title":"test","amount":100,"note":"test","tags":["test"]}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(dataJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec("INSERT INTO expenses").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	SetDB(db)

	if assert.NoError(t, CreateExpense(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, dataJSON+"\n", rec.Body.String())
	}
}
