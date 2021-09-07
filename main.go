package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/zackartz/cmdlr2"
	"github.com/zackartz/tsunami-bot/cmds"
	"github.com/zackartz/tsunami-bot/db"
	"github.com/zackartz/tsunami-bot/events"
)

var (
	TOKEN string
	log   zerolog.Logger
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Err(err).Msg("Could not read environment")
	}

	go db.Init()
	if err != nil {
		log.Err(err).Msg("Could not connect to db")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	log = zerolog.New(output).With().Timestamp().Logger()

	TOKEN = os.Getenv("TOKEN")
	if TOKEN == "" {
		log.Info().Msg("no token")
	}

	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
	})
	defer client.Gateway().StayConnectedUntilInterrupted()

	client.Gateway().MessageReactionAdd(events.EmojiAdd)
	client.Gateway().MessageReactionRemove(events.EmojiRemove)

	router := cmdlr2.Create(&cmdlr2.Router{
		Prefixes:         []string{"&"},
		Client:           client,
		BotsAllowed:      false,
		IgnorePrefixCase: true,
	})

	router.RegisterCMDList(cmds.CommandList)

	router.RegisterDefaultHelpCommand(client)

	router.Initialize(client)
}
