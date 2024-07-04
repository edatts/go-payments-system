package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	// "strings"
	"testing"

	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type mockStore struct {
	users []*types.User
}

func (m *mockStore) CreateUser(user *types.User) error {

	// Ensure no unique constraints are violated
	for _, existingUser := range m.users {
		if existingUser.Username == user.Username {
			return ErrUsernameExists
		}
		if existingUser.Email == user.Email {
			return ErrEmailExists
		}
	}

	m.users = append(m.users, user)

	return nil
}

func (m *mockStore) GetUser(username string) (*types.User, error) {

	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}

	return &types.User{}, pgx.ErrNoRows
}

func (m *mockStore) GetUserById(id int32) (*types.User, error) {

	for _, user := range m.users {
		if user.Id == id {
			return user, nil
		}
	}

	return &types.User{}, pgx.ErrNoRows
}

func (m *mockStore) GetUserByEmail(email string) (*types.User, error) {

	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}

	return &types.User{}, pgx.ErrNoRows
}

func TestAuthHandlers(t *testing.T) {

	mockStore := &mockStore{}

	for _, testCase := range registerUserRequestTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			registerUserReq := testCase.req

			reqJson, err := json.Marshal(registerUserReq)
			if err != nil {
				t.Errorf("failed marshalling json: %s", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqJson))
			if err != nil {
				t.Errorf("failed creating http request: %s", err)
			}

			handler := NewHandler(mockStore)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()

			router.HandleFunc("/register", handler.handleRegister)
			router.ServeHTTP(rr, req)

			if rr.Code != testCase.expectedHttpCode {
				t.Errorf("expected %s, got %s", http.StatusText(testCase.expectedHttpCode), http.StatusText(rr.Code))
			}
		})
	}

	for _, testCase := range loginTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			loginReq := testCase.req

			reqJson, err := json.Marshal(loginReq)
			if err != nil {
				t.Errorf("faile marshalling json: %s", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqJson))
			if err != nil {
				t.Errorf("failed creating http request: %s", err)
			}

			handler := NewHandler(mockStore)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()

			router.HandleFunc("/login", handler.handleLogin)
			router.ServeHTTP(rr, req)

			if rr.Code != testCase.expectedHttpCode {
				t.Errorf("expected %s, got %s", http.StatusText(testCase.expectedHttpCode), http.StatusText(rr.Code))
			}

			body := rr.Body.String()
			// log.Printf("response body: %s", body)
			// log.Printf("expected Err: %s", testCase.expectedError)
			if testCase.expectedError != nil {
				if !strings.Contains(body, testCase.expectedError.Error()) {
					t.Errorf("expected %s, got %s", testCase.expectedError, body)
				}
			}

		})

	}

}

func registerUserRequestTestCases() []registerUserRequestTestCase {
	return []registerUserRequestTestCase{
		{
			"Test valid input",
			types.RegisterUserRequest{
				FirstName: "Susan",
				LastName:  "Coleman",
				Username:  "SusieC",
				Email:     "susan.coleman@yahoo.com",
				Password:  "abcdefg",
			},
			http.StatusCreated,
		},
		{
			"Test required fields missing",
			types.RegisterUserRequest{
				FirstName: "",
				LastName:  "",
				Username:  "",
				Email:     "",
				Password:  "",
			},
			http.StatusBadRequest,
		},
		{
			"Test invalid email",
			types.RegisterUserRequest{
				FirstName: "Susan",
				LastName:  "Coleman",
				Username:  "SusieC",
				Email:     "invalidEmail",
				Password:  "abcdefg",
			},
			http.StatusBadRequest,
		},
		{
			"Test short password",
			types.RegisterUserRequest{
				FirstName: "Susan",
				LastName:  "Coleman",
				Username:  "SusieC",
				Email:     "susan.coleman@yahoo.com",
				Password:  "xyz",
			},
			http.StatusBadRequest,
		},
		{
			"Test long password",
			types.RegisterUserRequest{
				FirstName: "Susan",
				LastName:  "Coleman",
				Username:  "SusieC",
				Email:     "susan.coleman@yahoo.com",
				Password:  "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			http.StatusBadRequest,
		},
	}
}

type registerUserRequestTestCase struct {
	name             string
	req              types.RegisterUserRequest
	expectedHttpCode int
}

type loginTestCase struct {
	name             string
	req              types.LoginRequest
	expectedHttpCode int
	expectedError    error
}

func loginTestCases() []loginTestCase {
	return []loginTestCase{
		{
			name: "Test valid input",
			req: types.LoginRequest{
				Username: "SusieC",
				Password: "abcdefg",
			},
			expectedHttpCode: http.StatusOK,
			expectedError:    nil,
		},
		{
			name: "Test invalid input",
			req: types.LoginRequest{
				Username: "SusieC",
				Email:    "susan.coleman@yahoo.com",
				Password: "abcdefg",
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedError:    ErrFailedValidation,
		},
		{
			name: "Test invalid input",
			req: types.LoginRequest{
				Password: "abcdefg",
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedError:    ErrFailedValidation,
		},
		{
			name: "Test non-existent email",
			req: types.LoginRequest{
				Email:    "not.real@gmail.lol",
				Password: "password",
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedError:    ErrEmailNotExists,
		},
		{
			name: "Test non-existent username",
			req: types.LoginRequest{
				Username: "I do not exist",
				Password: "password",
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedError:    ErrUsernameNotExists,
		},
		{
			name: "Test wrong password",
			req: types.LoginRequest{
				Username: "SusieC",
				Password: "WrongPassword12",
			},
			expectedHttpCode: http.StatusBadRequest,
			expectedError:    ErrWrongPassword,
		},
	}
}
