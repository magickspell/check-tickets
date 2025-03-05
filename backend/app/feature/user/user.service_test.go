package featureUser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type UserService interface {
	GetUser(userId int) (User, error)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(userId int) (User, error) {
	args := m.Called(userId)
	return args.Get(0).(User), args.Error(1)
}

type Services struct {
	userService UserService
}

func (ss *Services) GetUserBalance(userId int) (User, error) {
	return ss.userService.GetUser(userId)
}

func NewService(userService UserService) *Services {
	return &Services{userService: userService}
}

func TestGetUserBalance(t *testing.T) {
	tests := []struct {
		name        string
		userId      int
		mockUser    User
		mockError   error
		expected    User
		expectedErr error
	}{
		{
			name:   "[user found]",
			userId: 1,
			mockUser: User{
				Id:      1,
				Balance: 100,
			},
			mockError: nil,
			expected: User{
				Id:      1,
				Balance: 100,
			},
			expectedErr: nil,
		},
		{
			name:        "[user not found]",
			userId:      2,
			mockUser:    User{},
			mockError:   errors.New("user not found"),
			expected:    User{},
			expectedErr: errors.New("user not found"),
		},
	}

	mockUserService := new(MockUserService)
	service := NewService(mockUserService)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserService.On("GetUser", test.userId).Return(test.mockUser, test.mockError)

			user, err := service.GetUserBalance(test.userId)

			assert.Equal(t, test.expected, user, "Expected user %v, got %v", test.expected, user)

			if test.expectedErr != nil {
				require.Error(t, err, "Expected error, got nil")
				assert.Equal(t, test.expectedErr.Error(), err.Error(), "Expected error %v, got %v", test.expectedErr, err)
			} else {
				require.NoError(t, err, "Expected no error, got %v", err)
			}

			mockUserService.AssertCalled(t, "GetUser", test.userId)
		})
	}
}
