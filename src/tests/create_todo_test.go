package todoctrltest

import (
	"bytes"
	"encoding/json"
	"github.com/alirezamastery/graph_task/middleware"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/prometheus/client_golang/prometheus/testutil"

	todoctrl "github.com/alirezamastery/graph_task/controllers/todo"
)

func TestCreateTodo_201_InsertsAndIncrementsGauge(t *testing.T) {
	db, mock, sqlDB := NewMockGormDB(t)
	t.Cleanup(func() { _ = sqlDB.Close() })

	ctl := todoctrl.NewTodoController(db)
	router := SetupRouter(ctl)

	before := testutil.ToFloat64(middleware.TasksCount)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta("INSERT")).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()

	body := []byte(`{"title":" test 1 ","description":" desc ","is_done":false}`)
	req := httptest.NewRequest(http.MethodPost, "/api/task/todos/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", recorder.Code, recorder.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v, body=%s", err, recorder.Body.String())
	}

	if resp["title"] != "test 1" {
		t.Fatalf("expected trimmed title, got %#v", resp["title"])
	}
	if resp["description"] != "desc" {
		t.Fatalf("expected trimmed description, got %#v", resp["description"])
	}
	if resp["is_done"] != false {
		t.Fatalf("expected is_done=false, got %#v", resp["is_done"])
	}

	after := testutil.ToFloat64(middleware.TasksCount)
	if after != before+1 {
		t.Fatalf("expected TasksCount +1, before=%v after=%v", before, after)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("db expectations not met: %v", err)
	}
}

func TestCreateTodo_400_EmptyTitle_NoDBCall(t *testing.T) {
	db, mock, sqlDB := NewMockGormDB(t)
	t.Cleanup(func() { _ = sqlDB.Close() })

	ctl := todoctrl.NewTodoController(db)
	router := SetupRouter(ctl)

	before := testutil.ToFloat64(middleware.TasksCount)

	body := []byte(`{"title":"   ","description":"test 1 desc","is_done":false}`)
	req := httptest.NewRequest(http.MethodPost, "/api/task/todos/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", recorder.Code, recorder.Body.String())
	}

	after := testutil.ToFloat64(middleware.TasksCount)
	if after != before {
		t.Fatalf("expected TasksCount unchanged, before=%v after=%v", before, after)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("db expectations not met: %v", err)
	}
}
