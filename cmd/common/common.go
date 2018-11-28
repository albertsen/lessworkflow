package common

import (
	"flag"
	"io/ioutil"
	"os"
)

func Flags() (*string, *string, *bool) {
	url := flag.String("s", "nats://localhost:4222", "URL of messaging server.")
	topic := flag.String("t", "actions", "Message topic.")
	help := flag.Bool("h", false, "This message.")
	return url, topic, help
}

func ParseFlags(ShowHelp *bool) {
	flag.Parse()
	if *ShowHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func ReadFileOrStdin(FileName string) ([]byte, error) {
	if FileName == "" {
		return ioutil.ReadAll(os.Stdin)
	} else {
		return ioutil.ReadFile(FileName)
	}
}
