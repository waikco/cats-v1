package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/julienschmidt/httprouter"
	"github.com/waikco/cats-v1/conf"
	"github.com/waikco/cats-v1/model"
)

func TestApp_CreateCat(t *testing.T) {
	var a App
	a.BootstrapServer()

	tests := []struct {
		description       string
		requestBody       []byte
		expectedStatus    int
		expectedResponse  catResponse
		expectedErr       error
		expectedMockCalls int
	}{
		{
			description:    "successful pet creation",
			requestBody:    []byte(`{"name":"cat-1","color":"color-1","age":1}`),
			expectedStatus: http.StatusCreated,
			expectedResponse: catResponse{
				Result: "success",
				ID:     "fe271e7e-83ca-477b-92fc-d0c3fa602d7d",
			},
			expectedErr:       nil,
			expectedMockCalls: 1,
		},
		{
			description:    "unsuccessful pet creation from bad request",
			requestBody:    []byte(`"name":"cat-1","color":"color-1","age":`),
			expectedStatus: http.StatusBadRequest,
			expectedResponse: catResponse{
				Result: "invalid json in request body",
			},
			expectedErr:       nil,
			expectedMockCalls: 0,
		},
		{
			description:    "unsuccessful pet creation from internal server error",
			requestBody:    []byte(`{"name":"cat-1","color":"color-1","age":1}`),
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: catResponse{
				Result: "error storing cat",
			},
			expectedErr:       errors.New("database error"),
			expectedMockCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := model.NewMockStorage(ctrl)
			s.EXPECT().
				Insert(tt.requestBody).
				Return(tt.expectedResponse.ID, tt.expectedErr).Times(tt.expectedMockCalls)
			a.Storage = s

			req, _ := http.NewRequest("POST", "/cats/v1/", bytes.NewBuffer(tt.requestBody))
			response := httptest.NewRecorder()
			a.Router.ServeHTTP(response, req)

			var r catResponse
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

func TestApp_DeleteCat(t *testing.T) {
	type fields struct {
		Server  *http.Server
		Storage model.Storage
		Router  http.Handler
		Config  conf.Config
	}
	type args struct {
		w  http.ResponseWriter
		r  *http.Request
		ps httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestApp_GetCat(t *testing.T) {
	type fields struct {
		Server  *http.Server
		Storage model.Storage
		Router  http.Handler
		Config  conf.Config
	}
	type args struct {
		w  http.ResponseWriter
		r  *http.Request
		ps httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestApp_GetCats(t *testing.T) {
	type fields struct {
		Server  *http.Server
		Storage model.Storage
		Router  http.Handler
		Config  conf.Config
	}
	type args struct {
		w  http.ResponseWriter
		r  *http.Request
		ps httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestApp_Health(t *testing.T) {
	type fields struct {
		Server  *http.Server
		Storage model.Storage
		Router  http.Handler
		Config  conf.Config
	}
	type args struct {
		w  http.ResponseWriter
		r  *http.Request
		ps httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestApp_MassCreateCat(t *testing.T) {
	type fields struct {
		Server  *http.Server
		Storage model.Storage
		Router  http.Handler
		Config  conf.Config
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestApp_UpdateCat(t *testing.T) {
	type fields struct {
		Server  *http.Server
		Storage model.Storage
		Router  http.Handler
		Config  conf.Config
	}
	type args struct {
		w  http.ResponseWriter
		r  *http.Request
		ps httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
