package main

import (
	"fmt"

	"github.com/FDUTCH/better_dragonfly/scheduler"
	"github.com/df-mc/dragonfly/server/world"
)

func main() {
	w := world.Config{}.New()
	defer w.Close()
	s := scheduler.Scheduler[string, string]{World: w}

	result := s.PipeExec("hi %s !" /* passing param */, func(tx *world.Tx, param string) string {
		// inner logic...
		return fmt.Sprintf(param, "FDUTCH")
	})
	fmt.Println(result)
}
