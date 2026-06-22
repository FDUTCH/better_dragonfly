package packet_hook

import (
	"reflect"

	"github.com/FDUTCH/better_dragonfly/internal/unsafe_fields"
	"github.com/df-mc/dragonfly/server/player"
)

// NewHooker ...
func NewHooker() *Hooker {
	return &Hooker{constructors: make(map[uint32]Hook)}
}

type Hook = func(player *player.Player, handler PacketHandler) PacketHandler

// Hooker ...
type Hooker struct {
	constructors map[uint32]Hook
}

// Register registers packet hook.
func (h *Hooker) Register(id uint32, hook Hook) {
	h.constructors[id] = hook
}

// Hook hooks registered packets of the player.
func (h *Hooker) Hook(pl *player.Player) {
	conf := pl.Data()
	handlers := unsafe_fields.FetchPrivateField(conf.Session, "handlers")
	r := handlers.MapRange()
	for r.Next() {
		var packetHandler PacketHandler
		if val := r.Value().Interface(); val != nil {
			packetHandler = val.(PacketHandler)
		} else {
			continue
		}
		key := r.Key()
		id := key.Interface().(uint32)
		if constructor := h.constructors[id]; constructor != nil {
			val := reflect.ValueOf(constructor(pl, packetHandler))
			handlers.SetMapIndex(key, val)
		}
	}
}
