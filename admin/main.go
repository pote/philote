package main

import(
	"code.google.com/p/go-uuid/uuid"
	"flag"
	"fmt"
	"os"
	"strings"
)


func printUsageAndExit() {
	fmt.Println(`Usage: philote-cli [create-key|publish] <flags>

Try philote-cli <command> --help for details.`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsageAndExit()
	}

	switch os.Args[1] {
	case "create-key":
		kflag := flag.NewFlagSet("create-key", flag.ExitOnError)
		token := kflag.String("token", uuid.New(), "authorization token")
		readableChannels := kflag.String("read", "test-channel", "comma-separated list of readable channels.")
		writeableChannels := kflag.String("write", "test-channel", "comma-separated list of readable channels.")
		allowedUses := kflag.Int("allowed-uses", 0, "allowed uses for token (use 0 for unlimited usage).")
		kflag.Parse(os.Args[2:])

		err := createAccessKey(
			*token, strings.Split(*readableChannels, ","),
			strings.Split(*writeableChannels, ","),
			*allowedUses)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(*token)

	case "publish":
		pflag := flag.NewFlagSet("publish", flag.ExitOnError)
		channel := pflag.String("channel", "test-channel", "channel in which to publish message")
		data := pflag.String("data", "Hello from philote-cli!", "message payload")
		pflag.Parse(os.Args[2:])

		listeners, err := publishMessage(*channel, *data); if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(listeners)

	default:
		printUsageAndExit()
	}
}
