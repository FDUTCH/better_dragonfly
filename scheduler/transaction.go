package scheduler

import (
	"unsafe"

	"github.com/df-mc/dragonfly/server/world"
)

type Tx struct {
	w      *world.World
	closed bool
}

func (tx *Tx) Get() *world.Tx {
	return (*world.Tx)(unsafe.Pointer(tx))
}

func (tx *Tx) close() {
	tx.closed = true
}

func pushToQueue(w *world.World, t transaction) {
	basePtr := unsafe.Pointer(w)
	queuePtr := unsafe.Pointer(uintptr(basePtr) + queueOffset)
	queue := *(*chan transaction)(queuePtr)
	queue <- t
}

func init() {
	queueOffset = mustFieldOffset[world.World]("queue")
}

var queueOffset uintptr

type transaction interface {
	Run(w *world.World)
}
