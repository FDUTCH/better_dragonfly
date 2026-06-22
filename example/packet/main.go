package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/FDUTCH/better_dragonfly/packet_hook"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/pelletier/go-toml"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func main() {
	// creating Hooker.
	hooker := packet_hook.NewHooker()

	// registering hooks.
	hooker.Register(packet.IDPlayerAuthInput, func(player *player.Player, handler packet_hook.PacketHandler) packet_hook.PacketHandler {
		return Hook{handler}
	})

	slog.SetLogLoggerLevel(slog.LevelDebug)
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	conf, err := readConfig(slog.Default())
	if err != nil {
		panic(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()
	for p := range srv.Accept() {
		// hooking packet handlers.
		hooker.Hook(p)
	}
}

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}

type Hook struct {
	h packet_hook.PacketHandler
}

func (h Hook) Handle(p packet.Packet, s *session.Session, tx *world.Tx, c session.Controllable) error {
	fmt.Printf("handled packet %T before original handler\n", p)
	return h.h.Handle(p, s, tx, c)
}
