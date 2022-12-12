package plugin

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-plugin-api/cluster"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

	"github.com/mattermost/mattermost-plugin-servicenow-virtual-agent/server/serializer"
)

func Test_MessageHasBeenPosted(t *testing.T) {
	defer monkey.UnpatchAll()

	for _, testCase := range []struct {
		description                       string
		Message                           string
		getChannelError                   *model.AppError
		getDirectChannelError             *model.AppError
		getUserError                      error
		parseAuthTokenError               error
		sendMessageToVirtualAgentAPIError error
		createMessageAttachmentError      error
		scheduleJobErr                    error
	}{
		{
			description: "Message is posted and successfully sent to Virtual Agent",
			Message:     "mockMessage",
		},
		{
			description:     "Message is posted but failed to get current channel",
			getChannelError: &model.AppError{},
			Message:         "mockMessage",
		},
		{
			description:  "Message is posted but failed to get user from KV store",
			getUserError: errors.New("error getting the user from KVstore"),
			Message:      "mockMessage",
		},
		{
			description:  "Message is posted but user is not connected to ServiceNow",
			getUserError: ErrNotFound,
			Message:      "mockMessage",
		},
		{
			description: "Message is posted but user is not connected to ServiceNow",
			Message:     "disconnect",
		},
		{
			description:         "Message is posted but failed to parse auth token",
			parseAuthTokenError: errors.New("error in parsing the auth token"),
			Message:             "mockMessage",
		},
		{
			description:                       "Message is posted but failed to parse auth token",
			sendMessageToVirtualAgentAPIError: errors.New("error in parsing the auth token"),
			Message:                           "mockMessage",
		},
		{
			description:                  "Message is posted but failed to create message attachment",
			createMessageAttachmentError: errors.New("error in creating message attachment"),
			Message:                      "mockMessage",
		},
		{
			description:           "Failed to get direct channel",
			getDirectChannelError: &model.AppError{},
			Message:               "mockMessage",
		},
		{
			description:    "Failed to schedule the job",
			scheduleJobErr: errors.New("error while scheduling the job"),
			Message:        "mockMessage",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			p := Plugin{}
			mockInterval := int64(1000)
			p.botUserID = "mock-botID"

			mockAPI := &plugintest.API{}
			mockAPI.On("LogError", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			mockAPI.On("GetChannel", "mockChannelID").Return(&model.Channel{
				Type: "D",
				Name: "mock-botID__mock",
			}, testCase.getChannelError)

			mockAPI.On("GetConfig").Return(&model.Config{
				ServiceSettings: model.ServiceSettings{
					TimeBetweenUserTypingUpdatesMilliseconds: &mockInterval,
				},
			})
			mockAPI.On("GetDirectChannel", mock.Anything, mock.Anything).Return(&model.Channel{
				Id: "mock-channelID",
			}, testCase.getDirectChannelError)

			monkey.Patch(cluster.Schedule, func(_ cluster.JobPluginAPI, _ string, _ cluster.NextWaitInterval, _ func()) (*cluster.Job, error) {
				return &cluster.Job{}, testCase.scheduleJobErr
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "Ephemeral", func(_ *Plugin, _, _, _ string, _ ...interface{}) {})

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "GetUser", func(_ *Plugin, _ string) (*serializer.User, error) {
				return &serializer.User{}, testCase.getUserError
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "DM", func(_ *Plugin, _, _ string, _ ...interface{}) (string, error) {
				return "mockPostID", nil
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "DMWithAttachments", func(_ *Plugin, _ string, _ ...*model.SlackAttachment) (string, error) {
				return "mockPostID", nil
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "ParseAuthToken", func(_ *Plugin, _ string) (*oauth2.Token, error) {
				return &oauth2.Token{}, testCase.parseAuthTokenError
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "MakeClient", func(_ *Plugin, _ context.Context, _ *oauth2.Token) Client {
				return &client{}
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(&p), "CreateMessageAttachment", func(_ *Plugin, _ string) (*MessageAttachment, error) {
				return &MessageAttachment{}, testCase.createMessageAttachmentError
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(&client{}), "SendMessageToVirtualAgentAPI", func(_ *client, _, _ string, _ bool, _ *MessageAttachment) error {
				return testCase.sendMessageToVirtualAgentAPIError
			})

			p.SetAPI(mockAPI)

			post := &model.Post{
				ChannelId: "mockChannelID",
				UserId:    "mock-userID",
				Message:   testCase.Message,
				FileIds:   []string{"mockFileID"},
			}

			p.MessageHasBeenPosted(&plugin.Context{}, post)

			mockAPI.AssertNumberOfCalls(t, "GetChannel", 1)

			if testCase.getChannelError != nil || testCase.parseAuthTokenError != nil || testCase.sendMessageToVirtualAgentAPIError != nil || (testCase.getUserError != nil && testCase.getUserError != ErrNotFound) {
				mockAPI.AssertNumberOfCalls(t, "LogError", 1)
			}
		})
	}
}
