package main

import (
	"github.com/godoji/algocore/pkg/algo"
	"github.com/godoji/algocore/pkg/env"
	"github.com/godoji/algocore/pkg/ritmic"
	candles "github.com/northberg/candlestick"
	"math"
	"sort"
)

var Params = []string{"historySize"}

type LocalStore struct {
	Initialized bool
	History     *env.FiLoStack
}

type PriceVolumePair struct {
	Price  float64
	Volume float64
	Time   int64
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
		histSize := param.GetInt("historySize")
		store.History = env.NewFiLoStack(histSize)
		store.Initialized = true
	}

	highsAndLows := chart.Algorithm("highs-and-lows", 7)
	if highsAndLows.HasEvents() {
		for _, event := range highsAndLows.CurrentEvents() {
			vol := 0.0
			for i := 0; i < 7; i++ {
				vol += chart.Interval(candles.Interval1d).FromLast(0).Volume
			}
			res.NewEvent("high").SetPrice(event.Price).SetTime(event.Time).SetColor("green").SetIcon("top")
			store.History.Push(&PriceVolumePair{
				Price:  event.Price,
				Volume: vol,
				Time:   event.Time,
			})
		}
	}

	// Calculate local maxima
	weights := make(map[float64]float64, 0)
	for _, i := range store.History.ToSlice() {
		if i == nil {
			break
		}
		c := i.(*PriceVolumePair)
		if chart.Time()-c.Time > 90*candles.Interval1d {
			continue
		}
		if _, ok := weights[c.Price]; ok {
			weights[c.Price] += c.Volume
		} else {
			weights[c.Price] = c.Volume
		}
	}

	if len(weights) == 0 {
		return
	}

	values := make([]float64, len(weights))
	{
		i := 0
		for v := range weights {
			values[i] = v
			i++
		}
	}
	sort.Float64s(values)

	// Create some events
	centroids := make([]*Centroid, 0)
	k := 4
	step := len(values) / k
	for i := 0; i < k; i++ {
		pseudoRandomIndex := i * step
		centroids = append(centroids, NewCentroid(values[pseudoRandomIndex]))
	}

	iterations := 10
	for i := 0; i < iterations; i++ {
		for _, v := range values {
			lowestDist := math.MaxFloat64
			bestCentroid := centroids[0]
			for _, centroid := range centroids {
				dist := math.Abs(centroid.Center() - v)
				if dist < lowestDist {
					lowestDist = dist
					bestCentroid = centroid
				}
			}
			bestCentroid.Add(v)
		}
		newCentroids := make([]*Centroid, 0)
		for _, centroid := range centroids {
			newCentroids = append(newCentroids, NewCentroid(centroid.Mean(weights)))
		}
		centroids = newCentroids
	}

	for _, cluster := range centroids {
		res.NewEvent("level").AddSegment(&algo.SegmentAnnotation{
			TimeBegin:  chart.Time(),
			TimeEnd:    chart.Time() + 180*candles.Interval1d,
			PriceBegin: cluster.Center(),
			PriceEnd:   cluster.Center(),
			Style:      "solid",
			Color:      "rgba(255,255,255,0.025)",
		})
	}
}

type Centroid struct {
	values         []float64
	representative float64
	size           int
}

func NewCentroid(v float64) *Centroid {
	return &Centroid{
		values:         make([]float64, 0),
		representative: v,
		size:           0,
	}
}

func (c *Centroid) Add(v float64) {
	c.values = append(c.values, v)
	c.size++
}

func (c *Centroid) Center() float64 {
	return c.representative
}

func (c *Centroid) Mean(weights map[float64]float64) float64 {
	if c.size == 0 {
		return 0
	}
	sumWeights := 0.0
	for _, v := range c.values {
		sumWeights += weights[v]
	}
	total := 0.0
	for _, v := range c.values {
		part := weights[v] / sumWeights
		total += v * part
	}
	return total
}

func main() {
	ritmic.Serve(Evaluate, Params)
}
