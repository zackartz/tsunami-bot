package role

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prisma/prisma-client-go/runtime/types"
	"github.com/zackartz/tsunami-bot/db"

	"github.com/andersfylling/disgord"
	"github.com/zackartz/cmdlr2"
)

var AddRoleCommand = &cmdlr2.Command{
	Name:        "add",
	Description: "add a role to the commandlist",
	Usage:       "add @role :emoji:",
	Example:     "add @DPS :DPS:",
	Handler: func(ctx *cmdlr2.Ctx) {
		if ctx.Event.Message.Author.ID == 133314498214756352 || ctx.Event.Message.Author.ID == 261558620632645633 {
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

			if err != nil {
				ctx.ResponseText("Try making a message")
				return
			}

			var emoji *disgord.Emoji

			if strings.HasPrefix(ctx.Args.Get(1).Raw(), "<") {
				regex := regexp.MustCompile("\\d+")
				id := regex.Find([]byte(ctx.Args.Get(1).Raw()))
				eID, err := strconv.ParseInt(string(id), 10, 64)
				if err != nil {
					ctx.ResponseText(fmt.Sprintf("%v", err))
				}
				emoji, err = ctx.Client.Guild(ctx.Event.Message.GuildID).Emoji(disgord.Snowflake(eID)).Get()
				if err != nil {
					ctx.ResponseText(fmt.Sprintf("%v", err))
					return
				}
				_, err = db.Client.Role.CreateOne(
					db.Role.ID.Set(strconv.Itoa(int(role.ID))),
					db.Role.Name.Set(role.Name),
					db.Role.Emoji.Set(emoji.Name),
					db.Role.EmojiID.Set(types.BigInt(eID)),
					db.Role.RoleMessage.Link(
						db.RoleMessage.ChannelID.Equals(types.BigInt(ctx.Event.Message.ID)),
					),
				).Exec(context.Background())
				if err != nil {
					ctx.ResponseText(fmt.Sprintf("%v", err))
					return
				}
			} else {
				emoji = &disgord.Emoji{Name: ctx.Args.Get(1).Raw()}
				_, err = db.Client.Role.CreateOne(
					db.Role.ID.Set(strconv.Itoa(int(role.ID))),
					db.Role.Name.Set(role.Name),
					db.Role.Emoji.Set(ctx.Args.Get(1).Raw()),
					db.Role.EmojiID.Set(types.BigInt(0)),
					db.Role.RoleMessage.Link(
						db.RoleMessage.ChannelID.Equals(types.BigInt(ctx.Event.Message.ID)),
					),
				).Exec(context.Background())
				if err != nil {
					ctx.ResponseText(fmt.Sprintf("%v", err))
					return
				}
			}

			msg, err := db.Client.RoleMessage.FindFirst(
				db.RoleMessage.ChannelID.Equals(types.BigInt(ctx.Event.Message.ChannelID)),
			).Exec(context.Background())

			embed := renderEmbed(msg)

			_, err = ctx.Client.Channel(ctx.Event.Message.ChannelID).Message(disgord.Snowflake(msg.ChannelID)).UpdateBuilder().SetEmbed(embed).Execute()

			err = ctx.Client.Channel(disgord.Snowflake(msg.ChannelID)).Message(disgord.Snowflake(msg.ChannelID)).Reaction(emoji).Create()
			if err != nil {
				ctx.ResponseText(fmt.Sprintf("%v", err))
			}

			_ = ctx.Client.Channel(ctx.Event.Message.ChannelID).Message(ctx.Event.Message.ID).Delete()
		}
	},
}

func getRoles(roles []*disgord.Role, rID int64) (*disgord.Role, error) {
	for _, r := range roles {
		if r.ID == disgord.Snowflake(rID) {
			return r, nil
		}
	}
	return nil, errors.New("couldn't find role")
}

func renderEmbed(msg *db.RoleMessageModel) *disgord.Embed {
	var fields []*disgord.EmbedField
	roles := msg.Roles()

	for _, x := range roles {
		var val string

		if x.EmojiID != 0 {
			val = fmt.Sprintf("<:%s:%d>", x.Emoji, x.EmojiID)
		} else {
			val = x.Emoji
		}

		fields = append(fields, &disgord.EmbedField{
			Name:   x.Name,
			Value:  val,
			Inline: true,
		})
	}

	return &disgord.Embed{
		Title:       "Roles",
		Type:        "rich",
		Description: "Use the following to pick roles!",
		Timestamp: disgord.Time{
			Time: time.Now(),
		},
		Color:  0xffff00,
		Fields: fields,
	}
}
