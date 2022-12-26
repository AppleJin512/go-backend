package stats

import (
	"moonbite/trending/internal/config"
	"sync"

	"github.com/urfave/cli/v2"
)

func Command(c *cli.Context) error {
	_, nodeListeners, err := NewNodeListeners(c.Context, config.Config.Instances)
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	for _, nl := range nodeListeners {
		wg.Add(1)
		go nl.Listen(wg)
	}
	wg.Wait()
	return nil
}
