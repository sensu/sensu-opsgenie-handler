package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/sensu/sensu-go/types"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-plugin-sdk/templates"
)

const (
	notFound = "NOT FOUND"
	source   = "Sensu Go"
)

// Config represents the handler plugin config.
type Config struct {
	sensu.PluginConfig
	IncludeEventInNote  bool
	FullDetails         bool
	WithAnnotations     bool
	WithLabels          bool
	MessageLimit        int
	DescriptionLimit    int
	APIRegion           string
	AuthToken           string
	Team                string
	Priority            string
	SensuDashboard      string
	MessageTemplate     string
	DescriptionTemplate string
	Actions             []string
	TagsTemplates       []string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-opsgenie-handler",
			Short:    "The Sensu Go OpsGenie handler for incident management",
			Keyspace: "sensu.io/plugins/sensu-opsgenie-handler/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "region",
			Env:       "OPSGENIE_REGION",
			Argument:  "region",
			Shorthand: "r",
			Default:   "us",
			Usage:     "The OpsGenie API Region (us or eu), use default from OPSGENIE_REGION env var",
			Value:     &plugin.APIRegion,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "auth",
			Env:       "OPSGENIE_AUTHTOKEN",
			Argument:  "auth",
			Shorthand: "a",
			Default:   "",
			Secret:    true,
			Usage:     "The OpsGenie V2 API authentication token, use default from OPSGENIE_AUTHTOKEN env var",
			Value:     &plugin.AuthToken,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "team",
			Env:       "OPSGENIE_TEAM",
			Argument:  "team",
			Shorthand: "t",
			Default:   "",
			Usage:     "The OpsGenie V2 API Team, use default from OPSGENIE_TEAM env var",
			Value:     &plugin.Team,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "sensuDashboard",
			Env:       "OPSGENIE_SENSU_DASHBOARD",
			Argument:  "sensuDashboard",
			Shorthand: "s",
			Default:   "",
			Usage:     "The OpsGenie Handler will use it to create a source Sensu Dashboard URL. Use OPSGENIE_SENSU_DASHBOARD. Example: http://sensu-dashboard.example.local/c/~/n",
			Value:     &plugin.SensuDashboard,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "messageTemplate",
			Env:       "OPSGENIE_MESSAGE_TEMPLATE",
			Argument:  "messageTemplate",
			Shorthand: "m",
			Default:   "{{.Entity.Name}}/{{.Check.Name}}",
			Usage:     "The template for the message to be sent",
			Value:     &plugin.MessageTemplate,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "messageLimit",
			Env:       "OPSGENIE_MESSAGE_LIMIT",
			Argument:  "messageLimit",
			Shorthand: "l",
			Default:   130,
			Usage:     "The maximum length of the message field",
			Value:     &plugin.MessageLimit,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "descriptionTemplate",
			Env:       "OPSGENIE_DESCRIPTION_TEMPLATE",
			Argument:  "descriptionTemplate",
			Shorthand: "d",
			Default:   "{{.Check.Output}}",
			Usage:     "The template for the description to be sent",
			Value:     &plugin.DescriptionTemplate,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "descriptionLimit",
			Env:       "OPSGENIE_DESCRIPTION_LIMIT",
			Argument:  "descriptionLimit",
			Shorthand: "L",
			Default:   15000,
			Usage:     "The maximum length of the description field",
			Value:     &plugin.DescriptionLimit,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "includeEventInNote",
			Env:       "",
			Argument:  "includeEventInNote",
			Shorthand: "i",
			Default:   false,
			Usage:     "Include the event JSON in the payload sent to OpsGenie",
			Value:     &plugin.IncludeEventInNote,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "priority",
			Env:       "OPSGENIE_PRIORITY",
			Argument:  "priority",
			Shorthand: "p",
			Default:   "P3",
			Usage:     "The OpsGenie Alert Priority, use default from OPSGENIE_PRIORITY env var",
			Value:     &plugin.Priority,
		},
		&sensu.SlicePluginConfigOption[string]{
			Path:      "actions",
			Env:       "",
			Argument:  "actions",
			Shorthand: "A",
			Default:   []string{},
			Usage:     "The OpsGenie custom actions to assign to the event",
			Value:     &plugin.Actions,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "withAnnotations",
			Env:       "",
			Argument:  "withAnnotations",
			Shorthand: "w",
			Default:   false,
			Usage:     "Include the event.metadata.Annotations in details to send to OpsGenie",
			Value:     &plugin.WithAnnotations,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "withLabels",
			Env:       "",
			Argument:  "withLabels",
			Shorthand: "W",
			Default:   false,
			Usage:     "Include the event.metadata.Labels in details to send to OpsGenie",
			Value:     &plugin.WithLabels,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "fullDetails",
			Env:       "",
			Argument:  "fullDetails",
			Shorthand: "F",
			Default:   false,
			Usage:     "Include the more details to send to OpsGenie like proxy_entity_name, occurrences and agent details arch and os",
			Value:     &plugin.FullDetails,
		},
		&sensu.SlicePluginConfigOption[string]{
			Path:      "tagTemplate",
			Env:       "",
			Argument:  "tagTemplate",
			Shorthand: "T",
			Default:   []string{"{{.Entity.Name}}", "{{.Check.Name}}", "{{.Entity.Namespace}}", "{{.Entity.EntityClass}}"},
			Usage:     "The template to assign for the incident in OpsGenie",
			Value:     &plugin.TagsTemplates,
		},
	}
)

func main() {
	handler := sensu.NewGoHandler(&plugin.PluginConfig, options, checkArgs, executeHandler)
	handler.Execute()
}

func checkArgs(_ *types.Event) error {
	if len(plugin.AuthToken) == 0 {
		return fmt.Errorf("authentication token is empty")
	}
	if len(plugin.Team) == 0 {
		return fmt.Errorf("team is empty")
	}
	return nil
}

// parseEventKeyTags func returns string, string, and []string with event data
// fist string contains custom template string to use in message
// second string contains Entity.Name/Check.Name to use in alias
// []string contains Entity.Name Check.Name Entity.Namespace, event.Entity.EntityClass to use as tags in Opsgenie
func parseEventKeyTags(event *types.Event) (title string, alias string, tags []string) {
	alias = fmt.Sprintf("%s/%s", event.Entity.Name, event.Check.Name)
	title, err := templates.EvalTemplate("title", plugin.MessageTemplate, event)
	if err != nil {
		return "", "", []string{}
	}
	for _, v := range plugin.TagsTemplates {
		tag, err := templates.EvalTemplate("tags", v, event)
		if err != nil {
			return "", "", []string{}
		}
		tags = append(tags, tag)
	}
	return trim(title, plugin.MessageLimit), alias, tags
}

// parseDescription func returns string with custom template string to use in description
func parseDescription(event *types.Event) (description string) {
	description, err := templates.EvalTemplate("description", plugin.DescriptionTemplate, event)
	if err != nil {
		return ""
	}
	// allow newlines to get expanded
	description = strings.Replace(description, `\n`, "\n", -1)
	return trim(description, plugin.DescriptionLimit)
}

// parseDetails func returns a map of string string with check information for the details field
func parseDetails(event *types.Event) map[string]string {
	details := make(map[string]string)
	details["output"] = event.Check.Output
	details["command"] = event.Check.Command
	details["proxy_entity_name"] = event.Check.ProxyEntityName
	details["state"] = event.Check.State
	details["status"] = fmt.Sprintf("%d", event.Check.Status)
	details["occurrences"] = fmt.Sprintf("%d", event.Check.Occurrences)
	details["occurrences_watermark"] = fmt.Sprintf("%d", event.Check.OccurrencesWatermark)
	if plugin.FullDetails {
		details["ttl"] = fmt.Sprintf("%d", event.Check.Ttl)
		details["interval"] = fmt.Sprintf("%d", event.Check.Interval)
		details["subscriptions"] = fmt.Sprintf("%v", event.Check.Subscriptions)
		details["handlers"] = fmt.Sprintf("%v", event.Check.Handlers)

		if event.Entity.EntityClass == "agent" {
			details["arch"] = event.Entity.System.GetArch()
			details["os"] = event.Entity.System.GetOS()
			details["platform"] = event.Entity.System.GetPlatform()
			details["platform_family"] = event.Entity.System.GetPlatformFamily()
			details["platform_version"] = event.Entity.System.GetPlatformVersion()
		}
	}

	if plugin.WithAnnotations {
		if event.Check.Annotations != nil {
			for key, value := range event.Check.Annotations {
				if !strings.Contains(key, plugin.PluginConfig.Keyspace) {
					checkKey := fmt.Sprintf("%s_annotation_%s", "check", key)
					details[checkKey] = value
				}
			}
		}
		if event.Entity.Annotations != nil {
			for key, value := range event.Entity.Annotations {
				if !strings.Contains(key, plugin.PluginConfig.Keyspace) {
					entityKey := fmt.Sprintf("%s_annotation_%s", "entity", key)
					details[entityKey] = value
				}
			}
		}
	}

	if plugin.WithLabels {
		if event.Check.Labels != nil {
			for key, value := range event.Check.Labels {
				checkKey := fmt.Sprintf("%s_label_%s", "check", key)
				details[checkKey] = value
			}
		}
		if event.Entity.Labels != nil {
			for key, value := range event.Entity.Labels {
				entityKey := fmt.Sprintf("%s_label_%s", "entity", key)
				details[entityKey] = value
			}
		}
	}

	if len(plugin.SensuDashboard) > 0 {
		details["sensuDashboard"] = fmt.Sprintf("source: %s/%s/events/%s/%s \n", plugin.SensuDashboard, event.Entity.Namespace, event.Entity.Name, event.Check.Name)
	}
	return details
}

// eventPriority func read priority in the event and return alert.PX
// check.Annotations override Entity.Annotations
func eventPriority() alert.Priority {
	switch plugin.Priority {
	case "P5":
		return alert.P5

	case "P4":
		return alert.P4

	case "P3":
		return alert.P3

	case "P2":
		return alert.P2

	case "P1":
		return alert.P1

	default:
		return alert.P3
	}
}

// stringInSlice checks if a slice contains a specific string
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// switchOpsgenieRegion func
func switchOpsgenieRegion() client.ApiUrl {
	var region client.ApiUrl
	apiRegionLowCase := strings.ToLower(plugin.APIRegion)
	switch apiRegionLowCase {
	case "eu":
		region = client.API_URL_EU
	case "us":
		region = client.API_URL
	default:
		region = client.API_URL
	}
	return region
}

func executeHandler(event *types.Event) error {
	alertClient, err := alert.NewClient(&client.Config{
		ApiKey:         plugin.AuthToken,
		OpsGenieAPIURL: switchOpsgenieRegion(),
	})
	if err != nil {
		return fmt.Errorf("failed to create opsgenie client: %s", err)
	}

	if event.Check.Status != 0 {
		return createIncident(alertClient, event)
	}

	// check if event has a alert
	hasAlert, _ := getAlert(alertClient, event)

	// close incident if status == 0
	if hasAlert != notFound && event.Check.Status == 0 {
		return closeAlert(alertClient, event, hasAlert)
	}

	return nil
}

// createIncident func create an alert in OpsGenie
func createIncident(alertClient *alert.Client, event *types.Event) error {
	var (
		note string
		err  error
	)

	if plugin.IncludeEventInNote {
		note, err = getNote(event)
		if err != nil {
			return err
		}
	}

	teams := []alert.Responder{
		{Type: alert.EscalationResponder, Name: plugin.Team},
		{Type: alert.ScheduleResponder, Name: plugin.Team},
	}
	title, alias, tags := parseEventKeyTags(event)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	createResult, err := alertClient.Create(ctx, &alert.CreateAlertRequest{
		Message:     title,
		Alias:       alias,
		Description: parseDescription(event),
		Responders:  teams,
		Actions:     plugin.Actions,
		Tags:        tags,
		Details:     parseDetails(event),
		Entity:      event.Entity.Name,
		Source:      source,
		Priority:    eventPriority(),
		Note:        note,
	})

	if err != nil {
		fmt.Println(err.Error())
	} else {
		// FUTURE: send to AH
		fmt.Println("Create request ID: " + createResult.RequestId)
	}
	return nil
}

// getAlert func get a alert using an alias.
func getAlert(alertClient *alert.Client, event *types.Event) (string, error) {
	_, alias, _ := parseEventKeyTags(event)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	getResult, err := alertClient.Get(ctx, &alert.GetAlertRequest{
		IdentifierType:  alert.ALIAS,
		IdentifierValue: alias,
	})

	if err != nil {
		return notFound, nil
	}
	fmt.Printf("ID: %s, Message: %s, Count: %d \n", getResult.Id, getResult.Message, getResult.Count)
	return getResult.Id, nil
}

// closeAlert func close an alert if status == 0
func closeAlert(alertClient *alert.Client, event *types.Event, alertid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	closeResult, err := alertClient.Close(ctx, &alert.CloseAlertRequest{
		IdentifierType:  alert.ALERTID,
		IdentifierValue: alertid,
		Source:          source,
		Note:            "Closed Automatically",
	})

	if err != nil {
		fmt.Printf("[ERROR] Not Closed: %s\n", err)
	}
	// FUTURE: send to AH
	fmt.Printf("RequestID %s to Close %s\n", alertid, closeResult.RequestId)

	return nil
}

// getNote func creates a note with whole event in json format
func getNote(event *types.Event) (string, error) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Event data update:\n\n%s", eventJSON), nil
}

// time func returns only the first n bytes of a string
func trim(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
