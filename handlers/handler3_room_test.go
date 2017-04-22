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

var createRoomIds []string

type roomStruct struct {
	RoomId string `json:"roomId,omitempty"`
}

func TestPostRoom(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			in: `
				{
					"name": "dennis room"
				}
			`,
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 2,
			in: `
				{
					"name": "dennis room",
					"isPublic": false
				}
			`,
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 3,
			in: `
				{
					"name": "dennis room",
					"isPublic": true
				}
			`,
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","metaData":{},"isPublic":true,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 4,
			in: `
				{
					"name": "dennis room",
					"pictureUrl": "http://localhost/images/dennis_room.png"
				}
			`,
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 5,
			in: `
				{
					"name": "dennis room",
					"pictureUrl": "http://localhost/images/dennis_room.png",
					"informationUrl": "http://localhost/dennis_room"
				}
			`,
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 6,
			in: `
				{
					"name": "dennis room",
					"pictureUrl": "http://localhost/images/dennis_room.png",
					"informationUrl": "http://localhost/dennis_room",
					"metaData": {"key": "value"}
				}
			`,
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","metaData":{"key":"value"},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 7,
			in: `
				{
					"roomId": "custom-room-id",
					"name": "dennis room"
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 8,
			in: `
				{
					"roomId": "custom-room-id-for-delete",
					"name": "dennis room"
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id-for-delete","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 9,
			in: `
				{
					"roomId": "custom-room-id",
					"name": "dennis room"
				}
			`,
			out:            `(?m)^{"title":"An error occurred while creating room item.","status":500,"detail":".*","errorName":"database-error"}$`,
			httpStatusCode: 500,
		},
		{
			testNo: 10,
			in: `
			`,
			out:            `(?m)^{"title":"Json parse error. \(Create room item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 11,
			in: `
				json
			`,
			out:            `(?m)^{"title":"Json parse error. \(Create room item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		res, err := http.Post(ts.URL+"/"+utils.API_VERSION+"/rooms", "application/json", reader)

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
			room := &roomStruct{}
			_ = json.Unmarshal(data, room)
			createRoomIds = append(createRoomIds, room.RoomId)
		}
	}
}

func TestGetRooms(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:         1,
			out:            `(?m)^{"rooms":[{"roomId":"[a-z0-9-]+","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+},{"roomId":"[a-z0-9-]+","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+},{"roomId":"[a-z0-9-]+","name":"dennis room","metaData":{},"isPublic":true,"created":[0-9]+,"modified":[0-9]+},{"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+},{"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+},{"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","metaData":{"key":"value"},"isPublic":false,"created":[0-9]+,"modified":[0-9]+},{"roomId":"custom-room-id","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+},{"roomId":"custom-room-id-for-delete","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}],"allCount":0}$`,
			httpStatusCode: 200,
		},
	}

	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/rooms")
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

func TestGetRoom(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	if len(createRoomIds) != 8 {
		t.Fatalf("createRoomIds length error \n[expected]%d\n[result  ]%d", 8, len(createRoomIds))
		t.Failed()
	}

	testTable := []testRecord{
		{
			testNo:         1,
			roomId:         createRoomIds[0],
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         2,
			roomId:         createRoomIds[1],
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         3,
			roomId:         createRoomIds[2],
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","metaData":{},"isPublic":true,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         4,
			roomId:         createRoomIds[3],
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         5,
			roomId:         createRoomIds[4],
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         6,
			roomId:         createRoomIds[5],
			out:            `(?m)^{"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","metaData":{"key":"value"},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         7,
			roomId:         createRoomIds[6],
			out:            `(?m)^{"roomId":"custom-room-id","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         8,
			roomId:         createRoomIds[7],
			out:            `(?m)^{"roomId":"custom-room-id-for-delete","name":"dennis room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         9,
			roomId:         "not-exist-room-id",
			out:            ``,
			httpStatusCode: 404,
		},
	}

	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/rooms/" + testRecord.roomId)

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

func TestPutRoom(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			roomId: "custom-room-id",
			in: `
				{
					"name": "Jeremy room"
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id","name":"Jeremy room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 2,
			roomId: "custom-room-id",
			in: `
				{
					"isPublic": true
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id","name":"Jeremy room","metaData":{},"isPublic":true,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 3,
			roomId: "custom-room-id",
			in: `
				{
					"isPublic": false
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id","name":"Jeremy room","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 4,
			roomId: "custom-room-id",
			in: `
				{
					"pictureUrl": "http://localhost/images/jeremy.png"
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id","name":"Jeremy room","pictureUrl":"http://localhost/images/jeremy.png","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 5,
			roomId: "custom-room-id",
			in: `
				{
					"informationUrl": "http://localhost/jeremy"
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id","name":"Jeremy room","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","metaData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 6,
			roomId: "custom-room-id",
			in: `
				{
					"metaData": {"key": "value"}
				}
			`,
			out:            `(?m)^{"roomId":"custom-room-id","name":"Jeremy room","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","metaData":{"key":"value"},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 7,
			roomId: "custom-room-id",
			in: `
			`,
			out:            `(?m)^{"title":"Json parse error. \(Update room item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 8,
			roomId: "custom-room-id",
			in: `
				json
			`,
			out:            `(?m)^{"title":"Json parse error. \(Update room item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 9,
			roomId: "not-exist-room-id",
			in: `
				{
					"name": "Not exist"
				}
			`,
			out:            `(?m)^$`,
			httpStatusCode: 404,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		req, _ := http.NewRequest("PUT", ts.URL+"/"+utils.API_VERSION+"/rooms/"+testRecord.roomId, reader)
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

func TestDeleteRoom(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:         1,
			roomId:         "custom-room-id-for-delete",
			out:            `(?m)^$`,
			httpStatusCode: 204,
		},
		{
			testNo:         2,
			roomId:         "not-exist-room-id",
			out:            `(?m)^$`,
			httpStatusCode: 404,
		},
	}

	for _, testRecord := range testTable {
		req, _ := http.NewRequest("DELETE", ts.URL+"/"+utils.API_VERSION+"/rooms/"+testRecord.roomId, nil)
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
