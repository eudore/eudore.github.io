package config

import (
	"fmt"
	"sort"
)

var reloads []reloadfn

type reloadfn struct {
	name  string
	index int
	fn    func() error
	exec  bool
}

// Reload all config func
func ReloadAll(cs ...string) error {
	sort.SliceStable(reloads, func(i, j int) bool {
		return reloads[i].index < reloads[j].index
	})
	//fmt.Println(len(cs),len(reloads))
	if len(cs) == 0 {
		for _, i := range reloads {
			err := i.fn()
			fmt.Println("reload", i.name)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		reloads[0].fn()
		sort.Search(len(reloads), func(i int) bool {
			return true
		})
	}
	return nil
}

// Set reload func
func SetReload(name string, index int, fn func() error) {
	reloads = append(reloads, reloadfn{name: name, index: index, fn: fn})
}
