package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/fairway-corp/swagchat-api/utils"
)

var createRoomUserIds []string

type roomUserStruct struct {
	RoomUserId string `json:"roomUserId,omitempty"`
}

func TestPutRoomUsers(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["custom-user-id-1"]
				}
			`,
			out:            `(?m)^{"roomUsers":\[{"roomId":"custom\-room\-id","userId":"custom\-user\-id\-1","unreadCount":0,"metaData":{},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}\]}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 2,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["custom-user-id-2","custom-user-id-3"]
				}
			`,
			out:            `(?m)^{"roomUsers":\[{"roomId":"custom\-room\-id","userId":"custom\-user\-id\-1","unreadCount":0,"metaData":{},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"roomId":"custom\-room\-id","userId":"custom\-user\-id\-2","unreadCount":0,"metaData":{},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"roomId":"custom\-room\-id","userId":"custom\-user\-id\-3","unreadCount":0,"metaData":{},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}\]}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 3,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": []
				}
			`,
			out:            `(?m)^{"title":"Request parameter error\. \(Create room's user list\)","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"userIds","reason":"Not set\."}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 4,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["not-exist-user-id"]
				}
			`,
			out:            `(?m)^{"title":"Request parameter error\. \(Create room's user list\)","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"userIds","reason":"It contains a userId that does not exist\."}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 5,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["custom-user-id", "not-exist-user-id"]
				}
			`,
			out:            `(?m)^{"title":"Request parameter error\. \(Create room's user list\)","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"userIds","reason":"It contains a userId that does not exist\."}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 6,
			roomId: "custom-room-id",
			in: `
				`,
			out:            `(?m)^{"title":"Json parse error. \(Adding room's user list\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 7,
			roomId: "custom-room-id",
			in: `
					json
				`,
			out:            `(?m)^{"title":"Json parse error. \(Adding room's user list\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 8,
			roomId: "not-exist-room-id",
			in: `
				{
					"userIds": ["custom-user-id"]
				}
				`,
			out:            `(?m)^$`,
			httpStatusCode: 404,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		req, _ := http.NewRequest("PUT", ts.URL+"/"+utils.API_VERSION+"/rooms/"+testRecord.roomId+"/users", reader)
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

func TestPutRoomUser(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			roomId: "custom-room-id",
			userId: "custom-user-id-1",
			in: `
				{
					"unreadCount": 100
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id","userId":"custom-user-id-1","unreadCount":100,"metaData":{},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 2,
			roomId: "custom-room-id",
			userId: "custom-user-id-1",
			in: `
				{
					"metaData": {"key":"value"}
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id","userId":"custom-user-id-1","unreadCount":100,"metaData":{"key":"value"},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 3,
			roomId: "custom-room-id",
			userId: "custom-user-id-1",
			in: `
				{
					"unreadCount": 200,
					"metaData": {"key2":"value2"}
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id","userId":"custom-user-id-1","unreadCount":200,"metaData":{"key2":"value2"},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 4,
			roomId: "custom-room-id",
			userId: "custom-user-id-1",
			in: `
				`,
			out:            `(?m)^{"title":"Json parse error. \(Update room's user item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 5,
			roomId: "custom-room-id",
			userId: "custom-user-id-1",
			in: `
					json
				`,
			out:            `(?m)^{"title":"Json parse error. \(Update room's user item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		req, _ := http.NewRequest("PUT", ts.URL+"/"+utils.API_VERSION+"/rooms/"+testRecord.roomId+"/users/"+testRecord.userId, reader)
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

func TestDeleteRoomUsers(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["custom-user-id-1"]
				}
			`,
			out:            `(?m)^$`,
			httpStatusCode: 200,
		},
		{
			testNo: 2,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["custom-user-id-2","custom-user-id-3"]
				}
			`,
			out:            `(?m)^$`,
			httpStatusCode: 200,
		},
		{
			testNo: 3,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["custom-user-id-2","custom-user-id-3"]
				}
			`,
			out:            `(?m)^$`,
			httpStatusCode: 200,
		},
		{
			testNo: 6,
			roomId: "custom-room-id",
			in: `
				`,
			out:            `(?m)^{"title":"Json parse error. \(Deleting room's user list\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 7,
			roomId: "custom-room-id",
			in: `
					json
				`,
			out:            `(?m)^{"title":"Json parse error. \(Deleting room's user list\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 8,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["not-exist-user-id"]
				}
				`,
			out:            `(?m)^{"title":"Request parameter error\. \(Create room's user list\)","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"userIds","reason":"It contains a userId that does not exist\."}\]}$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		req, _ := http.NewRequest("DELETE", ts.URL+"/"+utils.API_VERSION+"/rooms/"+testRecord.roomId+"/users", reader)
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

func TestPostRoomUsers(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["custom-user-id-1"]
				}
			`,
			out:            `(?m)^{"roomUsers":\[{"roomId":"custom\-room\-id","userId":"custom\-user\-id\-1","unreadCount":0,"metaData":{},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}]}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 2,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["custom-user-id-1","custom-user-id-2"]
				}
			`,
			out:            `(?m)^{"roomUsers":\[{"roomId":"custom\-room\-id","userId":"custom\-user\-id\-1","unreadCount":0,"metaData":{},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"roomId":"custom\-room\-id","userId":"custom\-user\-id\-2","unreadCount":0,"metaData":{},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}]}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 3,
			roomId: "custom-room-id",
			in: `
				`,
			out:            `(?m)^{"title":"Json parse error. \(Create room's user list\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 4,
			roomId: "custom-room-id",
			in: `
					json
				`,
			out:            `(?m)^{"title":"Json parse error. \(Create room's user list\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 5,
			roomId: "custom-room-id",
			in: `
				{
					"userIds": ["not-exist-user-id"]
				}
				`,
			out:            `(?m)^{"title":"Request parameter error\. \(Create room's user list\)","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"userIds","reason":"It contains a userId that does not exist\."}\]}$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		res, err := http.Post(ts.URL+"/"+utils.API_VERSION+"/rooms/"+testRecord.roomId+"/users", "application/json", reader)

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
			roomUser := &roomUserStruct{}
			_ = json.Unmarshal(data, roomUser)
			createRoomUserIds = append(createRoomUserIds, roomUser.RoomUserId)
		}
	}
}
