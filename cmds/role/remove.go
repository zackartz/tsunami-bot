package role

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake"
	"github.com/prisma/prisma-client-go/runtime/types"
	"github.com/zackartz/cmdlr2"
	"github.com/zackartz/tsunami-bot/db"
)

var RoleRemoveCommand = &cmdlr2.Command{
	Name:        "remove",
	Description: "Remove a role from the commandlist",
	Usage:       "remove @role",
	Example:     "remove @DPS",
	Handler: func(ctx *cmdlr2.Ctx) {
		if ctx.Event.Message.Author.ID == 133314498214756352 || ctx.Event.Message.Author.ID == 271787171889807360 {
			rID, err := strconv.ParseInt(ctx.Args.Get(0).AsRoleMentionID(), 10, 64)
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
				return
			}

			var role *disgord.Role
			roles, err := ctx.Client.Guild(ctx.Event.Message.GuildID).GetRoles()
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
				return
			}

			role, err = getRoles(roles, rID)
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
				return
			}

			msg, err := db.Client.RoleMessage.FindUnique(
				db.RoleMessage.ChannelID.Equals(types.BigInt(ctx.Event.Message.ChannelID)),
			).With(
				db.RoleMessage.Roles.Fetch(),
			).Exec(context.Background())
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
				return
			}

			var r *db.RoleModel
			for _, x := range msg.Roles() {
				if x.ID == role.ID.String() {
					r = &x
					break
				}
			}

			if r == nil {
				return
			}

			var emoji *disgord.Emoji

			if r.EmojiID != 0 {
				emojis, err := ctx.Client.Guild(ctx.Event.Message.GuildID).GetEmojis()
				for _, e := range emojis {
					if snowflake.Snowflake(e.ID) == snowflake.Snowflake(r.EmojiID) {
						emoji = e
					}
				}

				_, err = db.Client.RoleMessage.FindUnique(
					db.RoleMessage.ChannelID.Equals(types.BigInt(ctx.Event.Message.ChannelID)),
				).Update(
					db.RoleMessage.Roles.Unlink(
						db.Role.ID.Equals(role.ID.String()),
					),
				).Exec(context.Background())

				_, err = db.Client.Role.FindUnique(
					db.Role.ID.Equals(role.ID.String()),
				).Delete().Exec(context.Background())

				if err != nil {
					return
				}
			} else {
				emoji = &disgord.Emoji{Name: r.Emoji}
				_, err = db.Client.RoleMessage.FindUnique(
					db.RoleMessage.ChannelID.Equals(types.BigInt(ctx.Event.Message.ChannelID)),
				).Update(
					db.RoleMessage.Roles.Unlink(
						db.Role.ID.Equals(role.ID.String()),
					),
				).Exec(context.Background())

				_, err = db.Client.Role.FindUnique(
					db.Role.ID.Equals(role.ID.String()),
				).Delete().Exec(context.Background())

				if err != nil {
					return
				}
			}

			msg, err = db.Client.RoleMessage.FindUnique(
				db.RoleMessage.ChannelID.Equals(types.BigInt(ctx.Event.Message.ChannelID)),
			).With(
				db.RoleMessage.Roles.Fetch(),
			).Exec(context.Background())
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
				return
			}

			embed := renderEmbed(msg)

			_, err = ctx.Client.Channel(ctx.Event.Message.ChannelID).Message(disgord.Snowflake(msg.MessageID)).UpdateBuilder().SetEmbed(embed).Execute()

			err = ctx.Client.Channel(disgord.Snowflake(msg.ChannelID)).Message(disgord.Snowflake(msg.MessageID)).Reaction(emoji).DeleteOwn()
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
			}

			_ = ctx.Client.Channel(ctx.Event.Message.ChannelID).Message(ctx.Event.Message.ID).Delete()
		}
	},
}
