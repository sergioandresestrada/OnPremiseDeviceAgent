package server

import (
	"backend/pkg/mocks"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
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
		testName           string
	}{
		{nil, http.StatusBadRequest, "Empty request body"},
		{[]byte(`{"type":"HEARTBEAT"}`), http.StatusBadRequest, "Message is empty"},
		{[]byte(`{"message":"placeholder", "type":"JOB"}`), http.StatusBadRequest, "Type is incorrect"},
		{[]byte(`{"message":"placeholder", "type":"HEARTBEAT"}`), http.StatusOK, "Good request"},
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
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
		testName           string
	}{
		{nil, "", http.StatusBadRequest, "Empty request"}, // empty request
		{[]byte(`{"type":"HEARTBEAT", "IPAddress" : "127.0.0.1", "material":"HR PA 12GB"}`), "", http.StatusBadRequest, "Wrong type"}, // Wrong type
		{[]byte(`{"type":"JOB"}`), "", http.StatusBadRequest, "Missing field"},                                                        // Missing field
		{[]byte(`{"type":"JOB", "IPAddress" : "127.0.0.1", "material":"HR PA 12GB"}`), "", http.StatusBadRequest, "Missing file"},     // Missing file
		{[]byte(`{"type":"JOB", "IPAddress" : "127.0.0.1", "material":"HR PA 12GB"}`), "sample.pdf", http.StatusOK, "Good request"},   // All good
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if tt.data != nil {
				fw, _ := writer.CreateFormField("data")
				io.Copy(fw, strings.NewReader(string(tt.data)))
			}

			if tt.file != "" {
				var fw io.Writer
				switch filepath.Ext(tt.file) {
				case ".pdf":
					fw, _ = CustomCreateFormFile(writer, "file", tt.file, "application/pdf")
				default:
					fw, _ = writer.CreateFormFile("file", tt.file)
				}
				file, err := os.Open("../test_files/" + tt.file)
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

func CustomCreateFormFile(w *multipart.Writer, fieldName string, fileName string, content_type string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	h.Set("Content-Type", content_type)
	return w.CreatePart(h)
}

func TestUpload(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockQueue := mocks.NewMockQueue(mockCtrl)

	// Object Storage not used in Upload Handler but we need to pass one to the server struct
	mockObjStorage := mocks.NewMockObj_storage(mockCtrl)

	// The mocked queue will return nil as error when called with any value
	mockQueue.EXPECT().SendMessage(gomock.Any()).Return(nil).AnyTimes()

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, router)
	server.Routes()

	var tc = []struct {
		body               []byte
		content_type       string
		expectedStatusCode int
		testName           string
	}{
		{[]byte(`placeholder`), "text/plain", http.StatusBadRequest, "Invalid content type"},                                                                         // invalid content type
		{[]byte(`()!!)(""·!!))`), "application/json", http.StatusBadRequest, "Invalid JSON format"},                                                                  // invalid JSON format
		{[]byte(`{"IPAddress" : "127.0.0.1", "type":"JOB", "UploadInfo":"Jobs"}`), "application/json", http.StatusBadRequest, "Invalid Type field"},                  // Invalid Type field
		{[]byte(`{"IPAddress" : "999.999.0.0", "type":"UPLOAD", "UploadInfo":"Jobs"}`), "application/json", http.StatusBadRequest, "Invalid IP address"},             // Invalid IP field
		{[]byte(`{"IPAddress" : "127.0.0.1", "type":"UPLOAD", "UploadInfo":"placeholder"}`), "application/json", http.StatusBadRequest, "Invalid Upload Info field"}, // Invalid UploadInfo field
		{[]byte(`{"IPAddress" : "127.0.0.1", "type":"UPLOAD", "UploadInfo":"Jobs"}`), "application/json", http.StatusOK, "Good request"},                             // Invalid UploadInfo field
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			req := httptest.NewRequest("POST", "/upload", bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", tt.content_type)
			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}
}
