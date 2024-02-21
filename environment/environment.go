// Package environment provides some utilities for working with base images.
package environment

import (
	"log/slog"
	"sync"
)

// SoftwareList is the map that contains the list of
// the software packages installed in this environment.
type SoftwareList map[string]string

// DetermineSoftwareList determines this environment to get the SoftwareList.
func DetermineSoftwareList() SoftwareList {
	softwares := GetSoftwares()

	type item struct {
		key   string
		value string
	}
	itemCh := make(chan item)

	wg := sync.WaitGroup{}
	wg.Add(len(softwares))

	for _, s := range softwares {
		go func(s Software) {
			defer wg.Done()

			slog.Info("determine software", slog.String("name", s.Name()))

			version, ok := s.Version()
			if ok {
				itemCh <- item{s.Name(), version}
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(itemCh)
	}()

	softwareList := make(SoftwareList)
	for item := range itemCh {
		softwareList[item.key] = item.value
	}

	return softwareList
}
