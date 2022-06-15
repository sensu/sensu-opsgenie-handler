package main

import (
	"encoding/json"
	"testing"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetNote(t *testing.T) {
	event := corev2.FixtureEvent("foo", "bar")
	eventJSON, err := json.Marshal(event)
	assert.NoError(t, err)
	note, err := getNote(event)
	assert.NoError(t, err)
	assert.Contains(t, note, "Event data update:\n\n")
	assert.Contains(t, note, string(eventJSON))
}

func TestParseEventKeyTags(t *testing.T) {
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
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	assert.Error(checkArgs(event))
	plugin.AuthToken = "Testing"
	assert.Error(checkArgs(event))
	plugin.Team = "Testing"
	assert.NoError(checkArgs(event))
}

func TestStringInSlice(t *testing.T) {
	testSlice := []string{"foo", "bar", "test"}
	testString := "test"
	testResult := stringInSlice(testString, testSlice)
	assert.True(t, testResult)
	testString = "giraffe"
	testResult = stringInSlice(testString, testSlice)
	assert.False(t, testResult)
}

func TestTrim(t *testing.T) {
	testString := "This string is 33 characters long"
	assert.Equal(t, trim(testString, 40), testString)
	assert.Equal(t, trim(testString, 4), "This")
}

func TestSwitchOpsgenieRegion(t *testing.T) {
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
