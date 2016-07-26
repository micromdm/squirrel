package munkiserver_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/micromdm/squirrel/munki/munki"
)

func TestCreateDevice(t *testing.T) {
	server, _ := newServer(t)
	defer server.Close()
	devices := []munki.Device{
		munki.Device{
			SerialNumber: "FOOBARBAZ1",
		},
	}

	for _, d := range devices {
		testCreateDeviceHTTP(t, server, d.SerialNumber, d, http.StatusOK)
		os.Remove("testdata/testrepo/devices/" + d.SerialNumber)
	}
}

type createDeviceRequest struct {
	SerialNumber string `plist:"serial_number" json:"serial_number"`
	munki.Device
}

func testCreateDeviceHTTP(t *testing.T, server *httptest.Server, serial string, device munki.Device, expectedStatus int) *munki.Manifest {
	client := http.DefaultClient
	theURL := server.URL + "/api/v1/devices"
	var req = &createDeviceRequest{
		SerialNumber: serial,
		Device:       device,
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	body := ioutil.NopCloser(bytes.NewBuffer(data))
	resp, err := client.Post(theURL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedStatus {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}

	return nil
}
