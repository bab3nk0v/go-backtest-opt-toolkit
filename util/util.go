package util

import (
	"math"
	"strconv"
	"time"
	"trade-optimizer/dataloader"
)

const (
	MinUint uint = 0                 // binary: all zeroes

	// Perform a bitwise NOT to change every bit from 0 to 1
	MaxUint      = ^MinUint          // binary: all ones

	// Shift the binary number to the right (i.e. divide by two)
	// to change the high bit to 0
	MaxInt       = int(MaxUint >> 1) // binary: all ones except high bit

	// Perform another bitwise NOT to change the high bit to 1 and
	// all other bits to 0
	MinInt       = ^MaxInt           // binary: all zeroes except high bit
)

func CrossedAbove(x []float64, y []float64) bool {
	if (x[len(x)-1] > y[len(y)-1]) && (x[len(x)-2] <= y[len(y)-2]) {
		return true
	}
	return false
}

func CrossedBelow(x []float64, y []float64) bool {
	if x[len(x)-1] < y[len(y)-1] && (x[len(x)-2] >= y[len(y)-2]) {
		return true
	}
	return false
}

func meanUint64(x []uint64) float64 {
	r := 0.0
	for i, t := range x {
		w := float64(i) / float64(i+1)
		r *= w
		r += float64(t) / float64(i+1)
	}
	return r
}

func meanFloat64(x []float64) float64 {
	r := 0.0
	for i, t := range x {
		w := float64(i) / float64(i+1)
		r *= w
		r += float64(t) / float64(i+1)
	}
	return r
}

func SplitAtDate(data []dataloader.Candle, ts uint64) ([]dataloader.Candle, []dataloader.Candle) {
	var splitIndex int
	for i, c := range data {
		if c.Time > ts {
			splitIndex = i
			break
		}
	}
	return data[:splitIndex], data[splitIndex:]
}

func SumFloat64(x []float64) float64 {
	r := 0.0
	for _, y := range x {
		r += y
	}
	return r
}

func ConditionedSumFloat64(x []float64, cond func(float64) bool) float64 {
	r := 0.0
	for _, y := range x {
		if cond(y) {
			r += y
		}
	}
	return r
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func WalkForwardSplit(candles []dataloader.Candle, ddof int, candleMinutes int, nPhases int) [][][]dataloader.Candle {
	rest := candles
	n := len(candles)
	//d := float64(ddof)
	L := 3
	insampleSize := uint64(math.Floor((float64(L) * float64(n)) / (float64(nPhases) + float64(L))))
	outofSampleSize := uint64(math.Floor((float64(n) - float64(insampleSize)) / float64(nPhases)))
	//insampleSize := uint64(math.Floor((float64(n) * (d + 1)) / (d + float64(nPhases) + 1)))
	println("Insample size will be " + strconv.Itoa(int(math.Round(float64(insampleSize)*float64(candleMinutes)/60/24))) + " days")
	//outofSampleSize := uint64(math.Floor(float64(n) / (d + float64(nPhases) + 1)))
	println("Out of sample size will be " + strconv.Itoa(int(math.Round(float64(outofSampleSize)*float64(candleMinutes)/60/24))) + " days")
	var walkForwardData [][][]dataloader.Candle
	for i := 0; i < nPhases; i++ {
		insample := rest[:insampleSize]
		outofsample := rest[insampleSize : insampleSize+outofSampleSize]
		rest = rest[outofSampleSize:]
		walkForwardData = append(walkForwardData, [][]dataloader.Candle{insample, outofsample})
	}
	return walkForwardData
}

func CumSumFloat64(x []float64) []float64 {
	var r []float64
	s := 0.0
	for _, g := range x {
		s += g
		r = append(r, s)
	}
	return r
}

func RunningMaxFloat64(x []float64) ([]float64, []int) {
	var r []float64
	var rmi []int
	m := 0.0
	crmi := 0
	for i, t := range x {
		if t > m {
			m = t
			crmi = i
		}
		r = append(r, m)
		rmi = append(rmi, crmi)
	}
	return r, rmi
}

func MaxDrawdownFloat64(x []float64) float64 {
	cs := CumSumFloat64(x)
	rm, _ := RunningMaxFloat64(cs)
	d := 0.0
	for i, p := range cs {
		if (rm[i] - p) > d {
			d = rm[i] - p
		}
	}
	return d
}

func Eod(t time.Time) time.Time {
	year, month, day := t.Date()
	eod := time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
	return eod
}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	eod := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return eod
}

func ReverseInt64Slice(numbers []int64) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

func ReverseFloat64Slice(numbers []float64) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

func AvgDrawdown(equity []float64) float64 {
	avgDd := 0.0
	nDd := 0
	curDrawdown := 0.0
	curDrawdownStart := 0.0
	if len(equity) < 2 {
		return 0.0
	}
	for i, x := range equity {
		if i > 0 {
			if x < equity[i-1] {
				if curDrawdownStart == 0.0 {
					curDrawdownStart = math.Abs(equity[i-1])
					if curDrawdownStart == 0.0 {
						f64Eps := float64(7.)/3 - float64(4.)/3 - float64(1.)
						curDrawdownStart = f64Eps
					}
				}
				curDrawdown += equity[i-1] - x
			} else if curDrawdown != 0 {
				avgDd += (curDrawdown / curDrawdownStart)
				nDd++
				curDrawdown = 0.0
				curDrawdownStart = 0.0
			}
		}
	}
	if curDrawdown != 0 {
		avgDd += (curDrawdown / curDrawdownStart)
		nDd++
	}
	if nDd == 0.0 {
		return 0.0
	}
	avgDd /= float64(nDd)
	return avgDd
}
