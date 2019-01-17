package movingaverage
import (
	"time"
)
// @author Robin Verlangen
// Moving average implementation for Go

type MovingAverage struct {
	Window int
	Values []Values
	valPos int
	slotsFilled bool
}
type Values struct {
	Value float64
	Time time.Time
}

func (ma *MovingAverage) Avg() float64 {
	var sum = float64(0)
	var c = ma.Window-1

	// Are all slots filled? If not, ignore unused
	if !ma.slotsFilled {
		c = ma.valPos-1
		if c < 0 {
			// Empty register
			return 0
		}
	}

	// Sum Values
	var ic = 0
	for i := 0; i <= c; i++ {
		sum += ma.Values[i].Value
		ic++
	}

	// Finalize average and return
	avg := sum / float64(ic)
	return avg
}

func (ma *MovingAverage) Add(val Values) {
	// Put into Values array
	ma.Values[ma.valPos] = val

	// Increment value position
	ma.valPos = (ma.valPos + 1) % ma.Window

	// Did we just go back to 0, effectively meaning we filled all registers?
	if !ma.slotsFilled && ma.valPos == 0 {
		ma.slotsFilled = true
	}
}

func New(window int) *MovingAverage {
	return &MovingAverage{
		Window : window,
		Values : make([]Values, window),
		valPos : 0,
		slotsFilled : false,
	}
}
