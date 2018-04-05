package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var createMessageIds []string

type messageStruct struct {
	MessageIds []string `json:"messageIds,omitempty"`
}

func TestPostMessages(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-1",
							"type": "text",
							"payload": {
								"text": "Welcome to swagchat!"
							}
						}
					]
				}
			`,
			out:            `(?m)^{"messageIds":["[a-z0-9-]+"]}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 2,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-1",
							"type": "text",
							"payload": {
								"text": "Hi custom-room-id-1!"
							}
						},
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-1",
							"type": "text",
							"payload": {
								"text": "How's it going?"
							}
						}
					]
				}
			`,
			out:            `(?m)^{"messageIds":["[a-z0-9-]+","[a-z0-9-]+"]}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 3,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-1",
							"type": "text",
							"payload": {
								"text": "Good!"
							}
						},
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-2",
							"type": "text",
							"payload": {
								"text": "Welcome to swagchat!"
							}
						}
					]
				}
			`,
			out:            `(?m)^{"messageIds":["[a-z0-9-]+","[a-z0-9-]+"]}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 4,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-2",
							"type": "image",
							"payload": {
								"mime": "image/png",
								"sourceUrl": "http://example.com/source.png",
								"thumbnailUrl": "http://example.com/thumbnail.png"
							}
						}
					]
				}
			`,
			out:            `(?m)^{"messageIds":["[a-z0-9-]+"]}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 5,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-2",
							"type": "not-exist-type",
							"payload": {}
						}
					]
				}
			`,
			out:            `(?m)^{"messageIds":["[a-z0-9-]+"]}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 6,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-2",
							"type": "text",
							"payload": {
								"text": "Bye!"
							}
						},
						{
							"roomId": "custom-room-id-1",
							"userId": "not-exist-user-id",
							"type": "text",
							"payload": {
								"text": "Welcome to swagchat!"
							}
						}
					]
				}
			`,
			out:            `(?m)^{"messageIds":["[a-z0-9-]+"],"errors":\[{"title":"Request parameter error\. \(Create message item\)","status":400,"errorName":"invalid\-param","invalidParams":\[{"name":"userId","reason":"userId is invalid\. Not exist user\."}\]}\]}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 7,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-2",
							"type": "text",
							"payload": {}
						}
					]
				}
			`,
			out:            `(?m)^{"errors":\[{"title":"Request parameter error\. \(Create message item\)","status":400,"errorName":"invalid\-param","invalidParams":\[{"name":"payload","reason":"Text type needs text\."}\]}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 8,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-2",
							"type": "image",
							"payload": {
								"mime": "image/png"
							}
						}
					]
				}
			`,
			out:            `(?m)^{"errors":\[{"title":"Request parameter error\. \(Create message item\)","status":400,"errorName":"invalid\-param","invalidParams":\[{"name":"payload","reason":"Image type needs mime and sourceUrl\."}\]}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 9,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-2",
							"type": "image",
							"payload": {
								"sourceUrl": "http://example.com/source.png"
							}
						}
					]
				}
			`,
			out:            `(?m)^{"errors":\[{"title":"Request parameter error\. \(Create message item\)","status":400,"errorName":"invalid\-param","invalidParams":\[{"name":"payload","reason":"Image type needs mime and sourceUrl\."}\]}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 10,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-2",
							"type": "not-exist-type"
						}
					]
				}
			`,
			out:            `(?m)^{"errors":\[{"title":"Request parameter error\. \(Create message item\)","status":400,"errorName":"invalid\-param","invalidParams":\[{"name":"payload","reason":"payload is empty\."}\]}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 11,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "custom-user-id-2",
							"type": "text"
						}
					]
				}
			`,
			out:            `(?m)^{"errors":\[{"title":"Request parameter error\. \(Create message item\)","status":400,"errorName":"invalid\-param","invalidParams":\[{"name":"payload","reason":"payload is empty\."}\]}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 12,
			in: `
				{
					"messages" : [
						{
							"roomId": "not-exist-room-id",
							"userId": "custom-user-id-1",
							"type": "text",
							"payload": {
								"text": "Welcome to swagchat!"
							}
						}
					]
				}
			`,
			out:            `(?m)^{"errors":\[{"title":"Request parameter error\. \(Create message item\)","status":400,"errorName":"invalid\-param","invalidParams":\[{"name":"roomId","reason":"roomId is invalid\. Not exist room\."}\]}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 13,
			in: `
				{
					"messages" : [
						{
							"roomId": "custom-room-id-1",
							"userId": "not-exist-user-id",
							"type": "text",
							"payload": {
								"text": "Welcome to swagchat!"
							}
						}
					]
				}
			`,
			out:            `(?m)^{"errors":\[{"title":"Request parameter error\. \(Create message item\)","status":400,"errorName":"invalid\-param","invalidParams":\[{"name":"userId","reason":"userId is invalid\. Not exist user\."}\]}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 14,
			in: `
			`,
			out:            `(?m)^{"title":"Json parse error. \(Create message item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 15,
			in: `
				json
			`,
			out:            `(?m)^{"title":"Json parse error. \(Create message item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		req, _ := http.NewRequest("POST", ts.URL+"/messages", reader)
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

		if testRecord.httpStatusCode == 201 {
			message := &messageStruct{}
			_ = json.Unmarshal(data, message)
			for _, messageId := range message.MessageIds {
				createMessageIds = append(createMessageIds, messageId)
			}
		}
	}
}

func TestGetMessage(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	if len(createMessageIds) != 8 {
		t.Fatalf("createMessageIds length error \n[expected]%d\n[result  ]%d", 8, len(createMessageIds))
		t.Failed()
	}

	testTable := []testRecord{
		{
			testNo:         1,
			messageId:      createMessageIds[0],
			out:            `(?m)^{"messageId":"[a-z0-9-]+","roomId":"custom-room-id-1","userId":"custom-user-id-1","type":"text","payload":{"text":"Welcome to swagchat\!"},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         2,
			messageId:      createMessageIds[1],
			out:            `(?m)^{"messageId":"[a-z0-9-]+","roomId":"custom-room-id-1","userId":"custom-user-id-1","type":"text","payload":{"text":"Hi custom-room-id-1\!"},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         3,
			messageId:      createMessageIds[2],
			out:            `(?m)^{"messageId":"[a-z0-9-]+","roomId":"custom-room-id-1","userId":"custom-user-id-1","type":"text","payload":{"text":"How\'s it going\?"},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         4,
			messageId:      createMessageIds[3],
			out:            `(?m)^{"messageId":"[a-z0-9-]+","roomId":"custom-room-id-1","userId":"custom-user-id-1","type":"text","payload":{"text":"Good\!"},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         5,
			messageId:      createMessageIds[4],
			out:            `(?m)^{"messageId":"[a-z0-9-]+","roomId":"custom-room-id-1","userId":"custom-user-id-2","type":"text","payload":{"text":"Welcome to swagchat\!"},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         6,
			messageId:      createMessageIds[5],
			out:            `(?m)^{"messageId":"[a-z0-9-]+","roomId":"custom-room-id-1","userId":"custom-user-id-2","type":"image","payload":{"mime":"image\/png","sourceUrl":"http\:\/\/example.com\/source\.png","thumbnailUrl":"http\:\/\/example.com\/thumbnail\.png"},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         7,
			messageId:      createMessageIds[6],
			out:            `(?m)^{"messageId":"[a-z0-9-]+","roomId":"custom-room-id-1","userId":"custom-user-id-2","type":"not-exist-type","payload":{},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         8,
			messageId:      createMessageIds[7],
			out:            `(?m)^{"messageId":"[a-z0-9-]+","roomId":"custom-room-id-1","userId":"custom-user-id-2","type":"text","payload":{"text":"Bye\!"},"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         9,
			messageId:      "not-exist-message-id",
			out:            ``,
			httpStatusCode: 404,
		},
	}

	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/messages/" + testRecord.messageId)

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
