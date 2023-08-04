package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"gitlab.com/lcook/ports-webhook/internal/config"
	"gitlab.com/lcook/ports-webhook/internal/port"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

func main() {
	var Config string

	flag.StringVar(&Config, "c", "config.yaml", "Configuration file")
	flag.Parse()
	configFmt := "(\"" + Config + "\")"

	configFile, err := config.LoadConfig(Config)
	if err != nil {
		log.Fatalln("Error loading configuration file", err)
	}

	log.Println("Read configuration file", configFmt)
	if len(configFile.Whitelist) < 1 {
		log.Fatalln("No whitelist rules to be loaded")
	}

	log.Printf("Loaded %d port whitelist rule(s)\n", len(configFile.Whitelist))

	hook, _ := gitlab.New(gitlab.Options.Secret(configFile.HookSecret))

	http.HandleFunc(configFile.HookPath, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, gitlab.PushEvents)
		if err != nil {
			return
		}

		switch payload.(type) {
		case gitlab.PushEventPayload:
			data := payload.(gitlab.PushEventPayload)
			if data.Ref != "refs/heads/master" {
				break
			}

			portHook := port.NewHook(configFile, data)
			portHook.Commits()
			portHook.Sync()
		}
	})

	http.ListenAndServe(":"+strconv.Itoa(configFile.HookPort), nil)
}
