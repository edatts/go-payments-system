package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type mockStore struct{}

func (m *mockStore) CreateUser(user *types.User) error {

	return nil
}

func (m *mockStore) GetUserById(id uint64) (*types.User, error) {

	return &types.User{}, nil
}

func (m *mockStore) GetUserByEmail(email string) (*types.User, error) {

	return &types.User{}, pgx.ErrNoRows
}

func TestAuthHandlers(t *testing.T) {

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

			mockStore := &mockStore{}
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

}

func registerUserRequestTestCases() []registerUserRequestTestCase {
	return []registerUserRequestTestCase{
		{
			"Test valid input",
			types.RegisterUserRequest{
				FirstName: "Susan",
				LastName:  "Coleman",
				Username:  "SuperSue",
				Email:     "abc@yahoo.com",
				Password:  "xyzsted",
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
				Username:  "SuperSue",
				Email:     "invalidEmail",
				Password:  "xyzsted",
			},
			http.StatusBadRequest,
		},
		{
			"Test short password",
			types.RegisterUserRequest{
				FirstName: "Susan",
				LastName:  "Coleman",
				Username:  "SuperSue",
				Email:     "abc@yahoo.com",
				Password:  "xyz",
			},
			http.StatusBadRequest,
		},
		{
			"Test long password",
			types.RegisterUserRequest{
				FirstName: "Susan",
				LastName:  "Coleman",
				Username:  "SuperSue",
				Email:     "abc@yahoo.com",
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
