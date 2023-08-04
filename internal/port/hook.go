package port

import (
	"fmt"
	"os/exec"
	"strings"

	"gitlab.com/lcook/ports-webhook/internal/config"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

type Hook struct {
	Ports   []Port
	Config  config.Config
	Payload gitlab.PushEventPayload
}

func NewHook(cfg config.Config, payload gitlab.PushEventPayload) *Hook {
	return &Hook{[]Port{}, cfg, payload}
}

func (h *Hook) Commits() {
	commits := h.Payload.Commits
	for _, c := range commits {
		fmt.Printf("\nCommit: %s\nAuthor: %s <%s>\n\n", c.ID, c.Author.Name, c.Author.Email)
		for _, a := range c.Added {
			fmt.Printf("A\t%s\n", a)
			h.Add(a)
		}
		for _, m := range c.Modified {
			fmt.Printf("M\t%s\n", m)
			h.Add(m)
		}
		for _, r := range c.Removed {
			fmt.Printf("R\t%s\n", r)
			h.Add(r)
		}
		fmt.Println()
	}

	h.Dedupe()
}

func (h *Hook) Add(val string) {
	var port Port
	for i, p := range h.Config.Whitelist {
		c := strings.Split(p, "/")[0]
		p = strings.Split(p, "/")[1]
		port = Port{c, p}
		if strings.Contains(val, port.Fullname()) {
			break
		}
		if i == len(h.Config.Whitelist)-1 {
			return
		}
	}

	h.Ports = append(h.Ports, port)
}

func (h *Hook) Dedupe() {
	keys := make(map[Port]bool)
	var unique []Port
	for _, port := range h.Ports {
		if _, value := keys[port]; !value {
			keys[port] = true
			unique = append(unique, port)
		}
	}

	h.Ports = unique
}

func (h *Hook) Sync() {
	for _, port := range h.Ports {
		exec.Command("glport", "-u", "-p", port.Name, "-c", port.Category, "-m", "Sync with upstream", "-o").Run()
		fmt.Printf("%s: Synced with upstream repository\n", port.Fullname())
	}
}
