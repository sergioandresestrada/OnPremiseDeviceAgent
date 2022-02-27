package server

import (
	"backend/pkg/mocks"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestHeartbeat(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockQueue := mocks.NewMockQueue(mockCtrl)

	// Object Storage not used in Heartbeat Handler but we need to pass one to the server struct
	mockObjStorage := mocks.NewMockObj_storage(mockCtrl)

	// The mocked queue will return nil as error when called with any value
	mockQueue.EXPECT().SendMessage(gomock.Any()).Return(nil).AnyTimes()

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, router)
	server.Routes()

	var tc = []struct {
		body               []byte
		expectedStatusCode int
	}{
		{nil, http.StatusBadRequest},                                               // empty request body
		{[]byte(`{"type":"HEARTBEAT"}`), http.StatusBadRequest},                    // message is empty
		{[]byte(`{"message":"placeholder", "type":"JOB"}`), http.StatusBadRequest}, // type is incorrect
		{[]byte(`{"message":"placeholder", "type":"HEARTBEAT"}`), http.StatusOK},   // all fine
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v", i), func(t *testing.T) {
			req := httptest.NewRequest("POST", "/heartbeat", bytes.NewBuffer(tt.body))
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}

}

func TestJob(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockQueue := mocks.NewMockQueue(mockCtrl)
	mockObjStorage := mocks.NewMockObj_storage(mockCtrl)

	// The mocked queue will return nil as error when called with any value
	mockQueue.EXPECT().SendMessage(gomock.Any()).Return(nil).AnyTimes()

	// The mocked object storage will return nil as error when called with any values
	mockObjStorage.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, router)
	server.Routes()

	var tc = []struct {
		data               []byte
		file               string
		expectedStatusCode int
	}{
		{nil, "", http.StatusBadRequest}, // empty request
		{[]byte(`{"message":"placeholder", "type":"HEARTBEAT"}`), "", http.StatusBadRequest},                   // Wrong type
		{[]byte(`{"type":"JOB"}`), "", http.StatusBadRequest},                                                  // Missing field
		{[]byte(`{"message":"placeholder", "type":"JOB"}`), "", http.StatusBadRequest},                         // Missing file
		{[]byte(`{"message":"placeholder", "type":"JOB"}`), "../test_files/sample.pdf", http.StatusBadRequest}, // All good
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v", i), func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if tt.data != nil {
				fw, _ := writer.CreateFormField("name")
				io.Copy(fw, strings.NewReader(string(tt.data)))
			}

			if tt.file != "" {
				fw, _ := writer.CreateFormFile("file", tt.file)
				file, err := os.Open(tt.file)
				if err != nil {
					t.Errorf("File %s not found in test folder", tt.file)
				}
				io.Copy(fw, file)
				file.Close()
			}
			writer.Close()

			req := httptest.NewRequest("POST", "/job", bytes.NewBuffer(body.Bytes()))
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}

}
