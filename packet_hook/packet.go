package packet_hook

import (
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PacketHandler interface {
	Handle(p packet.Packet, s *session.Session, tx *world.Tx, c session.Controllable) error
}
