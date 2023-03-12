package main

import (
	"github.com/godoji/algocore/pkg/algo"
	"github.com/godoji/algocore/pkg/env"
	"github.com/godoji/algocore/pkg/ritmic"
	candles "github.com/northberg/candlestick"
	"log"
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

	// Initialize any memory, or append
	if !store.Initialized {
		histSize := param.GetInt("historySize")
		store.History = env.NewFiLoStack(histSize)
		for i := 0; i < histSize; i++ {
			c := chart.Interval(candles.Interval1d).FromLast(i)
			store.History.Push(c)
		}
		store.Initialized = true
	} else {
		c := chart.Interval(candles.Interval1d).Candle()
		store.History.Push(c)
	}

	// Calculate local maxima
	values := make([]float64, 0)
	for _, i := range store.History.ToSlice() {
		c := i.(*candles.Candle)
		if !c.Missing {
			values = append(values, c.High, c.Low)
		}
	}

	if len(values) < 100 {
		return
	}

	sort.Float64s(values)

	// Create some events
	centroids := make([]*Centroid, 0)
	k := 5
	step := len(values) / k
	for i := 0; i < 5; i++ {
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
			newCentroids = append(newCentroids, NewCentroid(centroid.Mean()))
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
	sum            float64
	size           int
}

func NewCentroid(v float64) *Centroid {
	return &Centroid{
		values:         make([]float64, 0),
		representative: v,
		sum:            0,
		size:           0,
	}
}

func (c *Centroid) Add(v float64) {
	c.values = append(c.values, v)
	c.sum += v
	c.size++
}

func (c *Centroid) Center() float64 {
	return c.representative
}

func (c *Centroid) Mean() float64 {
	if c.size == 0 {
		log.Fatalln("invalid mean size")
	}
	return c.sum / float64(c.size)
}

func main() {
	ritmic.Serve(Evaluate, Params)
}
