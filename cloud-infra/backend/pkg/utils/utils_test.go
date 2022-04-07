package utils

import (
	"backend/pkg/types"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestValidateFile(t *testing.T) {
	var tc = []struct {
		FileName    string
		MIMEType    string
		ExpectError bool // true mean error should NOT be nil
	}{
		{"../test_files/invalidExtension.json", "application/pdf", true}, // File with invalid Extension
		{"../test_files/sample.pdf", "notvalid/content-type", true},      // File with pdf extension but invalid content-type
		{"../test_files/sample.pdf", "application/pdf", false},           // valid pdf file
		{"../test_files/invalid.stl", "invalid/extension", true},         // file with stl extension but not valid stl file
		{"../test_files/goku_ss.stl", "notused/content-type", false},     // valid stl file
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v", i), func(t *testing.T) {
			file, err := os.Open(tt.FileName)
			if err != nil {
				t.Errorf("File %s not found in test folder", tt.FileName)
			}
			defer file.Close()

			err = ValidateFile(file, tt.FileName, tt.MIMEType)

			if tt.ExpectError && err == nil {
				t.Errorf("Expected error with file %s and Content-Type %s but got no error", tt.FileName, tt.MIMEType)
			} else if !tt.ExpectError && err != nil {
				t.Errorf("Did not expect error with file %s and Content-Type %s but got %s", tt.FileName, tt.MIMEType, err.Error())
			}
		})
	}
}

func TestValidateMaterial(t *testing.T) {
	var tc = []struct {
		Material string
		err      error
	}{
		{"HR PA 12GB", nil}, // Valid Material
		{"plastic", errors.New("invalid material received")}, // Invalid Material
		{"", errors.New("invalid material received")},        // Empty
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v", i), func(t *testing.T) {
			err := ValidateMaterial(tt.Material)
			fmt.Println(err == nil)
			if err != nil && (tt.err == nil || err.Error() != tt.err.Error()) {
				t.Fail()
			}
			if err == nil && tt.err != nil {
				t.Fail()
			}
		})
	}
}

func TestDevicesToPublicJSON(t *testing.T) {
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

	deviceNoName := types.Device{
		IP:         "placeholder",
		DeviceUUID: "placeholder",
	}

	var tc = []struct {
		devices        []types.Device
		expectedOutput []byte
		name           string
	}{
		{[]types.Device{}, []byte("[]"), "No devices"},
		{[]types.Device{deviceEmpty}, []byte("[]"), "Empty device"},
		{[]types.Device{deviceNoModel}, []byte("[{\"Name\":\"deviceNoModelName\"}]"), "Device no model"},
		{[]types.Device{deviceFull}, []byte("[{\"Name\":\"deviceFullName\",\"Model\":\"deviceFullModel\"}]"), "Device full"},
		{[]types.Device{deviceNoName}, []byte("[]"), "Device no name"},
		{[]types.Device{deviceFull, deviceEmpty, deviceNoModel, deviceNoName}, []byte("[{\"Name\":\"deviceFullName\",\"Model\":\"deviceFullModel\"},{\"Name\":\"deviceNoModelName\"}]"), "Various devices"},
	}

	for i, tt := range tc {
		t.Run(fmt.Sprintf("Test %v: %s", i, tt.name), func(t *testing.T) {
			output := DevicesToPublicJSON(tt.devices)
			if !reflect.DeepEqual(output, tt.expectedOutput) {
				t.Errorf("Expected output %v, got %v", string(tt.expectedOutput), string(output))
			}
		})
	}
}
