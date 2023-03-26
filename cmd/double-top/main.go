package main

import (
	"github.com/godoji/algocore/pkg/algo"
	"github.com/godoji/algocore/pkg/env"
	"github.com/godoji/algocore/pkg/ritmic"
	"github.com/northberg/candlestick"
	"math"
)

var Params = []string{"historySize"}

type LocalStore struct {
	Initialized bool
	History     *env.FiLoStack
}

func Evaluate(chart env.MarketSupplier, res *algo.ResultHandler, mem *env.Memory, param env.Parameters) {

	// A way of loading memory from disk
	var store *LocalStore
	if tmp := mem.Read(); tmp == nil {
		store = new(LocalStore)
	} else {
		store = tmp.(*LocalStore)
	}
	defer mem.Store(store)

	if !store.Initialized {
		store.History = env.NewFiLoStack(3)
		store.Initialized = true
	}

	highsAndLows := chart.Algorithm("highs-and-lows", param.Get("historySize"))
	if highsAndLows.HasEvents() {
		lastEvent := highsAndLows.CurrentEvents()[0]
		store.History.Push(lastEvent)
	}

	if !store.History.IsFull() {
		return
	}

	t1 := store.History.At(0).(*algo.Event)
	b1 := store.History.At(1).(*algo.Event)
	t2 := store.History.At(2).(*algo.Event)

	maxTop := math.Max(t1.Price, t2.Price)

	if !(t1.Label == "high" && b1.Label == "low" && t2.Label == "high") {
		return
	}

	if t1.Price > t2.Price*1.01 || t1.Price < t2.Price/1.01 {
		return
	}

	ch := chart.Interval(candlestick.Interval1d)
	leftIndex := ch.ToIndex(t1.Time)
	rightIndex := ch.ToIndex(t2.Time)

	for i := rightIndex; i <= leftIndex; i++ {
		c := ch.FromLast(int(i))
		if c.Time == t1.Time || c.Time == t2.Time || c.Time == b1.Time {
			continue
		}
		if c.Missing {
			continue
		}
		if c.High > maxTop {
			res.NewEvent("bad").SetColor("orange").SetTime(c.Time).SetPrice(c.High)
			return
		}
		if c.Low < b1.Price {
			res.NewEvent("bad").SetColor("red").SetTime(c.Time).SetPrice(c.Low)
			return
		}
	}

	var leftCandle *candlestick.Candle
	for i := 1; i < 200; i++ {
		j := int(leftIndex) + i
		c := ch.FromLast(j)
		if c.Missing {
			continue
		}
		if c.High >= maxTop {
			break
		}
		if c.Low <= b1.Price {
			leftCandle = c
			break
		}
	}
	if leftCandle == nil {
		return
	}

	var rightCandle *candlestick.Candle
	for i := 1; i < 200; i++ {
		j := int(rightIndex) - i
		if j < 0 {
			break
		}
		c := ch.FromLast(j)
		if c.Missing {
			continue
		}
		if c.High >= maxTop {
			break
		}
		if c.Low <= b1.Price {
			rightCandle = c
			break
		}
	}
	if rightCandle == nil {
		return
	}

	event := res.NewEvent("double-top").SetColor("red").SetIcon("down")
	event.AddSegment(&algo.SegmentAnnotation{
		TimeBegin:  leftCandle.Time,
		TimeEnd:    rightCandle.Time,
		PriceBegin: b1.Price,
		PriceEnd:   maxTop,
		Style:      "region",
		Color:      "rgb(100,100,100)",
	})
	event.AddEvent(t1, "t1")
	event.AddEvent(b1, "b1")
	event.AddEvent(t2, "t2")
	store.History.Clear()
}

func main() {
	ritmic.Serve(Evaluate, Params)
}
