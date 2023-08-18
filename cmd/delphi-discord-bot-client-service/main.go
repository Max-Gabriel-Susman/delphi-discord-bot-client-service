package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Max-Gabriel-Susman/delphi-discord-bot-client-service/internal/clients/inference"
	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

const (
	apolloCommandPrefix = "apollo"
	helpArg             = "help"
	statusArg           = "status"
	promptArg           = "prompt:"

	apolloCommandDescription         = "Apollo is a bot for uploading, managing, and operating LLM driven applications on the delphi-inferential-cluster.\n\n"
	usageHeading                     = "Usage:\n\n"
	apolloHelpSubCommandUsageHeading = "Use \"apollo help <command>\" for more information about a command.\n\n"
	additionalHelpTopicsHeading      = "Additional help topics:\n\n"
	outputHeading                    = "output: "
	availableCommandsHeading         = "The commands are:\n\n"
	availableHelpSubCommandsHeading  = "The commands are:\n\n"
	availableHelpTopicsHeading       = "Use \"apollo help <topic>\" for more information about that topic.\n"
	apolloCommandUsage               = "\t\tapollo <command> [arguments]\n\n"
	promptSubCommandDescription      = "\t\tprompt: submits a prompt to the delphi-inferential-cluster\n"
	helpSubCommandDescription        = "\t\thelp: prints this message\n"
	statusSubCommandDescription      = "\t\tstatus: prints the status of the bot\n"

	apolloStatusOnlineResponse  = "Apollo online"
	unknownSubCommandResponse   = "unknown sub command"
	unknownHelpArgumentResponse = "unknown help argument"
	tooManyArgumentsResponse    = "too many arguments"
	emptyPromptResponse         = "empty prompts will not be submitted to the model"
	commandUnvaliableResponse   = " command currently unavailable"
)

func main() {
	ctx := context.Background()
	run(ctx, os.Args)
}

func run(ctx context.Context, _ []string) error {
	var cfg struct {
		ServiceName string `env:"SERVICE_NAME" envDefault:"delphi-discord-bot-client-service"`
		Env         string `env:"ENV" envDefault:"local"`
		API         struct {
			Address string `env:"API_ADDRESS" envDefault:"http://localhost:8080"`
		}
	}
	if err := env.Parse(&cfg); err != nil {
		return errors.Wrap(err, "parsing configuration")
	}

	inferenceClient := inference.NewClient("inference", "http://localhost:8080")

	fmt.Println("before discord connection") // delete
	sess, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("after discord connection") // delete

	fmt.Println("before handler additon") // delete
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		hasPrompt := false
		sections := strings.Split(m.Content, "\"")
		args := strings.Split(sections[0], " ")
		if len(sections) >= 2 {
			prompt := sections[1]
			if prompt != "" {
				hasPrompt = true
			}
		}
		if args[0] != apolloCommandPrefix {
			return
		}

		if args[1] == statusArg {
			// TODO: execute healthcheck on inference service
			inferenceClient.HealthCheck(ctx)
			s.ChannelMessageSend(m.ChannelID, apolloStatusOnlineResponse)
		} else if args[1] == promptArg {
			if hasPrompt {
				var inferenceGenerationResponse string
				// TODO: execute inference generation on inference service
				s.ChannelMessageSend(m.ChannelID, outputHeading+inferenceGenerationResponse)
			} else {
				s.ChannelMessageSend(m.ChannelID, emptyPromptResponse)
			}
		} else if args[1] == helpArg {
			if len(args) == 2 {
				helpMessage :=
					apolloCommandDescription +
						usageHeading +
						apolloCommandUsage +
						apolloHelpSubCommandUsageHeading +
						promptSubCommandDescription +
						helpSubCommandDescription +
						additionalHelpTopicsHeading +
						availableHelpTopicsHeading
				s.ChannelMessageSend(m.ChannelID, helpMessage)
			} else if len(args) == 3 {
				switch args[2] {
				case helpArg:
					s.ChannelMessageSend(m.ChannelID, apolloHelpSubCommandUsageHeading)
				case statusArg:
					s.ChannelMessageSend(m.ChannelID, statusSubCommandDescription)
				case promptArg:
					s.ChannelMessageSend(m.ChannelID, promptSubCommandDescription)
				default:
					s.ChannelMessageSend(m.ChannelID, unknownHelpArgumentResponse)
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, tooManyArgumentsResponse)
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, unknownSubCommandResponse)
		}
	})
	fmt.Println("after handler additon") // delete

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	fmt.Println("before session open") // delete
	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()
	fmt.Println("after session open") // delete

	fmt.Println(apolloStatusOnlineResponse)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return nil
}
