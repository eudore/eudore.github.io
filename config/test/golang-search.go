package main;

import (
	"fmt"
	"github.com/dwayhs/go-search-engine"
	"github.com/dwayhs/go-search-engine/core"
	"github.com/dwayhs/go-search-engine/analysis/analyzers"
)

var index *gosearchengine.Index;

func init() {
	index = gosearchengine.NewIndex(
		gosearchengine.Mapping{
			Attributes: map[string]analyzers.Analyzer{
				"body": analyzers.NewSimpleAnalyzer(),
			},
		},
	)
}

func doc(data string) (core.Document){
	docA := core.Document{
		Attributes: map[string]string{
    		"body": data, 
		},
	}
	return docA
}


func main() {
	index.Index(doc("The quick brown fox jumps over the lazy dog fkdlf fd"))
	index.Index(doc("The quick fox jumps fd"))
	index.Index(doc("The quick fox jumps over the lazy dog fkdlf fd"))
	index.Index(doc("The quick dlf fd"))
	index.Index(doc("The quick brown fox jumps over the lazy dog"))
//	index.Index(doc("native Group Immediate Order By Id takes group whichGroup, integer order returns boolean"))
//	index.Index(doc("native Group PointOrder takes group whichGroup, string order, real x, real y returns boolean"))
//	index.Index(doc("native Group Point Order Loc takes group whichGroup, string order, location whichLocation returns boolean"))
//	index.Index(doc("native Group Point Order By Id takes group whichGroup, integer order, real x, real y returns boolean"))
//	index.Index(doc("native Group Target Order takes group whichGroup, string order, widget targetWidget returns boolean"))
//	index.Index(doc("native Group Point Order By Id Loc takes group whichGroup, integer order, location whichLocation returns boolean"))
	searchResult := index.Search("body", "quick")
	for k,v := range searchResult {
		fmt.Println(k,v)
	}
}
