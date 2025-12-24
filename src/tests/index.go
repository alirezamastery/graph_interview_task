package todoctrltest

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	todoctrl "github.com/alirezamastery/graph_task/controllers/todo"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func NewMockGormDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(
		postgres.Config{
			Conn: sqlDB,
			// PreferSimpleProtocol: true,
		}),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatalf("gorm error: %v", err)
	}

	return gdb, mock, sqlDB
}

func SetupRouter(ctl *todoctrl.Controller) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/task/todos/", ctl.CreateTodo())
	return r
}
