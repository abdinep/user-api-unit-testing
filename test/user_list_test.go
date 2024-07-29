package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project/controller"
	"project/initialzer"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserList(t *testing.T) {
	mock := initialzer.MockDbConfig(t)
	defer initialzer.Mockdb.Close()
	gin.SetMode(gin.TestMode)

	t.Run("successfull user list", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL").WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(1, "user1", "user1@gmail.com", "user1@123").
				AddRow(2, "user2", "user2@gmail.com", "user2@123"),
		)

		server := gin.Default()
		server.GET("/userlist", controller.ListUser)

		req, _ := http.NewRequest("GET", "/userlist", nil)
		req.Header.Set("Content-Type", "application/json")
		write := httptest.NewRecorder()
		server.ServeHTTP(write, req)
		assert.Equal(t, http.StatusOK, write.Code)

		var response map[string]interface{}
		err := json.Unmarshal(write.Body.Bytes(), &response)
		assert.NoError(t, err)
		datas := response["data"].([]interface{})
		user1 := datas[0].(map[string]interface{})
		assert.Equal(t, "user1", user1["name"])
		assert.Equal(t, "user1@gmail.com", user1["email"])
		user2 := datas[1].(map[string]interface{})
		assert.Equal(t, "user2", user2["name"])
		assert.Equal(t, "user2@gmail.com", user2["email"])
		assert.Equal(t, int(response["total_user"].(float64)), 2)

	})

	t.Run("failed to fetch data", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL").
			WillReturnError(gorm.ErrInvalidTransaction)

		server := gin.Default()
		server.GET("/userlist", controller.ListUser)

		req, _ := http.NewRequest("GET", "/userlist", nil)
		req.Header.Set("Content-Type", "application/json")
		write := httptest.NewRecorder()
		server.ServeHTTP(write, req)
		assert.Equal(t, http.StatusInternalServerError, write.Code)
		var response map[string]string
		err := json.Unmarshal(write.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Failed to fetch users", response["error"])
	})
}
