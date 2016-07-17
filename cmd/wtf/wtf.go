package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli"
	"pault.ag/go/wtf"
)

var API_URI string = "https://wtf.pault.ag/wtf/"

func wtfIs(word string) error {
	res, err := http.Get(fmt.Sprintf("%s/%s", API_URI, word))
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == 404 {
		fmt.Printf("%s: not found\n", word)
		return nil
	} else if res.StatusCode != 200 {
		log.Fatalf("Bad HTTP code: %d", res.StatusCode)
		return fmt.Errorf("WTF")
	}

	acronyms := wtf.Acronyms{}
	if err := json.NewDecoder(res.Body).Decode(&acronyms); err != nil {
		return err
	}

	fmt.Printf("%s:\n", word)
	for _, acronym := range acronyms {
		fmt.Printf("  %s\n", acronym.String())
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "wtf"
	app.Usage = "the fuck is that"
	app.Version = "0.0.1~alpha1"

	app.Action = func(c *cli.Context) error {
		args := c.Args()
		if len(args) > 0 && args[0] == "is" {
			args = args[1:]
		}
		for _, word := range args {
			if err := wtfIs(word); err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		}
		return nil
	}
	app.Run(os.Args)
}
