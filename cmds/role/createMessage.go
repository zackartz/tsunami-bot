package role

import (
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/prisma/prisma-client-go/runtime/types"
	"github.com/zackartz/cmdlr2"
	"github.com/zackartz/tsunami-bot/db"
)

var CreateMessageRoleCommand = &cmdlr2.Command{
	Name:        "create",
	Description: "Create a new role selection message.",
	Usage:       "create",
	Example:     "create",
	Handler: func(ctx *cmdlr2.Ctx) {
		embed := &disgord.Embed{
			Title:       "Roles",
			Description: "Use the following to pick roles!",
			Fields:      []*disgord.EmbedField{},
		}

		if ctx.Event.Message.Author.ID == 133314498214756352 || ctx.Event.Message.Author.ID == 271787171889807360 {
			m, _ := ctx.Client.Channel(ctx.Event.Message.ChannelID).CreateMessage(&disgord.CreateMessageParams{Embed: embed})

			_, err := db.Client.RoleMessage.CreateOne(
				db.RoleMessage.ChannelID.Set(types.BigInt(ctx.Event.Message.ChannelID)),
				db.RoleMessage.GuildID.Set(types.BigInt(ctx.Event.Message.GuildID)),
				db.RoleMessage.MessageID.Set(types.BigInt(m.ID)),
			).Exec(context.Background())
			if err != nil {
				_ = ctx.Client.Channel(ctx.Event.Message.ChannelID).DeleteMessages(&disgord.DeleteMessagesParams{Messages: []disgord.Snowflake{m.ID}})
				ctx.ResponseText(fmt.Sprintf("%v", err))
			}
		}
	},
}
