package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestUsers() func() {
	// 保存 users.
	defaultUser := users

	// 构造测试 users.
	users = []User{
		{ID: 1, Name: "test-user1"},
	}

	// 还原.
	return func() {
		users = defaultUser
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error")
}

func TestCreateUserHandler(t *testing.T) {
	cleanup := setupTestUsers()
	defer cleanup()

	w := httptest.NewRecorder()

	body := strings.NewReader(`{"name": "user2"}`)
	req := httptest.NewRequest("POST", "/users", body)

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "", w.Body.String())

	assert.Equal(t, 2, len(users))
	u2, _ := json.Marshal(users[1])
	assert.Equal(t, `{"id":2,"name":"user2"}`, string(u2))
}

func TestCreateUserHandler_BadRequest(t *testing.T) {
	cleanup := setupTestUsers()
	defer cleanup()

	w := httptest.NewRecorder()
	// 使用 errorReader, 使 io.ReadAll 失败.
	req := httptest.NewRequest("POST", "/users", &errorReader{})

	router := setupRouter()
	router.ServeHTTP(w, req)

	// 期待 400 Bad Request.
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	// 响应体里包含 read error 提示
	assert.Contains(t, w.Body.String(), `"msg":"read error"`)

	// 确保没有新增用户
	assert.Equal(t, 1, len(users))
}

func TestCreateUserHandler_UnmarshalError(t *testing.T) {
	cleanup := setupTestUsers()
	defer cleanup()

	w := httptest.NewRecorder()
	// 构造一个无法解析的 JSON 字符串
	body := strings.NewReader(`{"name": 123}`) // name 应该是 string
	req := httptest.NewRequest("POST", "/users", body)

	router := setupRouter()
	router.ServeHTTP(w, req)

	// 期待 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	// 响应体里包含 json.Unmarshal 的错误信息
	assert.Contains(t, w.Body.String(), `"msg":`)

	// 确保没有新增用户
	assert.Equal(t, 1, len(users))
}

func TestGetUserHandler(t *testing.T) {
	cleanup := setupTestUsers()
	defer cleanup()

	type want struct {
		code int
		body string
	}
	tests := []struct {
		name string
		args int
		want want
	}{
		{
			name: "get test-user1",
			args: 1,
			want: want{
				code: 200,
				body: `{"id":1,"name":"test-user1"}`,
			},
		},
		{
			name: "get user not found",
			args: 2,
			want: want{
				code: 404,
				body: `{"msg":"not found"}`,
			},
		},
	}

	router := setupRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d", tt.args), nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}
