package featureTransaction

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	user "backend/app/feature/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserService interface {
	GetUser(userId int) (user.User, error)
}

type TransactionService interface {
	GetTransactions(userId int) ([]Transaction, error)
	CreateTransaction(c *gin.Context) error
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(userId int) (user.User, error) {
	args := m.Called(userId)
	return args.Get(0).(user.User), args.Error(1)
}

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) GetTransactions(userId int) ([]Transaction, error) {
	args := m.Called(userId)
	return args.Get(0).([]Transaction), args.Error(1)
}

func (m *MockTransactionService) CreateTransaction(c *gin.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockTransactionService) InsertTransaction(amount int, userFromId *int, userToId int) {
	m.Called(amount, userFromId, userToId)
}

type Services struct {
	userService        UserService
	transactionService TransactionService
}

func NewService(userService UserService, transactionService TransactionService) *Services {
	return &Services{
		userService:        userService,
		transactionService: transactionService,
	}
}

func (ss *Services) GetUserBalance(userId int) (user.User, error) {
	return ss.userService.GetUser(userId)
}

func TestGetTransactions(t *testing.T) {
	tests := []struct {
		name        string
		userId      int
		mockResult  []Transaction
		mockError   error
		expected    []Transaction
		expectedErr error
	}{
		{
			name:        "[transactions found]",
			userId:      1,
			mockResult:  []Transaction{},
			mockError:   nil,
			expected:    []Transaction{},
			expectedErr: nil,
		},
		{
			name:        "[transactions not found]",
			userId:      2,
			mockResult:  nil,
			mockError:   errors.New("transactions not found"),
			expected:    nil,
			expectedErr: errors.New("transactions not found"),
		},
	}

	mockUserService := new(MockUserService)
	mockTransactionService := new(MockTransactionService)
	services := NewService(mockUserService, mockTransactionService)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockTransactionService.On("GetTransactions", test.userId).Return(test.mockResult, test.mockError)

			result, err := services.transactionService.GetTransactions(test.userId)

			assert.Equal(t, test.expected, result, "Expected transactions %v, got %v", test.expected, result)
			assert.Equal(t, test.expectedErr, err, "Expected error %v, got %v", test.expectedErr, err)

			mockTransactionService.AssertCalled(t, "GetTransactions", test.userId)
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		amount             int
		mockUserTo         user.User
		mockUserToErr      error
		mockUserFrom       user.User
		mockUserFromErr    error
		mockTransactionErr error
		expectedErr        string
	}{
		{
			name:               "[transaction complete]",
			requestBody:        `{"amount": 100, "userToId": 2, "userFromId": 1}`,
			amount:             100,
			mockUserTo:         user.User{Id: 2, Balance: 200},
			mockUserToErr:      nil,
			mockUserFrom:       user.User{Id: 1, Balance: 300},
			mockUserFromErr:    nil,
			mockTransactionErr: nil,
			expectedErr:        "",
		},
		{
			name:            "[balance error]",
			requestBody:     `{"amount": 100, "userFromId": 1, "userToId": 2}`,
			amount:          100,
			mockUserTo:      user.User{Id: 2, Balance: 150},
			mockUserToErr:   nil,
			mockUserFrom:    user.User{Id: 1, Balance: 50},
			mockUserFromErr: nil,
			expectedErr:     "баланс меньше суммы перевода",
		},
		{
			name:        "[cant process request error]",
			requestBody: `invalid json`,
			expectedErr: "error: cant process transaction request",
		},
	}

	mockUserService := new(MockUserService)
	mockTransactionService := new(MockTransactionService)
	services := NewService(mockUserService, mockTransactionService)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			if !strings.Contains(test.requestBody, ":") {
				err = fmt.Errorf("error: cant process transaction request")
			}
			if strings.Contains(test.requestBody, "userToId") {
				t.Log("[userToId]")
				mockUserService.On("GetUser", test.mockUserTo.Id).Return(test.mockUserTo, test.mockUserToErr)
				services.GetUserBalance(test.mockUserTo.Id)
			}
			if strings.Contains(test.requestBody, "userFromId") {
				t.Log("[userFromId]")
				mockUserService.On("GetUser", test.mockUserFrom.Id).Return(test.mockUserFrom, test.mockUserFromErr)
				services.GetUserBalance(test.mockUserFrom.Id)
			}
			if test.amount != 0 && test.amount > test.mockUserFrom.Balance {
				err = fmt.Errorf("баланс меньше суммы перевода")
			}

			mockTransactionService.On("CreateTransaction", mock.Anything).Return(test.mockTransactionErr)

			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			services.transactionService.CreateTransaction(c)

			if test.expectedErr != "" {
				assert.Error(t, err, "Expected error, got nil")
				assert.Contains(t, err.Error(), test.expectedErr, "Expected error %v, got %v", test.expectedErr, err)
			} else {
				assert.NoError(t, err, "Expected no error, got %v", err)
			}

			if err != nil {
				return
			}

			// проверяем что все ожидания выполнены
			mockUserService.AssertCalled(t, "GetUser", test.mockUserTo.Id)
			mockTransactionService.AssertExpectations(t)
		})
	}
}
