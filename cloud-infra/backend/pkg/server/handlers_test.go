package server

import (
	"backend/pkg/mocks"
	"backend/pkg/types"
	"backend/pkg/utils"
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
	mockObjStorage := mocks.NewMockObjStorage(mockCtrl)

	mockDatabase := mocks.NewMockDatabase(mockCtrl)

	// The mocked queue will return nil as error when called with any value
	mockQueue.EXPECT().SendMessage(gomock.Any()).Return(nil).AnyTimes()

	// We assume database never returns an error and always gives a valid IP and UUID back
	mockDatabase.EXPECT().DeviceIPAndUUIDFromName(gomock.Any()).Return("127.0.0.1", "placeholderUUID", nil).AnyTimes()

	// We assume database insert message never return an error
	mockDatabase.EXPECT().InsertMessage(gomock.Any()).Return(nil).AnyTimes()

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, mockDatabase, router)
	server.Routes()

	var tc = []struct {
		body               []byte
		expectedStatusCode int
		testName           string
	}{
		{nil, http.StatusBadRequest, "Empty request body"},
		{[]byte(`{"type":"HEARTBEAT"}`), http.StatusBadRequest, "Message is empty"},
		{[]byte(`{"type":"HEARTBEAT", "message":"placeholder"}`), http.StatusBadRequest, "Device Name is empty"},
		{[]byte(`{"message":"placeholder", "type":"JOB", "DeviceName":"device"}`), http.StatusBadRequest, "Type is incorrect"},
		{[]byte(`{"message":"placeholder", "type":"HEARTBEAT", "DeviceName":"device"}`), http.StatusOK, "Good request"},
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
	mockObjStorage := mocks.NewMockObjStorage(mockCtrl)
	mockDatabase := mocks.NewMockDatabase(mockCtrl)

	// The mocked queue will return nil as error when called with any value
	mockQueue.EXPECT().SendMessage(gomock.Any()).Return(nil).AnyTimes()

	// We assume database never returns an error and always gives a valid IP and UUID back
	mockDatabase.EXPECT().DeviceIPAndUUIDFromName(gomock.Any()).Return("127.0.0.1", "placeholderUUID", nil).AnyTimes()

	// The mocked object storage will return nil as error when called with any values
	mockObjStorage.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// We assume database insert message never return an error
	mockDatabase.EXPECT().InsertMessage(gomock.Any()).Return(nil).AnyTimes()

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, mockDatabase, router)
	server.Routes()

	var tc = []struct {
		data               []byte
		file               string
		expectedStatusCode int
		testName           string
	}{
		{nil, "", http.StatusBadRequest, "Empty request"}, // empty request
		{[]byte(`{"type":"HEARTBEAT", "DeviceName" : "device", "material":"HR PA 12GB"}`), "", http.StatusBadRequest, "Wrong type"}, // Wrong type
		{[]byte(`{"type":"JOB", "DeviceName" : "device"}`), "", http.StatusBadRequest, "Missing material"},
		{[]byte(`{"type":"JOB", "material" : "HR PA 12GB"}`), "", http.StatusBadRequest, "Missing device name"},                   // Missing field
		{[]byte(`{"type":"JOB", "DeviceName" : "device", "material":"HR PA 12GB"}`), "", http.StatusBadRequest, "Missing file"},   // Missing file
		{[]byte(`{"type":"JOB", "DeviceName" : "device", "material":"HR PA 12GB"}`), "sample.pdf", http.StatusOK, "Good request"}, // All good
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if tt.data != nil {
				fw, _ := writer.CreateFormField("data")
				_, err := io.Copy(fw, strings.NewReader(string(tt.data)))
				if err != nil {
					t.Errorf("Error while copying data in test %v: %v", tt.testName, err)
				}
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
				_, err = io.Copy(fw, file)
				if err != nil {
					t.Errorf("Error while copying data in test %v: %v", tt.testName, err)
				}
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

func CustomCreateFormFile(w *multipart.Writer, fieldName string, fileName string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}

func TestUpload(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockQueue := mocks.NewMockQueue(mockCtrl)

	// Object Storage not used in Upload Handler but we need to pass one to the server struct
	mockObjStorage := mocks.NewMockObjStorage(mockCtrl)

	mockDatabase := mocks.NewMockDatabase(mockCtrl)

	// We assume database never returns an error and always gives a valid IP and UUID back
	mockDatabase.EXPECT().DeviceIPAndUUIDFromName(gomock.Any()).Return("127.0.0.1", "placeholderUUID", nil).AnyTimes()

	// The mocked queue will return nil as error when called with any value
	mockQueue.EXPECT().SendMessage(gomock.Any()).Return(nil).AnyTimes()

	// We assume database insert message never return an error
	mockDatabase.EXPECT().InsertMessage(gomock.Any()).Return(nil).AnyTimes()

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, mockDatabase, router)
	server.Routes()

	var tc = []struct {
		body               []byte
		contentType        string
		expectedStatusCode int
		testName           string
	}{
		{[]byte(`placeholder`), "text/plain", http.StatusBadRequest, "Invalid content type"},
		{[]byte(`()!!)(""??!!))`), "application/json", http.StatusBadRequest, "Invalid JSON format"},
		{[]byte(`{"DeviceName" : "device", "type":"JOB", "UploadInfo":"Jobs"}`), "application/json", http.StatusBadRequest, "Invalid Type field"},
		{[]byte(`{"type":"UPLOAD", "UploadInfo":"Jobs"}`), "application/json", http.StatusBadRequest, "Missing device name field"},
		{[]byte(`{"DeviceName" : "device", "type":"UPLOAD", "UploadInfo":"placeholder"}`), "application/json", http.StatusBadRequest, "Invalid Upload Info field"},
		{[]byte(`{"DeviceName" : "device", "type":"UPLOAD", "UploadInfo":"Jobs"}`), "application/json", http.StatusOK, "Good request"},
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			req := httptest.NewRequest("POST", "/upload", bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}
}

func TestUploadIdentification(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Mocked queue not used in this handler but we need to pass one to the server struct
	mockQueue := mocks.NewMockQueue(mockCtrl)
	mockObjStorage := mocks.NewMockObjStorage(mockCtrl)

	// The mocked object storage will return nil as error when called with any values
	mockObjStorage.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, nil, router)
	server.Routes()

	var tc = []struct {
		body               []byte
		contentType        string
		expectedStatusCode int
		deviceName         string
		testName           string
	}{
		{[]byte(`{}`), "text/plain", http.StatusBadRequest, "deviceName", "Invalid content type"},
		{[]byte(`()!!)(""??!!))`), "application/json", http.StatusBadRequest, "deviceName", "Invalid JSON format"},
		{[]byte(`{"identification": "placeholder"}`), "application/json", http.StatusBadRequest, "", "Empty device name"},
		{[]byte(`{"identification": "placeholder"}`), "application/json", http.StatusOK, "deviceName", "Good request"},
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			req := httptest.NewRequest("POST", "/uploadIdentification", bytes.NewBuffer(tt.body))
			req.Header.Set("X-Device", tt.deviceName)
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}
}

func TestUploadJobs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Mocked queue not used in this handler but we need to pass one to the server struct
	mockQueue := mocks.NewMockQueue(mockCtrl)
	mockObjStorage := mocks.NewMockObjStorage(mockCtrl)

	// The mocked object storage will return nil as error when called with any values
	mockObjStorage.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, nil, router)
	server.Routes()

	var tc = []struct {
		body               []byte
		contentType        string
		expectedStatusCode int
		deviceName         string
		testName           string
	}{
		{[]byte(`{}`), "text/plain", http.StatusBadRequest, "deviceName", "Invalid content type"},
		{[]byte(`()!!)(""??!!))`), "application/json", http.StatusBadRequest, "deviceName", "Invalid JSON format"},
		{[]byte(`{"jobs": "placeholder"}`), "application/json", http.StatusBadRequest, "", "Empty device name"},
		{[]byte(`{"jobs": "placeholder"}`), "application/json", http.StatusOK, "deviceName", "Good request"},
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			req := httptest.NewRequest("POST", "/uploadJobs", bytes.NewBuffer(tt.body))
			req.Header.Set("X-Device", tt.deviceName)
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}
}

func TestAvailableInformation(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Mocked queue not used in this handler but we need to pass one to the server struct
	mockQueue := mocks.NewMockQueue(mockCtrl)
	mockObjStorage := mocks.NewMockObjStorage(mockCtrl)

	var validInformation types.Information
	validInformation.Identification = []string{"Identification-127_0_0_1.json", "Identification-192_168_0_1.json"}
	validInformation.Jobs = []string{"Jobs-127_0_0_1.json"}

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, nil, router)
	server.Routes()

	var tc = []struct {
		returnedError      error
		expectedStatusCode int
		testName           string
	}{
		{nil, 200, "All fine"},
		{fmt.Errorf("Error getting available information"), 500, "Error while getting available information"},
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {

			mockObjStorage.EXPECT().AvailableInformation().Return(validInformation, tt.returnedError).Times(1)
			req := httptest.NewRequest("GET", "/availableInformation", nil)

			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}

}

func TestGetInformationFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Mocked queue not used in this handler but we need to pass one to the server struct
	mockQueue := mocks.NewMockQueue(mockCtrl)
	mockObjStorage := mocks.NewMockObjStorage(mockCtrl)

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, nil, router)
	server.Routes()

	var tc = []struct {
		requestAdditionalPath    string
		GetFileMockReturnedError error
		expectedStatusCode       int
		usesObjStorage           bool
		testName                 string
	}{
		{"", nil, 400, false, "Missing URL parameter"},
		{"?archive=Jobs-192_2_1_1.json", nil, 400, false, "Invalid URL parameter key"},
		{"?file=jobs-192_2_1_1.json", nil, 400, false, "Invalid requested file prefix"},
		{"?file=Jobs-192_2_1_1.pdf", nil, 400, false, "Invalid requested file suffix"},
		{"?file=Jobs-192_2_1_1.json", fmt.Errorf("Error on GetFile"), 500, true, "Server error while getting requested file"},
		{"?file=Jobs-192_2_1_1.json", nil, 200, true, "All good"},
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			if tt.usesObjStorage {
				mockObjStorage.EXPECT().GetFile(gomock.Any(), gomock.Any()).Return(tt.GetFileMockReturnedError).Times(1)
			}

			req := httptest.NewRequest("GET", "/getInformationFile"+tt.requestAdditionalPath, nil)

			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}
}

func TestGetPublicDevices(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockQueue := mocks.NewMockQueue(mockCtrl)

	// Object Storage not used in Upload Handler but we need to pass one to the server struct
	mockObjStorage := mocks.NewMockObjStorage(mockCtrl)

	mockDatabase := mocks.NewMockDatabase(mockCtrl)

	router := mux.NewRouter()
	server := NewServer(mockQueue, mockObjStorage, mockDatabase, router)
	server.Routes()

	deviceEmpty := types.Device{}

	deviceFull := types.Device{
		IP:         "placeholder",
		DeviceUUID: "placeholder",
		Name:       "deviceFullName",
		Model:      "deviceFullModel",
	}

	deviceNoModel := types.Device{
		IP:         "placeholder",
		DeviceUUID: "placeholder",
		Name:       "deviceNoModelName",
	}

	// function utils.DevicesToPublicJSON is tested on its own so we assume here that it works fine

	var tc = []struct {
		returnedDevices      []types.Device
		databaseError        error
		expectedStatusCode   int
		expectedResponseBody []byte
		testName             string
	}{
		{[]types.Device{}, nil, 200, []byte(utils.DevicesToPublicJSON([]types.Device{})), "OK no devices"},
		{[]types.Device{deviceEmpty, deviceFull, deviceNoModel}, nil, 200, []byte(utils.DevicesToPublicJSON([]types.Device{deviceEmpty, deviceFull, deviceNoModel})), "OK with devices"},
		{[]types.Device{}, fmt.Errorf("error"), 500, []byte{}, "Server error"},
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			mockDatabase.EXPECT().GetDevices().Return(tt.returnedDevices, tt.databaseError).Times(1)
			req := httptest.NewRequest("GET", "/getPublicDevices", bytes.NewBuffer(nil))
			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
			body, _ := io.ReadAll(w.Result().Body)
			if string(body) != string(tt.expectedResponseBody) {
				t.Errorf("Expected response body %v, got %v", string(tt.expectedResponseBody), string(body))
			}
		})
	}
}

func TestNewDevice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockQueue := mocks.NewMockQueue(mockCtrl)

	mockObjStorage := mocks.NewMockObjStorage(mockCtrl)

	mockDatabase := mocks.NewMockDatabase(mockCtrl)

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, mockDatabase, router)
	server.Routes()

	var testCasesNoDBinvolved = []struct {
		body               []byte
		contentType        string
		expectedStatusCode int
		testName           string
	}{
		{[]byte(`placeholder`), "text/plain", http.StatusBadRequest, "Invalid content type"},
		{[]byte(`()!!)(""??!!))`), "application/json", http.StatusBadRequest, "Invalid JSON format"},
		{[]byte(`{"Name":"devName"}`), "application/json", http.StatusBadRequest, "Missing device IP field"},
		{[]byte(`{"IP":"127.0.0.1"}`), "application/json", http.StatusBadRequest, "Missing device name field"},
		{[]byte(`{"IP":"123456789","Name":"devName"}`), "application/json", http.StatusBadRequest, "Invalid IP provided"},
	}

	for i, tt := range testCasesNoDBinvolved {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			req := httptest.NewRequest("POST", "/devices", bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}

	var testCasesDBinvolved = []struct {
		body               []byte
		contentType        string
		expectedStatusCode int
		alreadyExists      bool
		insertError        error
		testName           string
	}{
		{[]byte(`{"IP":"127.0.0.1","Name":"devName"}`), "application/json", http.StatusBadRequest, true, fmt.Errorf("Device already exists"), "Device already exists"},
		{[]byte(`{"IP":"127.0.0.1","Name":"devName"}`), "application/json", http.StatusOK, false, nil, "Device inserted correctly"},
	}

	for i, tt := range testCasesDBinvolved {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			mockDatabase.EXPECT().DeviceExistWithNameAndIP(gomock.Any(), gomock.Any()).Return(tt.alreadyExists, nil).Times(1)

			if !tt.alreadyExists {
				mockDatabase.EXPECT().InsertDevice(gomock.Any()).Return(tt.insertError).Times(1)
			}

			req := httptest.NewRequest("POST", "/devices", bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}
}

func TestReceiveResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockQueue := mocks.NewMockQueue(mockCtrl)

	mockObjStorage := mocks.NewMockObjStorage(mockCtrl)

	mockDatabase := mocks.NewMockDatabase(mockCtrl)

	router := mux.NewRouter()

	server := NewServer(mockQueue, mockObjStorage, mockDatabase, router)
	server.Routes()

	var testCasesNoDBinvolved = []struct {
		contentType        string
		deviceUUID         string
		messageUUID        string
		body               []byte
		expectedStatusCode int
		testName           string
	}{
		{"text/plain", "placeholderDeviceUUID", "placeholderMessageUUID", []byte(`placeholder`), http.StatusBadRequest, "Invalid content type"},
		{"application/json", "invalidDeviceUUID", "placeholderMessageUUID", []byte(`placeholder`), http.StatusBadRequest, "Invalid device UUID"},
		{"application/json", "111c4951-31ba-4f8c-bca8-b17528810ee9", "invalidMessageUUID", []byte(`placeholder`), http.StatusBadRequest, "Invalid message UUID"},
		{"application/json", "111c4951-31ba-4f8c-bca8-b17528810ee9", "111c4951-31ba-4f8c-bca8-b17528810ee9", []byte(`()!!)(""??!!))`), http.StatusBadRequest, "Invalid JSON provided as body"},
		{"application/json", "111c4951-31ba-4f8c-bca8-b17528810ee9", "111c4951-31ba-4f8c-bca8-b17528810ee9", []byte(`{"Result":"SUCCESS"}`), http.StatusBadRequest, "Missing Timestamp in body"},
		{"application/json", "111c4951-31ba-4f8c-bca8-b17528810ee9", "111c4951-31ba-4f8c-bca8-b17528810ee9", []byte(`{"Timestamp": 1650795291931}`), http.StatusBadRequest, "Missing Result in body"},
		{"application/json", "111c4951-31ba-4f8c-bca8-b17528810ee9", "111c4951-31ba-4f8c-bca8-b17528810ee9", []byte(`{"Result":"SUCCESS", "Timestamp": "invalidTimestamp"}`), http.StatusBadRequest, "Invalid timestamp in body"},
	}

	for i, tt := range testCasesNoDBinvolved {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			url := "/responses" + "/" + tt.deviceUUID + "/" + tt.messageUUID
			req := httptest.NewRequest("POST", url, bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}

	var testCasesDBinvolved = []struct {
		contentType        string
		deviceUUID         string
		messageUUID        string
		body               []byte
		expectedStatusCode int
		insertError        error
		testName           string
	}{
		{"application/json", "111c4951-31ba-4f8c-bca8-b17528810ee9", "111c4951-31ba-4f8c-bca8-b17528810ee9", []byte(`{"Result":"SUCCESS", "Timestamp": 1650795291931}`), http.StatusInternalServerError, fmt.Errorf("Server error"), "Error while inserting result"},
		{"application/json", "111c4951-31ba-4f8c-bca8-b17528810ee9", "111c4951-31ba-4f8c-bca8-b17528810ee9", []byte(`{"Result":"SUCCESS", "Timestamp": 1650795291931}`), http.StatusOK, nil, "All good"}}

	for i, tt := range testCasesDBinvolved {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.testName), func(t *testing.T) {
			url := "/responses" + "/" + tt.deviceUUID + "/" + tt.messageUUID
			mockDatabase.EXPECT().InsertResult(gomock.Any()).Return(tt.insertError).Times(1)
			req := httptest.NewRequest("POST", url, bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, req)
			if w.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected code %v, got %v", tt.expectedStatusCode, w.Result().StatusCode)
			}
		})
	}

}
