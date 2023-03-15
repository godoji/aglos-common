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

	// Add candle to stack
	{
		c := chart.Interval(candles.Interval1d).Candle()
		store.History.Push(c)
	}

	// Calculate local maxima
	weights := make(map[float64]float64, 0)
	for _, i := range store.History.ToSlice() {
		if i == nil {
			break
		}
		c := i.(*candles.Candle)
		if !c.Missing {
			if _, ok := weights[c.High]; ok {
				weights[c.High] += c.Volume
			} else {
				weights[c.High] = c.Volume
			}
			if _, ok := weights[c.Low]; ok {
				weights[c.Low] += c.Volume
			} else {
				weights[c.Low] = c.Volume
			}
		}
	}

	if len(weights) < 50 {
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

	iterations := 20
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
		res.NewEvent("level").SetColor("blue").SetIcon("h").SetPrice(cluster.Center())
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
