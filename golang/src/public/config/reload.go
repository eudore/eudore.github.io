package config

import (
	"fmt"
	"sort"
)


var reloads		[]reloadfn

type reloadfn struct {
	name	string
	index	int
	fn		func() error
}

func ReloadAll(cs ...string) error {
	sort.SliceStable(reloads,func(i, j int) bool {
		return reloads[i].index < reloads[j].index
	})
	fmt.Println(len(cs),len(reloads))
	fmt.Println(reloads)
	if len(cs) == 0 {
		for _,i := range reloads {
			i.fn()
			fmt.Println("reload", i.name)
		}
	}else {
		reloads[0].fn()
		sort.Search(len(reloads), func(i int) bool {
			return true
		})
	}
	return nil
}

func SetReload(name string, index int, fn func() error) {
	reloads = append(reloads, reloadfn{name: name,index: index,fn: fn})
}