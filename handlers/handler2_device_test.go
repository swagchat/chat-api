package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/swagchat/chat-api/utils"
)

var createDeviceIds []string

type deviceStruct struct {
	DeviceId string `json:"deviceId,omitempty"`
}

func TestPostDevice(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:   1,
			userId:   "custom-user-id-1",
			platform: "1",
			in: `
				{
					"token": "abc"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-1","platform":1,"token":"abc","notificationDeviceId":"abc"}$`,
			httpStatusCode: 201,
		},
		{
			testNo:   2,
			userId:   "custom-user-id-1",
			platform: "2",
			in: `
				{
					"token": "def"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-1","platform":2,"token":"def","notificationDeviceId":"def"}$`,
			httpStatusCode: 201,
		},
		{
			testNo:   3,
			userId:   "custom-user-id-1",
			platform: "2",
			in: `
				{
					"token": "def"
				}
			`,
			out:            `(?m)^{"title":"An error occurred while creating device item.","status":500,"detail":"UNIQUE constraint failed: device\.user_id, device.platform","errorName":"database-error"}$`,
			httpStatusCode: 500,
		},
		{
			testNo:   4,
			userId:   "custom-user-id-1",
			platform: "3",
			in: `
					{
						"token": "abc"
					}
				`,
			out:            `(?m)^{"title":"Request parameter error\. \(Create device item\)","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"device\.platform","reason":"platform is invalid\. Currently only 1\(iOS\) and 2\(Android\) are supported\."}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo:   5,
			userId:   "custom-user-id-1",
			platform: "1",
			in: `
				`,
			out:            `(?m)^{"title":"Json parse error. \(Create device item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo:   6,
			userId:   "custom-user-id-1",
			platform: "2",
			in: `
					json
				`,
			out:            `(?m)^{"title":"Json parse error. \(Create device item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		res, err := http.Post(ts.URL+"/"+utils.API_VERSION+"/users/"+testRecord.userId+"/devices/"+testRecord.platform, "application/json", reader)

		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("TestNo %d\nError by ioutil.ReadAll(): %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}

		if testRecord.httpStatusCode == 201 {
			device := &deviceStruct{}
			_ = json.Unmarshal(data, device)
			createDeviceIds = append(createDeviceIds, device.DeviceId)
		}
	}
}

func TestGetDevices(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:         1,
			userId:         "custom-user-id-1",
			out:            `(?m)^{"devices":[{"userId":"[a-z0-9-]+","platform":1,"token":"abc","notificationDeviceId":"abc"},{"userId":"[a-z0-9-]+","platform":2,"token":"def","notificationDeviceId":"def"}]}$`,
			httpStatusCode: 200,
		},
	}

	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/users/" + testRecord.userId + "/devices")
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}

func TestGetDevice(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	if len(createDeviceIds) != 2 {
		t.Fatalf("createDeviceIds length error \n[expected]%d\n[result  ]%d", 2, len(createDeviceIds))
		t.Failed()
	}
	testTable := []testRecord{
		{
			testNo:         1,
			userId:         "custom-user-id-1",
			platform:       "1",
			out:            `(?m)^{"userId":"custom-user-id-1","platform":1,"token":"abc","notificationDeviceId":"abc"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         2,
			userId:         "custom-user-id-1",
			platform:       "2",
			out:            `(?m)^{"userId":"custom-user-id-1","platform":2,"token":"def","notificationDeviceId":"def"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         3,
			userId:         "not-exist-device-id",
			platform:       "2",
			out:            ``,
			httpStatusCode: 404,
		},
	}
	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/users/" + testRecord.userId + "/devices/" + testRecord.platform)

		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("TestNo %d\nError by ioutil.ReadAll(): %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}

func TestPutDevice(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:   1,
			userId:   "custom-user-id-1",
			platform: "1",
			in: `
				{
					"token": "ghi"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-1","platform":1,"token":"ghi","notificationDeviceId":"ghi"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:   2,
			userId:   "custom-user-id-1",
			platform: "2",
			in: `
				{
					"token": "jkl"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-1","platform":2,"token":"jkl","notificationDeviceId":"jkl"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:   3,
			userId:   "custom-user-id-1",
			platform: "1",
			in: `
					{
						"token": ""
					}
				`,
			out:            `(?m)^{"title":"Request parameter error\. \(Create device item\)","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"token","reason":"token is required, but it's empty\."}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo:   4,
			userId:   "custom-user-id-1",
			platform: "1",
			in: `
					{
						"token": "ghi"
					}
				`,
			out:            `(?m)^$`,
			httpStatusCode: 304,
		},
		{
			testNo:   5,
			userId:   "custom-user-id-1",
			platform: "1",
			in: `
					{
						"token2": "ghi"
					}
				`,
			out:            `(?m)^{"title":"Request parameter error\. \(Create device item\)","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"token","reason":"token is required, but it's empty\."}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo:   6,
			userId:   "custom-user-id-1",
			platform: "1",
			in: `
				`,
			out:            `(?m)^{"title":"Json parse error. \(Create device item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo:   7,
			userId:   "custom-user-id-1",
			platform: "1",
			in: `
					json
				`,
			out:            `(?m)^{"title":"Json parse error. \(Create device item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo:   8,
			userId:   "not-exist-user-id",
			platform: "1",
			in: `
					{
						"token": "jkl"
					}
				`,
			out:            `(?m)^$`,
			httpStatusCode: 404,
		},
		{
			testNo:   9,
			userId:   "custom-user-id-1",
			platform: "3",
			in: `
					{
						"token": "jkl"
					}
				`,
			out:            `(?m)^$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		req, _ := http.NewRequest("PUT", ts.URL+"/"+utils.API_VERSION+"/users/"+testRecord.userId+"/devices/"+testRecord.platform, reader)
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}

func TestDeleteDevice(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:         1,
			userId:         "custom-user-id-1",
			platform:       "1",
			out:            `(?m)^$`,
			httpStatusCode: 204,
		},
		{
			testNo:         2,
			userId:         "not-exist-user-id",
			platform:       "1",
			out:            `(?m)^$`,
			httpStatusCode: 404,
		},
		{
			testNo:         3,
			userId:         "custom-user-id-1",
			platform:       "3",
			out:            `(?m)^$`,
			httpStatusCode: 404,
		},
	}

	for _, testRecord := range testTable {
		req, _ := http.NewRequest("DELETE", ts.URL+"/"+utils.API_VERSION+"/users/"+testRecord.userId+"/devices/"+testRecord.platform, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}
