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
	"golang.org/x/crypto/bcrypt"
)

func TestSignin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successfull Login", func(t *testing.T) {
		password, _ := bcrypt.GenerateFromPassword([]byte("user@123"), bcrypt.DefaultCost)
		Mock := initialzer.MockDbConfig(t)
		defer initialzer.Mockdb.Close()

		Mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email=\\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
			WithArgs("user@gmail.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(1, "user", "user@gmail.com", password))
		server := gin.Default()
		server.POST("/signin", controller.Signin)
		credentials := controller.SigninInput{
			Email:    "user@gmail.com",
			Password: "user@123",
		}
		values, _ := json.Marshal(credentials)
		req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(values))
		req.Header.Set("Content-Type", "application/json")

		write := httptest.NewRecorder()
		server.ServeHTTP(write, req)
		assert.Equal(t, http.StatusOK, write.Code)
		assert.Contains(t, write.Body.String(), "Successfully logined")
	})

	t.Run("Wrong Password", func(t *testing.T) {
		password, _ := bcrypt.GenerateFromPassword([]byte("user@123"), bcrypt.DefaultCost)
		mock := initialzer.MockDbConfig(t)
		defer initialzer.Mockdb.Close()

		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email=\\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
			WithArgs("user@exmple.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(1, "user", "user@gmail.com", password))
		server := gin.Default()
		server.POST("/signin", controller.Signin)
		credentials := controller.SigninInput{
			Email:    "user@exmple.com",
			Password: "userss@123",
		}
		values, _ := json.Marshal(credentials)
		req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(values))
		req.Header.Set("Content-Type", "application/json")

		write := httptest.NewRecorder()
		server.ServeHTTP(write, req)
		assert.Equal(t, http.StatusUnauthorized, write.Code)
		assert.Contains(t, write.Body.String(), "invalid username or password")
	})

	t.Run("invalid email", func(t *testing.T) {
		password, _ := bcrypt.GenerateFromPassword([]byte("user@123"), bcrypt.DefaultCost)
		mock := initialzer.MockDbConfig(t)
		defer initialzer.Mockdb.Close()

		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email=\\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
			WithArgs("user@exmple.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(1, "user", "user@gmail.com", password))

		server := gin.Default()
		server.POST("/signin", controller.Signin)
		credentials := controller.SigninInput{
			Email:    "hello@exmple.com",
			Password: "user@123",
		}
		values, _ := json.Marshal(credentials)
		req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(values))
		req.Header.Set("Content-Type", "application/json")

		write := httptest.NewRecorder()
		server.ServeHTTP(write, req)
		assert.Equal(t, http.StatusUnauthorized, write.Code)
		assert.Contains(t, write.Body.String(), "invalid email")
	})
}
