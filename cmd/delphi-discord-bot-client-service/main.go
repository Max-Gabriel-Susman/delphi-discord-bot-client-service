package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	apolloCommandPrefix = "apollo"
	helpArg             = "help"
	statusArg           = "status"
	promptArg           = "prompt:"
	appArg              = "app"
	modelArg            = "model"
	conversationArg     = "conversation"
	indexArg            = "index"
	memoryArg           = "memory"
	chainArg            = "chain"
	agentArg            = "agent"

	apolloCommandDescription          = "Apollo is a bot for uploading, managing, and operating LLM driven applications on the delphi-inferential-cluster.\n\n"
	usageHeading                      = "Usage:\n\n"
	apolloHelpSubCommandUsageHeading  = "Use \"apollo help <command>\" for more information about a command.\n\n"
	additionalHelpTopicsHeading       = "Additional help topics:\n\n"
	outputHeading                     = "output: "
	availableCommandsHeading          = "The commands are:\n\n"
	availableHelpSubCommandsHeading   = "The commands are:\n\n"
	availableHelpTopicsHeading        = "Use \"apollo help <topic>\" for more information about that topic.\n"
	apolloCommandUsage                = "\t\tapollo <command> [arguments]\n\n"
	promptSubCommandDescription       = "\t\tprompt: submits a prompt to the delphi-inferential-cluster\n"
	conversationSubCommandDescription = ""
	modelSubCommandDescription        = ""
	appSubCommandDescription          = "\t\tapp: application management for instances hosted on the delphi-inferential-cluster\n\n"
	indexSubCommandDescription        = "\t\t\tindex: index management for instances hosted on the delphi-inferential-cluster\n"
	memorySubCommandDescription       = "\t\t\tmemory: memory management for instances hosted on the delphi-inferential-cluster\n"
	chainSubCommandDescription        = "\t\t\tchain: chain management for instances hosted on the delphi-inferential-cluster\n"
	agentSubCommandDescription        = "\t\t\tagent: agent management for instances hosted on the delphi-inferential-cluster\n"
	helpSubCommandDescription         = "\t\thelp: prints this message\n"
	statusSubCommandDescription       = "\t\tstatus: prints the status of the bot\n"

	apolloStatusOnlineResponse  = "Apollo online"
	unknownSubCommandResponse   = "unknown sub command"
	unknownHelpArgumentResponse = "unknown help argument"
	tooManyArgumentsResponse    = "too many arguments"
	emptyPromptResponse         = "empty prompts will not be submitted to the model"
	commandUnvaliableResponse   = " command currently unavailable"
)

func main() {
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
			s.ChannelMessageSend(m.ChannelID, apolloStatusOnlineResponse)
		} else if args[1] == promptArg {
			if hasPrompt {
				s.ChannelMessageSend(m.ChannelID, outputHeading+passPrompt(sections[1]))
			} else {
				s.ChannelMessageSend(m.ChannelID, emptyPromptResponse)
			}
		} else if args[1] == appArg {
			s.ChannelMessageSend(m.ChannelID, appArg+commandUnvaliableResponse)
		} else if args[1] == modelArg {
			s.ChannelMessageSend(m.ChannelID, modelArg+commandUnvaliableResponse)
		} else if args[1] == helpArg {
			if len(args) == 2 {
				helpMessage :=
					apolloCommandDescription +
						usageHeading +
						apolloCommandUsage +
						apolloHelpSubCommandUsageHeading +
						promptSubCommandDescription +
						conversationSubCommandDescription +
						appSubCommandDescription +
						modelSubCommandDescription +
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
				case conversationArg:
					s.ChannelMessageSend(m.ChannelID, conversationSubCommandDescription)
				case appArg:
					s.ChannelMessageSend(m.ChannelID, appSubCommandDescription)
				case modelArg:
					s.ChannelMessageSend(m.ChannelID, modelSubCommandDescription)
				case indexArg:
					s.ChannelMessageSend(m.ChannelID, modelSubCommandDescription)
				case memoryArg:
					s.ChannelMessageSend(m.ChannelID, modelSubCommandDescription)
				case chainArg:
					s.ChannelMessageSend(m.ChannelID, modelSubCommandDescription)
				case agentArg:
					s.ChannelMessageSend(m.ChannelID, modelSubCommandDescription)
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
}

func passPrompt(prompt string) string {
	// make a network call to the appropriate inference pipeline

	// return output
	return prompt
}
