package server

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	json "github.com/json-iterator/go"
	"github.com/waikco/cats-v1/model"
)

func TestApp_CreateCat(t *testing.T) {
	var a App
	a.BootstrapServer()

	tests := []struct {
		description string
		// given
		requestBody []byte

		mockResponse struct {
			one string
			two error
		}
		expectedMockCalls int
		// then
		expectedStatus   int
		expectedResponse Response
	}{
		{
			description:    "successful pet creation",
			requestBody:    []byte(`{"name":"cat-1","color":"color-1","age":1}`),
			expectedStatus: http.StatusCreated,
			expectedResponse: Response{
				Result: "fe271e7e-83ca-477b-92fc-d0c3fa602d7d",
			},
			mockResponse: struct {
				one string
				two error
			}{
				one: "fe271e7e-83ca-477b-92fc-d0c3fa602d7d",
				two: nil},
			expectedMockCalls: 1,
		},
		{
			description:    "unsuccessful pet creation from bad json request",
			requestBody:    []byte(`name:"cat-1","color":"color-1","age":`),
			expectedStatus: http.StatusBadRequest,
			expectedResponse: Response{
				Error: Error{
					Status:  http.StatusBadRequest,
					Message: "invalid json in request body",
				},
			},
			mockResponse: struct {
				one string
				two error
			}{one: "", two: nil},
			expectedMockCalls: 0,
		},
		{
			description:    "unsuccessful pet creation from internal server error",
			requestBody:    []byte(`{"name":"cat-1","color":"color-1","age":1}`),
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: Response{
				Error: Error{
					Status:  http.StatusInternalServerError,
					Message: "error storing cat",
				},
			},
			mockResponse: struct {
				one string
				two error
			}{
				one: "",
				two: errors.New("database error"),
			},
			expectedMockCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := model.NewMockStorage(ctrl)
			s.EXPECT().
				Insert(gomock.Any()).
				Return(tt.mockResponse.one, tt.mockResponse.two).
				Times(tt.expectedMockCalls)
			a.Storage = s

			req, _ := http.NewRequest(http.MethodPost, "/cats/v1/", bytes.NewBuffer(tt.requestBody))
			response := httptest.NewRecorder()
			a.Router.ServeHTTP(response, req)

			var r Response
			err := json.Unmarshal(response.Body.Bytes(), &r)

			if err != nil {
				t.Fatalf("unxpected error: %v", err)
			}
			if response.Code != tt.expectedStatus {
				t.Errorf("unxpected status code: got %d, expected %d", response.Code, tt.expectedStatus)
			}
			if !reflect.DeepEqual(r, tt.expectedResponse) {
				t.Errorf("unxpected response body: got %v, expected %v", r, tt.expectedResponse)
			}
		})
	}
}

func TestApp_GetCat(t *testing.T) {
	var a App
	a.BootstrapServer()

	tests := []struct {
		description string
		// given
		request string

		mockResponse struct {
			one []byte
			two error
		}
		// then
		expectedStatus    int
		expectedResponse  interface{}
		expectedMockCalls int
	}{
		{
			description:    "successful",
			request:        "/cats/v1/cats/fe271e7e-83ca-477b-92fc-d0c3fa602d7d",
			expectedStatus: http.StatusOK,
			expectedResponse: model.Cat{
				Name:  "cat-1",
				Color: "color-1",
				Age:   1},
			mockResponse: struct {
				one []byte
				two error
			}{
				one: []byte(`{"name": "cat-1", "color": "color-1", "age": 1}`),
				two: nil},
			expectedMockCalls: 1,
		},
		{
			description:    "pet not found",
			request:        "/cats/v1/cats/fakepet",
			expectedStatus: http.StatusNotFound,
			expectedResponse: Response{
				Result: "cat not found",
			},
			mockResponse: struct {
				one []byte
				two error
			}{one: []byte(`[]`), two: sql.ErrNoRows},
			expectedMockCalls: 1,
		},
		{
			description:    "server error",
			request:        "/cats/v1/cats/fe271e7e-83ca-477b-92fc-d0c3fa602d7d",
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: Response{
				Result: "error getting cat",
			},
			mockResponse: struct {
				one []byte
				two error
			}{one: []byte(`[]`), two: errors.New("internal error")},
			expectedMockCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := model.NewMockStorage(ctrl)
			s.EXPECT().
				Select(gomock.Any()).
				Return(tt.mockResponse.one, tt.mockResponse.two).
				Times(tt.expectedMockCalls)
			a.Storage = s

			req := httptest.NewRequest(http.MethodGet, tt.request, nil)
			response := httptest.NewRecorder()
			a.Router.ServeHTTP(response, req)

			if response.Code != tt.expectedStatus {
				t.Errorf("unxpected status code: got %d, expected %d", response.Code, tt.expectedStatus)
			}

			switch tt.expectedResponse.(type) {
			case model.Cat:
				var got model.Cat
				err := json.Unmarshal(response.Body.Bytes(), &got)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, tt.expectedResponse.(model.Cat)) {
					t.Errorf("unxpected response got: got %v, expected %v",
						got, tt.expectedResponse)
				}
			case Response:
				var got Response
				err := json.Unmarshal(response.Body.Bytes(), &got)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, tt.expectedResponse.(Response)) {
					t.Errorf("unxpected response got: got %v, expected %v",
						got, tt.expectedResponse)
				}
			}
		})
	}
}

func TestApp_DeleteCat(t *testing.T) {
	var a App
	a.BootstrapServer()

	tests := []struct {
		description string
		// given
		request string

		mockResponse struct {
			one error
		}
		// then
		expectedStatus    int
		expectedResponse  Response
		expectedMockCalls int
	}{
		{
			description: "success",
			request:     "/cats/v1/cats/fe271e7e-83ca-477b-92fc-d0c3fa602d7d",
			mockResponse: struct {
				one error
			}{one: nil},
			expectedStatus:    http.StatusOK,
			expectedResponse:  Response{Result: "success"},
			expectedMockCalls: 1,
		},
		{
			description: "missing pet id",
			request:     "/cats/v1/cats/3",
			mockResponse: struct {
				one error
			}{one: nil},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: Response{Error: Error{
				Status:  http.StatusBadRequest,
				Message: "invalid cat id: 3",
			}},
		},
		{
			description: "pet not found",
			request:     "/cats/v1/cats/fe271e7e-83ca-477b-92fc-d0c3fa602d7d",
			mockResponse: struct {
				one error
			}{one: sql.ErrNoRows},
			expectedStatus: http.StatusNotFound,
			expectedResponse: Response{Error: Error{
				Status:  http.StatusNotFound,
				Message: "cat id fe271e7e-83ca-477b-92fc-d0c3fa602d7d not found",
			}},
			expectedMockCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := model.NewMockStorage(ctrl)
			s.EXPECT().
				Delete(gomock.Any()).
				Return(tt.mockResponse.one).
				Times(tt.expectedMockCalls)
			a.Storage = s

			req := httptest.NewRequest(http.MethodDelete, tt.request, nil)
			response := httptest.NewRecorder()
			a.Router.ServeHTTP(response, req)

			if response.Code != tt.expectedStatus {
				t.Errorf("unxpected status code: got %d, expected %d", response.Code, tt.expectedStatus)
			}

			var got Response
			err := json.Unmarshal(response.Body.Bytes(), &got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expectedResponse) {
				t.Errorf("unxpected response got: got %v, expected %v",
					got, tt.expectedResponse)
			}

		})
	}
}

func TestApp_GetCats(t *testing.T) {
	var a App
	a.BootstrapServer()

	tests := []struct {
		description string
		// given
		request string

		mockResponse struct {
			one []byte
			two error
		}
		// then
		expectedStatus    int
		expectedResponse  interface{}
		expectedMockCalls int
	}{
		{
			description:    "successful",
			request:        "/cats/v1/cats",
			expectedStatus: http.StatusOK,
			expectedResponse: []model.Cat{model.Cat{
				Name:  "cat-1",
				Color: "color-1",
				Age:   1},
				model.Cat{
					Name:  "cat-2",
					Color: "color-2",
					Age:   2},
			},
			mockResponse: struct {
				one []byte
				two error
			}{
				one: []byte(`[{"name": "cat-1", "color": "color-1", "age": 1},{"name": "cat-2", "color": "color-2", "age": 2}]`),
				two: nil},
			expectedMockCalls: 1,
		},
		{
			description:      "no items available",
			request:          "/cats/v1/cats",
			expectedStatus:   http.StatusOK,
			expectedResponse: []model.Cat{},
			mockResponse: struct {
				one []byte
				two error
			}{
				one: []byte(`[]`),
				two: sql.ErrNoRows},
			expectedMockCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := model.NewMockStorage(ctrl)
			s.EXPECT().
				SelectAll(gomock.Any(), gomock.Any()).
				Return(tt.mockResponse.one, tt.mockResponse.two).
				Times(tt.expectedMockCalls)
			a.Storage = s

			req := httptest.NewRequest(http.MethodGet, tt.request, nil)
			response := httptest.NewRecorder()
			a.Router.ServeHTTP(response, req)

			if response.Code != tt.expectedStatus {
				t.Errorf("unxpected status code: got %d, expected %d", response.Code, tt.expectedStatus)
			}

			switch tt.expectedResponse.(type) {
			case []model.Cat:
				var got []model.Cat
				err := json.Unmarshal(response.Body.Bytes(), &got)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, tt.expectedResponse.([]model.Cat)) {
					t.Errorf("unxpected response got: got %v, expected %v",
						got, tt.expectedResponse)
				}
			case Response:
				var got Response
				err := json.Unmarshal(response.Body.Bytes(), &got)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, tt.expectedResponse.(Response)) {
					t.Errorf("unxpected response got: got %v, expected %v",
						got, tt.expectedResponse)
				}
			default:
				t.Errorf("unexpected type: %v", tt.expectedResponse)
			}
		})
	}
}

func TestApp_Health(t *testing.T) {
	var a App
	a.BootstrapServer()

	tests := []struct {
		description string
		// given
		request string
		// then
		expectedStatus   int
		expectedResponse interface{}
	}{{
		description:      "api reports healthy",
		request:          "/cats/v1/health",
		expectedStatus:   200,
		expectedResponse: health{Status: "Pets is up and available"},
	}}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.request, nil)
			response := httptest.NewRecorder()
			a.Router.ServeHTTP(response, req)

			if response.Code != tt.expectedStatus {
				t.Errorf("unxpected status code: got %d, expected %d", response.Code, tt.expectedStatus)
			}

			var got health
			err := json.Unmarshal(response.Body.Bytes(), &got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expectedResponse.(health)) {
				t.Errorf("unxpected response got: got %v, expected %v",
					got, tt.expectedResponse)
			}
		})
	}
}

func TestApp_UpdateCat(t *testing.T) {
	var a App
	a.BootstrapServer()

	tests := []struct {
		description string
		// given
		request      string
		requestBody  []byte
		mockResponse struct {
			one error
		}
		// then
		expectedStatus    int
		expectedResponse  interface{}
		expectedMockCalls int
	}{
		{
			description:    "successful",
			request:        "/cats/v1/fe271e7e-83ca-477b-92fc-d0c3fa602d7d",
			expectedStatus: http.StatusOK,
			requestBody:    []byte(`{"name":"cat-1","color":"color-1","age":1}`),
			expectedResponse: Response{
				Result: "success",
			},
			mockResponse: struct {
				one error
			}{
				one: nil},
			expectedMockCalls: 1,
		},
		{
			description:    "cat not found",
			request:        "/cats/v1/1DF8D025-E80D-425E-94BD-900D6738E7BC",
			requestBody:    []byte(`{"name":"cat-1","color":"color-1","age":1}`),
			expectedStatus: http.StatusNotFound,
			expectedResponse: Response{
				Error: Error{
					Status:  http.StatusNotFound,
					Message: "cat id 1DF8D025-E80D-425E-94BD-900D6738E7BC not found"},
			},
			mockResponse: struct {
				one error
			}{
				one: sql.ErrNoRows,
			},
			expectedMockCalls: 1,
		},
		{
			description:    "bad request invalid json",
			request:        "/cats/v1/1DF8D025-E80D-425E-94BD-900D6738E7BC",
			requestBody:    []byte(`{"name":"cat-1","color":"color-1","`),
			expectedStatus: http.StatusBadRequest,
			expectedResponse: Response{
				Error: Error{
					Status:  http.StatusBadRequest,
					Message: "invalid json in request body"},
			},
			mockResponse: struct {
				one error
			}{
				one: sql.ErrNoRows,
			},
			expectedMockCalls: 0,
		},
		// todo add test case or scenario where an error reading body occurs
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := model.NewMockStorage(ctrl)
			s.EXPECT().
				Update(gomock.Any(), gomock.Any()).
				Return(tt.mockResponse.one).
				Times(tt.expectedMockCalls)
			a.Storage = s

			req := httptest.NewRequest(http.MethodPut, tt.request, bytes.NewBuffer(tt.requestBody))
			response := httptest.NewRecorder()
			a.Router.ServeHTTP(response, req)

			if response.Code != tt.expectedStatus {
				t.Errorf("unxpected status code: got %d, expected %d", response.Code, tt.expectedStatus)
			}

			var got Response
			err := json.Unmarshal(response.Body.Bytes(), &got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expectedResponse.(Response)) {
				t.Errorf("unxpected response got: got %v, expected %v",
					got, tt.expectedResponse)
			}
		})
	}
}
