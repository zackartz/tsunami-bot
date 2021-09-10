package events

import (
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/prisma/prisma-client-go/runtime/types"
	"github.com/zackartz/tsunami-bot/db"
)

func EmojiRemove(s disgord.Session, h *disgord.MessageReactionRemove) {
	r, err := db.Client.Role.FindFirst(
		db.Role.EmojiID.Equals(types.BigInt(h.PartialEmoji.ID)),
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

	err = s.Guild(disgord.Snowflake(rm.GuildID)).Member(h.UserID).RemoveRole(disgord.ParseSnowflakeString(r.ID))
	if err != nil {
		s.Channel(h.ChannelID).CreateMessage(&disgord.CreateMessageParams{
			Content: fmt.Sprintf("%v", err),
		})
	}
}
