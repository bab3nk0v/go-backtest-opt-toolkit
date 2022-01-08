package objectives

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"trade-optimizer/util"
)

func TestAvgDrawdown(t *testing.T) {
	data := []float64{
		0.12153621779289538,
		-14.609256280694693,
		-13.221099541984449,
		-12.549395091942465,
	}
	avgDd := 121.20496067756277
	res := util.AvgDrawdown(data)
	assert.InDelta(t, avgDd, res, 0.001)
}

func TestAvgDrawdownLast(t *testing.T) {
	data := []float64{
		0.12153621779289538,
		14.609256280694693,
		-13.221099541984449,
		-14.549395091942465,
	}
	avgDd := 1.99590
	res := util.AvgDrawdown(data)
	assert.InDelta(t, avgDd, res, 0.001)
}

func TestAvgZigZag(t *testing.T) {
	data := []float64{
		1, -2, 3, -4, 5, -6, 7, -8,
	}
	avgDd := 2.419047619047619
	res := util.AvgDrawdown(data)
	assert.InDelta(t, avgDd, res, 0.001)
}

func TestNan(t *testing.T) {
	data := []float64{
		-18.406072106261856,
		-17.51378402468123,
		-0.17402855969419662,
		2.0793426255435534,
		2.1030336940107532,
	}
	res := util.AvgDrawdown(data)
	assert.InDelta(t, 0.0, res, 0.001)
}
