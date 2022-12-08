package plugin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type FileStruct struct {
	ID     string
	Expiry time.Time
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	// If the message is posted by bot simply return
	if post.UserId == p.botUserID {
		return
	}

	if presentInCache := p.dmChannelCache.Has(post.ChannelId); !presentInCache {
		channel, appErr := p.API.GetChannel(post.ChannelId)
		if appErr != nil {
			p.API.LogError("Error occurred while fetching channel by ID. ChannelID: %s. Error: %s", post.ChannelId, appErr.Error())
			return
		}

		if channel.Type != model.CHANNEL_DIRECT {
			return
		}

		channelName := strings.Split(channel.Name, "__")
		if channelName[0] != p.botUserID && channelName[1] != p.botUserID {
			return
		}

		if err := p.dmChannelCache.SetWithExpire(post.ChannelId, true, time.Minute*time.Duration(p.getConfiguration().ChannelCacheTTL)); err != nil {
			p.API.LogDebug("Failed to add channel in cache", "Error", err.Error())
		}
	}

	mattermostUserID := post.UserId
	// Check if the user is connected to ServiceNow
	user, err := p.GetUser(mattermostUserID)
	if err != nil {
		if err == ErrNotFound {
			_, _ = p.DM(mattermostUserID, WelcomePretextMessage, fmt.Sprintf("%s%s", p.GetPluginURL(), PathOAuth2Connect))
		} else {
			p.logAndSendErrorToUser(mattermostUserID, post.ChannelId, fmt.Sprintf("Error occurred while fetching user by ID. UserID: %s. Error: %s", mattermostUserID, err.Error()))
		}
		return
	}

	if strings.ToLower(post.Message) == DisconnectKeyword {
		_, _ = p.DMWithAttachments(post.UserId, p.CreateDisconnectUserAttachment())
		return
	}

	if len(post.FileIds) > 1 {
		p.logAndSendErrorToUser(mattermostUserID, post.ChannelId, "Cannot send more than one file attachment at a time.")
		return
	}

	token, err := p.ParseAuthToken(user.OAuth2Token)
	if err != nil {
		p.logAndSendErrorToUser(mattermostUserID, post.ChannelId, fmt.Sprintf("Error occurred while decrypting token. Error: %s", err.Error()))
		return
	}

	var attachment *MessageAttachment
	if len(post.FileIds) == 1 {
		attachment, err = p.CreateMessageAttachment(post.FileIds[0])
		if err != nil {
			p.logAndSendErrorToUser(mattermostUserID, post.ChannelId, err.Error())
			return
		}
	}

	client := p.MakeClient(context.Background(), token)
	if err = client.SendMessageToVirtualAgentAPI(user.UserID, post.Message, true, attachment); err != nil {
		p.logAndSendErrorToUser(mattermostUserID, post.ChannelId, err.Error())
	}
}
