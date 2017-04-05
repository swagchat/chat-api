package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/go-zoo/bone"
)

type testRecord struct {
	testNo         int
	in             string
	out            string
	httpStatusCode int
}

func TestMain(m *testing.M) {
	datastoreProvider := datastore.GetProvider()
	err := datastoreProvider.Connect()
	if err != nil {
		log.Println(err.Error())
	}
	datastoreProvider.Init()
	Mux = bone.New().Prefix("/" + utils.API_VERSION)
	SetUserMux()
	SetRoomMux()
	SetRoomUserMux()
	testRC := m.Run()
	err = datastoreProvider.DropDatabase()
	if err != nil {
		log.Println(err.Error())
	}
	os.Exit(testRC)
}

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
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
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
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
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
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","customData":{},"isPublic":true,"created":[0-9]+,"modified":[0-9]+}$`,
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
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
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
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 6,
			in: `
				{
					"name": "dennis room",
					"pictureUrl": "http://localhost/images/dennis_room.png",
					"informationUrl": "http://localhost/dennis_room",
					"customData": {"key": "value"}
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","customData":{"key":"value"},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 7,
			in: `
				{
					"roomId": "custom-id",
					"name": "dennis room"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"roomId":"custom-id","name":"dennis room","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 8,
			in: `
				{
					"roomId": "custom-id",
					"name": "dennis room"
				}
			`,
			out:            `(?m)^{"title":"An error occurred while creating room item.","status":500,"detail":".*","errorName":"database-error"}$`,
			httpStatusCode: 500,
		},
		{
			testNo: 9,
			in: `
			`,
			out:            `(?m)^{"title":"Json parse error. \(Create room item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 10,
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
	log.Println(len(createRoomIds))

}

func TestGetRooms(t *testing.T) {
}

func TestGetRoom(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	if len(createRoomIds) != 7 {
		t.Fatalf("createRoomIds length error \n[expected]%d\n[result  ]%d", 7, len(createRoomIds))
		t.Failed()
	}

	testTable := []testRecord{
		{
			testNo:         1,
			in:             createRoomIds[0],
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         2,
			in:             createRoomIds[1],
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         3,
			in:             createRoomIds[2],
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","customData":{},"isPublic":true,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         4,
			in:             createRoomIds[3],
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         5,
			in:             createRoomIds[4],
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         6,
			in:             createRoomIds[5],
			out:            `(?m)^{"id":[0-9]+,"roomId":"[a-z0-9-]+","name":"dennis room","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","customData":{"key":"value"},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         7,
			in:             createRoomIds[6],
			out:            `(?m)^{"id":[0-9]+,"roomId":"custom-id","name":"dennis room","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         8,
			in:             "not-exist-room-id",
			out:            ``,
			httpStatusCode: 404,
		},
	}

	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/rooms/" + testRecord.in)

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
			in: `
				{
					"name": "Jeremy room"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"roomId":"custom-id","name":"Jeremy room","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 204,
		},
		{
			testNo: 2,
			in: `
				{
					"isPublic": true
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"roomId":"custom-id","name":"Jeremy room","customData":{},"isPublic":true,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 204,
		},
		{
			testNo: 3,
			in: `
				{
					"isPublic": false
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"roomId":"custom-id","name":"Jeremy room","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 204,
		},
		{
			testNo: 4,
			in: `
				{
					"pictureUrl": "http://localhost/images/jeremy.png"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"roomId":"custom-id","name":"Jeremy room","pictureUrl":"http://localhost/images/jeremy.png","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 204,
		},
		{
			testNo: 5,
			in: `
				{
					"informationUrl": "http://localhost/jeremy"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"roomId":"custom-id","name":"Jeremy room","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","customData":{},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 204,
		},
		{
			testNo: 6,
			in: `
				{
					"customData": {"key": "value"}
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"roomId":"custom-id","name":"Jeremy room","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","customData":{"key":"value"},"isPublic":false,"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 204,
		},
		{
			testNo: 7,
			in: `
			`,
			out:            `(?m)^{"title":"Json parse error. \(Update room item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 8,
			in: `
				json
			`,
			out:            `(?m)^{"title":"Json parse error. \(Update room item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		req, _ := http.NewRequest("PUT", ts.URL+"/"+utils.API_VERSION+"/rooms/custom-id", reader)
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		if testRecord.httpStatusCode == 204 {
			res, err = http.Get(ts.URL + "/" + utils.API_VERSION + "/rooms/custom-id")
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
			in:             "custom-id",
			out:            `(?m)^$`,
			httpStatusCode: 204,
		},
	}

	for _, testRecord := range testTable {
		req, _ := http.NewRequest("DELETE", ts.URL+"/"+utils.API_VERSION+"/rooms/"+testRecord.in, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		if testRecord.httpStatusCode == 204 {
			res, err = http.Get(ts.URL + "/" + utils.API_VERSION + "/rooms/" + testRecord.in)
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}
