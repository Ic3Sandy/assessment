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
	mock.ExpectQuery("INSERT INTO expenses").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	SetDB(db)
	SetTableName("expenses")

	if assert.NoError(t, CreateExpense(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, `{"id":1,"title":"test","amount":100,"note":"test","tags":["test"]}`+"\n", rec.Body.String())
	}
}
