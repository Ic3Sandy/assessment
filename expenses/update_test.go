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

func TestUpdateExpense(t *testing.T) {
	dataJSON := `{"title":"test","amount":100,"note":"test","tags":["test"]}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses/1", strings.NewReader(dataJSON))
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
	mock.ExpectExec("UPDATE expenses").WillReturnResult(sqlmock.NewResult(1, 1))
	SetDB(db)

	if assert.NoError(t, UpdateExpense(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"id":1,"title":"test","amount":100,"note":"test","tags":["test"]}`+"\n", rec.Body.String())
	}
}
