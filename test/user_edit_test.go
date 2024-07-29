package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project/controller"
	"project/initialzer"
	"project/models"
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
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "email"=\$1,"name"=\$2,"updated_at"=\$3 WHERE id = \$4 AND "users"."deleted_at" IS NULL`).
			WithArgs("userEdit1@gmail.com", "userEdit1", sqlmock.AnyArg(), 11).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		router := gin.Default()
		router.PATCH("/user/edit/:id", controller.EditUser)

		user := models.User{
			Name:  "userEditsd1",
			Email: "userEdit1@gmail.com",
		}
		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPatch, "/user/edit/11", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Successfully updated user")
	})
}
