package services

import (
	"context"
	"errors"
	"testing"

	"github.com/chamanbetra/user-management-app/database"
	"github.com/chamanbetra/user-management-app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) error {
	args := Called(value)
	return args.Error(0)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *MockDB {
	m.Called(query, args)
	return m
}

func (m *MockDB) First(out interface{}) error {
	args := m.Called(out)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockDB := new(MockDB)
	database.DB = mockDB

	user := &models.User{Email: "test@example.com"}

	// Expect the Where call to return no user (simulating a new user)
	mockDB.On("Where", "email = ?", "test@example.com").Return(mockDB)
	mockDB.On("First", mock.Anything).Return(errors.New("not found"))

	// Expect the Create call
	mockDB.On("Create", user).Return(nil)

	err := CreateUser(context.Background(), user)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestGetUserByEmail(t *testing.T) {
	mockDB := new(MockDB)
	database.DB = mockDB

	email := "test@example.com"
	expectedUser := &models.User{Email: email}

	mockDB.On("Where", "email = ?", email).Return(mockDB)
	mockDB.On("First", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.User)
		*arg = *expectedUser
	}).Return(nil)

	user, err := GetUserByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockDB.AssertExpectations(t)
}
