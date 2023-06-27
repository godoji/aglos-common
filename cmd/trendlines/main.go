package main

import (
	"github.com/godoji/algocore/pkg/algo"
	"github.com/godoji/algocore/pkg/env"
	"github.com/godoji/algocore/pkg/ritmic"
)

var Params = []string{"historySize"}

type TrendLine struct {
}

type LocalStore struct {
	ExistingLines []*TrendLine
}

func Evaluate(chart env.MarketSupplier, res *algo.ResultHandler, mem *env.Memory, param env.Parameters) {

	var store *LocalStore
	if tmp := mem.Read(); tmp == nil {
		store = &LocalStore{
			ExistingLines: make([]*TrendLine, 0),
		}
	} else {
		store = tmp.(*LocalStore)
	}
	defer mem.Store(store)

	highsAndLows := chart.Algorithm("highs-and-lows", param.Get("historySize"))
	if !highsAndLows.HasEvents() {
		return
	}

	lastEvent := highsAndLows.LastEvents()[0]
	previousEvents := highsAndLows.PastEvents()

	if len(previousEvents) == 0 {
		return
	}

	//color := "green"
	//if lastEvent.Label == "low" {
	//	color = "red"
	//}
	//counter := 0
	//for i := 0; i < len(previousEvents); i++ {
	//	prev := previousEvents[len(previousEvents)-i-1]
	//	if counter > 1 {
	//		break
	//	}
	//	res.NewEvent("trendline").AddSegment(&algo.SegmentAnnotation{
	//		TimeBegin:  prev.Time,
	//		TimeEnd:    lastEvent.Time,
	//		PriceBegin: prev.Price,
	//		PriceEnd:   lastEvent.Price,
	//		Style:      "solid-line",
	//		Color:      color,
	//	}).AddEvent(lastEvent, lastEvent.Label)
	//	counter++
	//}

}

func main() {
	ritmic.Serve(Evaluate, Params)
}
