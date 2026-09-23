package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearConfig() {
	plugin.ApiUrl = ""
	plugin.AuthToken = ""
}

func TestCreateAlert(t *testing.T) {
	clearConfig()
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	event.Check.Status = 1
	event.Metrics = nil

	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		expectedBody := `"message":"entity1/check1"`
		assert.Contains(string(body), expectedBody)
		w.Header().Add("X-RateLimit-State", "OK")
		w.WriteHeader(http.StatusAccepted)
		_, err := w.Write([]byte(`{"result": "Request will be processed",
                                           "took": 0.302,
                                           "requestId": "43a29c5c-3dbf-4fa4-9c26-f4f71023e120"
                                           }`))
		require.NoError(t, err)
	}))
	url, err := url.Parse(apiStub.URL)
	assert.NoError(err)

	plugin.ApiUrl = url.Host
	plugin.AuthToken = "test_token"

	plugin.MessageTemplate = "{{.Entity.Name}}/{{.Check.Name}}"
	plugin.MessageLimit = 200
	title, _, _ := parseEventKeyTags(event)
	assert.Equal("entity1/check1", title)
	alertClient, err := alert.NewClient(&client.Config{
		ApiKey:         plugin.AuthToken,
		OpsGenieAPIURL: client.ApiUrl(plugin.ApiUrl),
	})

	assert.NoError(err)

	err = createIncident(alertClient, event)
	assert.NoError(err)
}

func TestGetNote(t *testing.T) {
	clearConfig()
	event := corev2.FixtureEvent("foo", "bar")
	eventJSON, err := json.Marshal(event)
	assert.NoError(t, err)
	note, err := getNote(event)
	assert.NoError(t, err)
	assert.Contains(t, note, "Event data update:\n\n")
	assert.Contains(t, note, string(eventJSON))
}

func TestParseEventKeyTags(t *testing.T) {
	clearConfig()
	event := corev2.FixtureEvent("foo", "bar")
	_, err := json.Marshal(event)
	assert.NoError(t, err)
	plugin.MessageTemplate = "{{.Entity.Name}}/{{.Check.Name}}"
	plugin.MessageLimit = 100
	plugin.TagsTemplates = []string{"{{.Entity.Name}}", "{{.Check.Name}}", "{{.Entity.Namespace}}", "{{.Entity.EntityClass}}"}
	title, alias, tags := parseEventKeyTags(event)
	assert.Contains(t, title, "foo")
	assert.Contains(t, alias, "foo")
	assert.Contains(t, tags, "foo")
}

func TestParseDescription(t *testing.T) {
	clearConfig()
	event := corev2.FixtureEvent("foo", "bar")
	event.Check.Output = "Check OK"
	_, err := json.Marshal(event)
	assert.NoError(t, err)
	plugin.DescriptionTemplate = "{{.Check.Output}}"
	plugin.DescriptionLimit = 100
	description := parseDescription(event)
	assert.Equal(t, description, "Check OK")
}

func TestParseDetails(t *testing.T) {
	clearConfig()
	event := corev2.FixtureEvent("foo", "bar")
	event.Check.Output = "Check OK"
	event.Check.State = "passing"
	event.Check.Status = 0
	_, err := json.Marshal(event)
	assert.NoError(t, err)
	det := parseDetails(event)
	assert.Equal(t, det["output"], "Check OK")
	assert.Equal(t, det["state"], "passing")
	assert.Equal(t, det["status"], "0")
}

func TestEventPriority(t *testing.T) {
	clearConfig()
	testcases := []struct {
		myPriority         string
		mismatchedPriority string
		alertsPriority     alert.Priority
	}{
		{"P1", "P2", alert.P1},
		{"P2", "P3", alert.P2},
		{"P3", "P4", alert.P3},
		{"P4", "P5", alert.P4},
		{"P5", "P1", alert.P5},
		{"Default", "P4", alert.P3},
	}

	for _, tc := range testcases {
		assert := assert.New(t)
		plugin.Priority = tc.myPriority
		priority := eventPriority()
		expectedValue := tc.alertsPriority
		assert.Equal(priority, expectedValue)
	}
}

func TestCheckArgs(t *testing.T) {
	clearConfig()
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	assert.Error(checkArgs(event))
	plugin.AuthToken = "Testing"
	assert.Error(checkArgs(event))
	plugin.Team = "Testing"
	assert.NoError(checkArgs(event))
}

func TestStringInSlice(t *testing.T) {
	clearConfig()
	testSlice := []string{"foo", "bar", "test"}
	testString := "test"
	testResult := stringInSlice(testString, testSlice)
	assert.True(t, testResult)
	testString = "giraffe"
	testResult = stringInSlice(testString, testSlice)
	assert.False(t, testResult)
}

func TestTrim(t *testing.T) {
	clearConfig()
	testString := "This string is 33 characters long"
	assert.Equal(t, trim(testString, 40), testString)
	assert.Equal(t, trim(testString, 4), "This")
}

func TestSwitchOpsgenieRegion(t *testing.T) {

	clearConfig()
	plugin.ApiUrl = "test-host:test-port"
	testVal := switchOpsgenieRegion()
	assert.Equal(t, testVal, client.ApiUrl(plugin.ApiUrl))
	clearConfig()
	expectedValueUS := client.API_URL
	expectedValueEU := client.API_URL_EU

	testUS := switchOpsgenieRegion()

	assert.Equal(t, testUS, expectedValueUS)

	plugin.APIRegion = "eu"

	testEU := switchOpsgenieRegion()

	assert.Equal(t, testEU, expectedValueEU)

	plugin.APIRegion = "EU"

	testEU2 := switchOpsgenieRegion()

	assert.Equal(t, testEU2, expectedValueEU)
}
