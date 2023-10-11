package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	pb "github.com/Max-Gabriel-Susman/delphi-discord-bot-client-service/inference"
	"github.com/Max-Gabriel-Susman/delphi-discord-bot-client-service/internal/clients/inference"
	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "172.17.0.2:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

const (
	delphiCommandPrefix = "delphi"
	helpArg             = "help"
	statusArg           = "status"
	promptArg           = "prompt:"
	startArg            = "start"
	stopArg             = "stop"
	restartArg          = "restart"

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

	delphiStartingClusterResponse   = "Starting delphi cluster"
	delphiStoppingClusterResponse   = "Stopping delphi cluster"
	delphiRestartingClusterResponse = "Restarting delphi cluster"
	delphiStatusOnlineResponse      = "delphi online"
	unknownSubCommandResponse       = "unknown sub command"
	rawPromptResponse               = "raw prompt response"
	unknownHelpArgumentResponse     = "unknown help argument"
	tooManyArgumentsResponse        = "too many arguments"
	emptyPromptResponse             = "empty prompts will not be submitted to the model"
	invalidPromptResponse           = "invalid prompt, please use the following format: delphi prompt: \"<prompt>\""
	commandUnvaliableResponse       = " command currently unavailable"
)

func main() {
	ctx := context.Background()
	run(ctx, os.Args)
}

func run(ctx context.Context, _ []string) error {
	var cfg struct {
		ServiceName           string `env:"SERVICE_NAME" envDefault:"delphi-discord-bot-client-service"`
		InfrastructureService struct {
			Host string `env:"INFRASTRUCTURE_SERVICE_HOST" envDefault:"localhost"`
			Port string `env:"INFRASTRUCTURE_SERVICE_PORT" envDefault:"8080"`
		}
		Env string `env:"ENV" envDefault:"local"`
		API struct {
			Address string `env:"API_ADDRESS" envDefault:"http://localhost:80"`
		}
	}
	if err := env.Parse(&cfg); err != nil {
		return errors.Wrap(err, "parsing configuration")
	}

	// applicationServiceClient := application.NewClient("application", "http://localhost:8000")
	inferentialServiceClient := inference.NewClient("application", "http://localhost:8082")

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
		if args[0] != delphiCommandPrefix {
			conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewGreeterClient(conn)

			// Contact the server and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Greeting: %s", r.GetMessage())
			s.ChannelMessageSend(m.ChannelID, rawPromptResponse)
		} else {
			if args[1] == statusArg {
				// TODO: execute healthcheck on inference service
				fmt.Println("Healthcheck started") // delete
				inferentialServiceClient.HealthCheck(ctx)
				s.ChannelMessageSend(m.ChannelID, delphiStatusOnlineResponse)
			} else if args[1] == promptArg {
				if hasPrompt {
					promptResponse := ""
					fmt.Println("Prompt request started") // delete
					resp, err := inferentialServiceClient.Prompt(ctx, inference.PromptRequest{Prompt: prompt})
					if err != nil {
						fmt.Printf("error requesting prompt: %v", err)
					}
					promptResponse = resp.Response
					s.ChannelMessageSend(m.ChannelID, promptResponseOutputHeading+promptResponse)
				} else {
					if len(args) > 2 {
						s.ChannelMessageSend(m.ChannelID, invalidPromptResponse)
					} else {
						s.ChannelMessageSend(m.ChannelID, emptyPromptResponse)
					}
				}
			} else if args[1] == startArg {
				// check user priveleges maybe l8r idk
				fmt.Println("Initializing Delphi cluster") // delete
				s.ChannelMessageSend(m.ChannelID, delphiStartingClusterResponse)

				// request deploy from monitoring svc(needs new name) monitoring, meta, administrative, infrastructure
				// let's do ifrastructure svc

				// process reponse

			} else if args[1] == stopArg {
				fmt.Println("Stopping Delphi cluster") // delete
				s.ChannelMessageSend(m.ChannelID, delphiStoppingClusterResponse)
			} else if args[1] == restartArg {
				fmt.Println("Initializing Delphi cluster") // delete
				s.ChannelMessageSend(m.ChannelID, delphiRestartingClusterResponse)
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
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println(delphiStatusOnlineResponse)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return nil
}
