package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LiangNing7/go-example/test/mysql/store"
	"github.com/LiangNing7/go-example/test/mysql/store/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_CreateUser_by_mock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := mocks.NewMockUserStore(ctrl)
	mockUserStore.EXPECT().Create(&store.User{
		Name: "user1",
	}).Return(nil)

	handler := &UserHandler{store: mockUserStore}
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	body := `{"name": "user1"}`
	reader := strings.NewReader(body)
	req := httptest.NewRequest("POST", "users", reader)
	router.ServeHTTP(w, req)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "", w.Body.String())
}

func TestUserHandler_GetUser_by_mock(t *testing.T) {
	// 创建一个 Controller 对象，用来管理整个测试中 Mock 对象的生命周期.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Finish() 会检查所有通过 EXPECT() 设置的期望调用（expectations）是否都被实际触发

	// 实例化一个 MockUserStore.
	mockUserStore := mocks.NewMockUserStore(ctrl)
	// EXPECT(): 告诉 Mock “接下来我要设置某个方法的期望调用”。
	// Create(&store.User{Name: "user1"})：表示期望在测试过程中，Mock 的 Create 方法会被以 &store.User{Name:"user1"} 作为参数调用一次。
	// Return(nil)：当上述调用发生时，Mock 就会返回 nil（即“创建成功，没有错误”）。
	mockUserStore.EXPECT().
		Get(2).
		Return(&store.User{
			ID:   2,
			Name: "user2",
		}, nil)

	handler := &UserHandler{store: mockUserStore}
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":2,"name":"user2"}`, w.Body.String())
}
