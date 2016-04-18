package wktparse

// LineString, LineString Z, LineString M, LineString ZM
//     strip brackets
//     split on comma
//     x, y, z, m
//
// Polygon
//     If ),( present
//         Split on ),(
//     else
//         Split on comma
//
// MultiPoint
//     Strip all brackets
//     Split on comma
//
// MultiLineString
//     Split on ),(
//         Strip all brackets
//         Split on comma
//
// TIN, TIN Z, TIN M, TIN ZM

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
  "reflect"
)

type Coordinate struct {
	X float64
	Y float64
	Z float64
	M float64
}

type CoordinateSet struct {
	Coordinates []Coordinate //(0 0 0, 0 1 0, 1 1 0, 1 0 0, 0 0 0)
	Holes       []Coordinate
}

// Define WTK
type WKT struct {
	WTKType    string
	Geometries []CoordinateSet
}

func ParseGeometry(WKTString string) (string, CoordinateSet) {

	var WKType string
	var Geometries CoordinateSet

	WKTString = strings.ToUpper(WKTString)

	if strings.HasPrefix(WKTString, "POINT") {

		WKType, Geometries = Point(WKTString, "POINT")

	} else if strings.HasPrefix(WKTString, "LINESTRING") {

		WKType, Geometries = Line(WKTString, "LINESTRING")

	} else if strings.HasPrefix(WKTString, "POLYGON") {

		WKType, Geometries = Polygon(WKTString, "POLYGON")
		//Geometries = Point(WKTString)
	} else if strings.HasPrefix(WKTString, "MULTIPOINT") {
		//Geometries = Point(WKTString)
	} else if strings.HasPrefix(WKTString, "MULTILINESTRING") {
		//Geometries = Point(WKTString)
	} else if strings.HasPrefix(WKTString, "MULTIPOLYGON") {
		//Geometries = Point(WKTString)
	}

	return WKType, Geometries

}

// POINT, POINT M, POINT Z, POINT ZM
// POINT (6 10)
func Point(WKTString string, ParentType string) (string, CoordinateSet) {

	point := strings.Split(RemoveWrappingGeom(WKTString), " ")
	var coords []float64
	for _, c := range point {
		if f, err := strconv.ParseFloat(c, 64); err == nil {
			coords = append(coords, f)
		}
	}

	wkttype := GetGeometryType(WKTString, ParentType)
	coordinateset := GetCoordinate(coords, wkttype)
	coordinates := CoordinateSet{Coordinates: []Coordinate{coordinateset}}

	return wkttype, coordinates
}

//LINESTRING (30 10, 10 30, 40 40)
func Line(WKTString string, ParentType string) (string, CoordinateSet) {

	coordinateset := CoordinateSet{}
	var wkttype string

	wkttype = GetGeometryType(WKTString, ParentType)

	lineparts := strings.Split(RemoveWrappingGeom(WKTString), ",")

	// For each part of the line
	for _, l := range lineparts {

		l = strings.TrimSpace(l)
		line := strings.Split(l, " ")
		var coords []float64

		// For each coordinate in the line
		for _, c := range line {
			c = strings.TrimSpace(c)
			if c != " " || c != "" {
				if f, err := strconv.ParseFloat(c, 64); err == nil {
					coords = append(coords, f)
				} else {
					fmt.Print("Error converting ", c, "to float64", " ", reflect.TypeOf(c), "\n ")
				}
			}
		}

		if (len(coords) > 0) {

			coordinateset.Coordinates = append(coordinateset.Coordinates, GetCoordinate(coords, wkttype))

		} else {

			fmt.Print("Not enough coordinates in this line :", WKTString)

		}

	}

	return wkttype, coordinateset
}

//POLYGON ((35 10, 45 45, 15 40, 10 20, 35 10), (20 30, 35 35, 30 20, 20 30))
func Polygon(WKTString string, ParentType string) (string, CoordinateSet) {

	coordinateset := CoordinateSet{}
	var wkttype string

	wkttype = GetGeometryType(WKTString, ParentType)

	polygonpartsstr := RemoveWrappingGeom(WKTString)
	polygonparts := []string{}

	if strings.Contains(polygonpartsstr, "),(") {

		polygonparts = strings.Split(polygonpartsstr, "),(") // We remove other brackets later

	} else {

		polygonparts = []string{ polygonpartsstr }

	}

	// For each part of the polgon
	for _, part := range polygonparts {

		part = RemoveWrappingGeom(part)
		partslice := strings.Split(part, ",")

		for _, coordstr := range partslice {

			coordstr = strings.TrimSpace(coordstr)
			coordstr = RemoveAllBrackets(coordstr)
			line := strings.Split(coordstr, " ")
			var coords []float64

			// For each coordinate in the line
			for _, c := range line {
				c = strings.TrimSpace(c)
				if c != " " || c != "" {
					if f, err := strconv.ParseFloat(c, 64); err == nil {
						coords = append(coords, f)
					} else {
						fmt.Print("Error converting ", c, "to float64", " ", reflect.TypeOf(c), "\n ")
					}
				}
			}

			if (len(coords) > 0 && len(coordinateset.Coordinates) < 1) {

				coordinateset.Coordinates = append(coordinateset.Coordinates, GetCoordinate(coords, wkttype))

			} else if (len(coords) > 0 && len(coordinateset.Coordinates) >= 1) {

				coordinateset.Holes =  append(coordinateset.Holes, GetCoordinate(coords, wkttype))

			} else {
				fmt.Print("Not enough coordinates in this line :", WKTString)
			}

		}

	}

	return wkttype, coordinateset

}

func GetCoordinate(coords []float64, wkttype string) Coordinate {

	var coordinate Coordinate

	if strings.Contains(wkttype, "Z") && !strings.Contains(wkttype, "ZM") {

		coordinate = Coordinate{X: coords[0], Y: coords[1], Z: coords[2]}

	} else if strings.Contains(wkttype, "ZM") {

		coordinate = Coordinate{X: coords[0], Y: coords[1], Z: coords[2], M: coords[3]}

	} else if strings.Contains(wkttype, "M") {

		coordinate = Coordinate{X: coords[0], Y: coords[1], M: coords[2]}

	} else {

		coordinate = Coordinate{X: coords[0], Y: coords[1]}

	}

	return coordinate

}

func GetGeometryType(WKTString string, ParentType string) string {

	var wkttype string

	if strings.HasPrefix(WKTString, ParentType + " Z") && !strings.HasPrefix(WKTString, ParentType + " ZM") {

		wkttype = ParentType + " Z"

	} else if strings.HasPrefix(WKTString, ParentType + " M") {

		wkttype = ParentType + " M"

	} else if strings.HasPrefix(WKTString, ParentType + " ZM") {

		wkttype = ParentType + " ZM"

	} else {

		wkttype = ParentType

	}

	return wkttype

}


func RemoveWrappingGeom(str string) string {
	if strings.Contains(str, "(") && strings.Contains(str, ")") {
		return str[strings.Index(str, "(")+1 : strings.LastIndex(str, ")")]
	} else {
		return str
	}

}

func RemoveAllAlphabet(str string) string {
	reg, _ := regexp.Compile("[A-Za-z]")
	return reg.ReplaceAllString(str, "")
}

func RemoveAllBrackets(str string) string {
	str = strings.Replace(str, "(", "", -1)
	str = strings.Replace(str, ")", "", -1)
	return str
}
