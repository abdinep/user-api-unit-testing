package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project/controller"
	"project/initialzer"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserEdit(t *testing.T) {
	mock := initialzer.MockDbConfig(t)
	defer initialzer.Mockdb.Close()
	gin.SetMode(gin.TestMode)

	t.Run("successfully edited user", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT 1`).
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "user1", "user1@gmail.com"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "email"=\$1,"name"=\$2,"updated_at"=\$3 WHERE id = \$4 AND "users"."deleted_at" IS NULL`).
			WithArgs("new@gmail.com", "userEdit1", sqlmock.AnyArg(), 11).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		server := gin.Default()
		server.PATCH("/user/edit/:ID", controller.EditUser)

		input := controller.EditUserInput{
			Name:  "userEdit1",
			Email: "userEdit1@gmail.com",
		}
		jsonValue, _ := json.Marshal(input)
		req, _ := http.NewRequest("PATCH", "/user/edit/:ID", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		write := httptest.NewRecorder()
		server.ServeHTTP(write, req)

		assert.Equal(t, http.StatusOK, write.Code)
		assert.Contains(t, write.Body.String(), "user updated")
	})
}
