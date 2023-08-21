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
	invalidPromptResponse       = "invalid prompt, please use the following format: apollo prompt: \"<prompt>\""
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
			Address string `env:"API_ADDRESS" envDefault:"http://localhost:80"`
		}
	}
	if err := env.Parse(&cfg); err != nil {
		return errors.Wrap(err, "parsing configuration")
	}

	inferenceClient := inference.NewClient("inference", "http://localhost:8082")

	// TODO: move all this discord bot logic to a separate package
	sess, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		hasPrompt := false
		sections := strings.Split(m.Content, "\"")
		args := strings.Split(sections[0], " ")
		var prompt string
		if len(sections) >= 2 {
			prompt = sections[1]
			fmt.Println("prompt is " + prompt) // delete
			if prompt != "" {
				hasPrompt = true
			}
		}
		if args[0] != apolloCommandPrefix {
			return
		}

		if args[1] == statusArg {
			// TODO: execute healthcheck on inference service
			fmt.Println("Healthcheck started") // delete
			inferenceClient.HealthCheck(ctx)
			s.ChannelMessageSend(m.ChannelID, apolloStatusOnlineResponse)
		} else if args[1] == promptArg {
			if hasPrompt {
				promptResponse := ""
				fmt.Println("Prompt request started") // delete
				resp, err := inferenceClient.Prompt(ctx, inference.PromptRequest{Prompt: prompt})
				if err != nil {
					fmt.Printf("error reuesting prompt: %v", err)
				}
				promptResponse = resp.Response
				s.ChannelMessageSend(m.ChannelID, outputHeading+promptResponse)
			} else {
				if len(args) > 2 {
					s.ChannelMessageSend(m.ChannelID, invalidPromptResponse)
				} else {
					s.ChannelMessageSend(m.ChannelID, emptyPromptResponse)
				}
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

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println(apolloStatusOnlineResponse)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return nil
}
