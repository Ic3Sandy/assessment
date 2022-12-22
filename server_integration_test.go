//go:build integration

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Ic3Sandy/assessment/expenses"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationCreateExpense(t *testing.T) {
	os.Setenv("DATABASE_URL", os.Getenv("DATABASE_URL_TEST"))
	tableName := "test_create_expense"
	db := ConnectAndCreateTable(tableName)
	defer db.Close()
	expenses.SetDB(db)
	expenses.SetTableName(tableName)

	dataJSON := `{"title":"test-title","amount":100,"note":"test-note","tags":["test-tags"]}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(dataJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, expenses.CreateExpense(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "test-title", response["title"])
		assert.Equal(t, 100.0, response["amount"])
		assert.Equal(t, "test-note", response["note"])
		assert.Equal(t, []interface{}{"test-tags"}, response["tags"])
	}
}

func TestIntegrationGetExpenses(t *testing.T) {
	os.Setenv("DATABASE_URL", os.Getenv("DATABASE_URL_TEST"))
	tableName := "test_get_expenses"
	db := ConnectAndCreateTable(tableName)
	defer db.Close()
	expenses.SetDB(db)
	expenses.SetTableName(tableName)

	db.Exec("INSERT INTO "+tableName+" (title, amount, note, tags) VALUES ($1, $2, $3, $4)",
		"test-title", 100, "test-note", pq.Array([]string{"test-tags"}))

	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	c := e.NewContext(req, rec)

	if assert.NoError(t, expenses.GetExpenses(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response []map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.GreaterOrEqual(t, len(response), 1)
		assert.Equal(t, "test-title", response[0]["title"])
		assert.Equal(t, 100.0, response[0]["amount"])
		assert.Equal(t, "test-note", response[0]["note"])
		assert.Equal(t, []interface{}{"test-tags"}, response[0]["tags"])
	}
}

func TestIntegrationGetExpensesById(t *testing.T) {
	os.Setenv("DATABASE_URL", os.Getenv("DATABASE_URL_TEST"))
	tableName := "test_get_expenses_by_id"
	db := ConnectAndCreateTable(tableName)
	defer db.Close()
	expenses.SetDB(db)
	expenses.SetTableName(tableName)

	db.Exec("INSERT INTO "+tableName+" (title, amount, note, tags) VALUES ($1, $2, $3, $4)",
		"test-title", 100, "test-note", pq.Array([]string{"test-tags"}))

	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/expenses/1", nil)
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, expenses.GetExpensesById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "test-title", response["title"])
		assert.Equal(t, 100.0, response["amount"])
		assert.Equal(t, "test-note", response["note"])
		assert.Equal(t, []interface{}{"test-tags"}, response["tags"])
	}
}

func TestIntegrationUpdateExpense(t *testing.T) {
	os.Setenv("DATABASE_URL", os.Getenv("DATABASE_URL_TEST"))
	tableName := "test_update_expense"
	db := ConnectAndCreateTable(tableName)
	defer db.Close()
	expenses.SetDB(db)
	expenses.SetTableName(tableName)

	db.Exec("INSERT INTO "+tableName+" (title, amount, note, tags) VALUES ($1, $2, $3, $4)",
		"test-title", 100, "test-note", pq.Array([]string{"test-tags"}))

	dataJSON := `{"title":"update-title","amount":123,"note":"update-note","tags":["update-tags"]}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses/1", strings.NewReader(dataJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, expenses.UpdateExpense(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "update-title", response["title"])
		assert.Equal(t, 123.0, response["amount"])
		assert.Equal(t, "update-note", response["note"])
		assert.Equal(t, []interface{}{"update-tags"}, response["tags"])
	}
}
