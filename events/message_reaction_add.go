package events

import (
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/zackartz/tsunami-bot/db"
)

func EmojiAdd(s disgord.Session, h *disgord.MessageReactionAdd) {
	r, err := db.Client.Role.FindFirst(
		db.Role.Emoji.Equals(h.PartialEmoji.Name),
	).With(
		db.Role.RoleMessage.Fetch(),
	).Exec(context.Background())
	if err != nil {
		return
	}

	u, err := s.CurrentUser().Get()
	if err != nil {
		return
	}

	if h.UserID == u.ID {
		return
	}

	rm, _ := r.RoleMessage()

	if rm == nil {
		return
	}

	err = s.Guild(disgord.Snowflake(rm.GuildID)).Member(h.UserID).AddRole(disgord.ParseSnowflakeString(r.ID))
	if err != nil {
		s.Channel(h.ChannelID).CreateMessage(&disgord.CreateMessageParams{
			Content: fmt.Sprintf("%v", err),
		})
	}
}
