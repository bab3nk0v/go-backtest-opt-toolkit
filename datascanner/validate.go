package datascanner

import (
	"fmt"
	"github.com/pkg/errors"
	"trade-optimizer/dataloader"
)

var errZeroVolume = errors.New("Critical amount of zero volume candles")
var errNonSorted = errors.New("Your candle array is not sorted")

func AnalyzeCandles(candles dataloader.Candles) (bool, error) {
	zeroVolume := 0
	for _, x := range candles {
		if x.Volume == 0 {
			zeroVolume++
		}
	}
	ratio := float64(zeroVolume) / float64(len(candles))
	if ratio > 0.1 {
		return false, errors.Wrap(errZeroVolume, fmt.Sprintf("%f of zero volume candles found", ratio))
	}

	for i, x := range candles {
		if i > 0 {
			if x.Time < candles[i-1].Time {
				return false, errNonSorted
			}
		}
	}

	return true, nil
}
