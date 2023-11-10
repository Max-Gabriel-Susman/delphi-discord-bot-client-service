package user

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/Max-Gabriel-Susman/delphi-discord-bot-client-service/inference"
	"github.com/bwmarrin/discordgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: I need to figure out a good way to mock this, probably going to need to dramatically refactor

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "0.0.0.0:50054", "the address to connect to")
	// addr = flag.String("addr", "172.17.0.2:50054", "the address to connect to")
)

// delphi list available inference

const (
	delphiCommandPrefix = "delphi"

	helpArg   = "help"
	statusArg = "status"
	setArg    = "set"
	getArg    = "get"
	listArg   = "list"
	trainArg  = "train"

	availableArg = "available"
	active       = "active"

	inferenceArg = "inference"

	modelArg  = "model" // do we actually need this argument or is it just redundant?
	modelsArg = "model"

	gptArg    = "gpt"
	falconArg = "falcon"
	t5Arg     = "t5"

	delphiCommandDescription         = "delphi is a bot for uploading, managing, and operating LLM driven applications on the delphi-inferential-cluster.\n\n"
	usageHeading                     = "Usage:\n\n"
	delphiHelpSubCommandUsageHeading = "Use \"delphi help <command>\" for more information about a command.\n\n"
	additionalHelpTopicsHeading      = "Additional help topics:\n\n"
	promptResponseOutputHeading      = "Model output: "
	availableCommandsHeading         = "The commands are:\n\n"
	availableHelpSubCommandsHeading  = "The commands are:\n\n"
	availableHelpTopicsHeading       = "Use \"delphi help <topic>\" for more information about that topic.\n"
	delphiCommandUsage               = "\t\tdelphi <command> [arguments]\n\n"
	promptSubCommandDescription      = "\t\tprompt: submits a prompt to the delphi-inferential-cluster\n"
	helpSubCommandDescription        = "\t\thelp: prints this message\n"
	statusSubCommandDescription      = "\t\tstatus: prints the status of the bot\n"

	delphiStatusOnlineResponse  = "delphi online"
	unknownSubCommandResponse   = "unknown sub command"
	rawPromptResponse           = "raw prompt response"
	unknownHelpArgumentResponse = "unknown help argument"
	tooManyArgumentsResponse    = "too many arguments"
	emptyPromptResponse         = "empty prompts will not be submitted to the model"
	invalidPromptResponse       = "invalid prompt, please use the following format: delphi prompt: \"<prompt>\""
	commandUnvaliableResponse   = " command currently unavailable"
)

func InitiateDiscordBotSession(ctx context.Context) {
	// TODO: move all this discord bot logic to a separate package
	sess, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		sections := strings.Split(m.Content, "\"")
		args := strings.Split(sections[0], " ")
		if args[0] != delphiCommandPrefix {
			conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewGreeterClient(conn)

			// Contact the server and print out its response.// Set a timeout for the context.
			timeout := 60 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			r, err := c.SayHello(ctx, &pb.HelloRequest{Name: m.Content})
			if err != nil {
				// log.Fatalf("could not greet: %v", err)
				fmt.Printf("could not greet: %v", err)
			}
			log.Printf("Greeting: %s", r.GetMessage())
			s.ChannelMessageSend(m.ChannelID, r.Message)
		} else {
			if args[1] == statusArg {
				// TODO: execute healthcheck on inference service
				fmt.Println("Healthcheck started") // delete
				s.ChannelMessageSend(m.ChannelID, delphiStatusOnlineResponse)
			} else if args[1] == helpArg {
				if len(args) == 2 {
					helpMessage :=
						delphiCommandDescription +
							usageHeading +
							delphiCommandUsage +
							delphiHelpSubCommandUsageHeading +
							promptSubCommandDescription +
							helpSubCommandDescription +
							additionalHelpTopicsHeading +
							availableHelpTopicsHeading
					s.ChannelMessageSend(m.ChannelID, helpMessage)
				} else if len(args) == 3 {
					switch args[2] {
					case helpArg:
						s.ChannelMessageSend(m.ChannelID, delphiHelpSubCommandUsageHeading)
					case statusArg:
						s.ChannelMessageSend(m.ChannelID, statusSubCommandDescription)
					default:
						s.ChannelMessageSend(m.ChannelID, unknownHelpArgumentResponse)
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, tooManyArgumentsResponse)
				}
			} else if args[1] == trainArg {

			} else {
				s.ChannelMessageSend(m.ChannelID, unknownSubCommandResponse)
			}
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println(delphiStatusOnlineResponse)
}
