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
	"github.com/stretchr/testify/require"
)

func TestSignup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid signup", func(t *testing.T) {
		mock := initialzer.MockDbConfig(t)
		defer initialzer.Mockdb.Close()
		server := gin.Default()
		server.POST("/signup", controller.Signup)
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users" \("created_at","updated_at","deleted_at","name","email","password"\) VALUES  \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "test user", "testuser@example.com", sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		payload := `{"name":"test user","email":"testuser@example.com","password":"user@123"}`
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		write := httptest.NewRecorder()
		server.ServeHTTP(write, req)

		if write.Code != http.StatusOK {
			t.Logf("Response Code:%d", write.Code)
			t.Logf("Response Body:%s", write.Body.String())
		}

		var response map[string]string
		err := json.Unmarshal(write.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "User created successfully", response["message"])

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations were not met: %v", err)
		}
	})

	t.Run("invalid signup", func(t *testing.T) {
		mock := initialzer.MockDbConfig(t)
		defer initialzer.Mockdb.Close()
		server := gin.Default()
		server.POST("/signup", controller.Signup)

		payload := `{"name":"","email":"","password":""}`
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		write := httptest.NewRecorder()
		server.ServeHTTP(write, req)

		if write.Code != http.StatusBadRequest {
			t.Logf("Response Code:%d", write.Code)
			t.Logf("Response Body:%s", write.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(write.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "All fields are required", response["error"])
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations were not met: %v", err)
		}
	})

}
