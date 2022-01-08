package dataloader

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strconv"
)

type Candle struct {
	Time   uint64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	Trades int64
}

type Candles []Candle

func (c Candles) Open() []float64 {
	r := make([]float64, 0, len(c))
	for _, c := range c {
		r = append(r, c.Open)
	}
	return r
}

func (c Candles) High() []float64 {
	r := make([]float64, 0, len(c))
	for _, c := range c {
		r = append(r, c.High)
	}
	return r
}

func (c Candles) Low() []float64 {
	r := make([]float64, 0, len(c))
	for _, c := range c {
		r = append(r, c.Low)
	}
	return r
}

func (c Candles) Close() []float64 {
	r := make([]float64, 0, len(c))
	for _, c := range c {
		r = append(r, c.Close)
	}
	return r
}

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "error in reading input csv")
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse file as CSV for "+filePath)
	}
	return records, nil
}

func toCandleArray(candleStrings [][]string, timeMillis bool) ([]Candle, error) {
	candles := make([]Candle, 0, len(candleStrings))
	for _, candleString := range candleStrings {
		ts, err := strconv.ParseUint(candleString[0], 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "error while parsing timestamp from field "+candleString[0])
		}
		open, err := strconv.ParseFloat(candleString[1], 64)
		if err != nil {
			return nil, errors.Wrap(err, "error while parsing open price from field "+candleString[1])
		}
		high, err := strconv.ParseFloat(candleString[2], 64)
		if err != nil {
			return nil, errors.Wrap(err, "error while parsing high price from field "+candleString[2])
		}
		low, err := strconv.ParseFloat(candleString[3], 64)
		if err != nil {
			return nil, errors.Wrap(err, "error while parsing low price from field "+candleString[3])
		}
		cloze, err := strconv.ParseFloat(candleString[4], 64)
		if err != nil {
			return nil, errors.Wrap(err, "error while parsing close price from field "+candleString[4])
		}
		volume, err := strconv.ParseFloat(candleString[5], 64)
		if err != nil {
			return nil, errors.Wrap(err, "error while parsing volume from field "+candleString[5])
		}
		trades := int64(-1)
		if len(candleString) >= 7 {
			trades, err = strconv.ParseInt(candleString[6], 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "error while parsing num trades from field "+candleString[6])
			}
		}
		if timeMillis {
			ts /= 1000
		}
		candles = append(candles, Candle{
			Time:   ts,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  cloze,
			Volume: volume,
			Trades: trades,
		})
	}
	return candles, nil
}

func ReadOHLCV(path string, timeMillis bool) ([]Candle, error) {
	candleStrings, err := readCsvFile(path)
	if err != nil {
		return nil, err
	}
	candles, err := toCandleArray(candleStrings, timeMillis)
	if err != nil {
		return nil, err
	}
	sort.Slice(candles, func(i, j int) bool {
		return candles[i].Time < candles[j].Time
	})
	return candles, nil
}
