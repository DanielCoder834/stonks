package main

import (
	"math"
	"strconv"
	"time"
)

// source: https://karthaus.nl/rdp/

func DouglasPeucker(args []Data, epsilon int) ([]Data, bool) {
	var emptyResults []Data
	dmax := 0.0
	index := 0
	end := len(args)
	for i := 2; i < (end - 1); i++ {
		d, ok := perpendicularDistance(args[0].Time, args[i], args[1], args[end])
		if !ok {
			return emptyResults, false
		}
		if d > dmax {
			index = i
			dmax = d
		}
	}
	var result []Data
	if int(math.Floor(dmax)) < epsilon {
		recResults1, ok1 := DouglasPeucker(args[1:index], epsilon)
		recResults2, ok2 := DouglasPeucker(args[index:end], epsilon)

		if !ok1 || !ok2 {
			return emptyResults, false
		}
		minusOneRec1Len := len(recResults1) - 1
		result = append(recResults1[1:minusOneRec1Len], recResults2[1:len(recResults2)]...)
	} else {
		result = args[1:end]
	}
	return result, true
}

func perpendicularDistance(firstXPoint time.Time, farawayPoint Data, startLinePoint Data, endLinePoint Data) (float64, bool) {
	farawayPointX := farawayPoint.Time.Sub(firstXPoint).Seconds()
	farawayPointY, errFar := strconv.ParseFloat(farawayPoint.Value.bulk, 64)
	if errFar != nil {
		return 0, false
	}

	startLinePointX := startLinePoint.Time.Sub(firstXPoint).Seconds()
	startLinePointY, errStart := strconv.ParseFloat(startLinePoint.Value.bulk, 64)
	if errStart != nil {
		return 0, false
	}

	endLinePointX := endLinePoint.Time.Sub(firstXPoint).Seconds()
	endLinePointY, errEnd := strconv.ParseFloat(endLinePoint.Value.bulk, 64)
	if errEnd != nil {
		return 0, false
	}

	area := math.Abs(.5 * (startLinePointY*endLinePointX + endLinePointY*farawayPointX +
		farawayPointY*startLinePointX - endLinePointY*startLinePointX -
		farawayPointY*endLinePointX - startLinePointY*farawayPointX))
	bottom := math.Sqrt(math.Pow(startLinePointY-endLinePointY, 2) + math.Pow(startLinePointX-endLinePointX, 2))
	return area / bottom * 2, true
}
