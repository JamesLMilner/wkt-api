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
	if strings.HasPrefix(WKTString, "POINT") {
		WKType, Geometries = Point(WKTString)
	} else if strings.HasPrefix(WKTString, "LINESTRING") {
		WKType, Geometries = Line(WKTString)
	} else if strings.HasPrefix(WKTString, "POLYGON") {
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
func Point(WKTString string) (string, CoordinateSet) {

	point := strings.Split(RemoveWrappingGeom(WKTString), " ")
	var coords []float64
	for _, c := range point {
		if f, err := strconv.ParseFloat(c, 64); err == nil {
			coords = append(coords, f)
		}
	}

	var wkttype string
	coordinateset := Coordinate{}

	if strings.HasPrefix(WKTString, "POINT Z") && !strings.HasPrefix(WKTString, "POINT ZM") {
		wkttype = "POINT Z"
		coordinateset = Coordinate{X: coords[0], Y: coords[1], Z: coords[2]}
	} else if strings.HasPrefix(WKTString, "POINT M") {
		wkttype = "POINT M"
		coordinateset = Coordinate{X: coords[0], Y: coords[1], M: coords[2]}
	} else if strings.HasPrefix(WKTString, "POINT ZM") {
		wkttype = "POINT ZM"
		coordinateset = Coordinate{X: coords[0], Y: coords[1], Z: coords[2], M: coords[3]}
	} else {
		wkttype = "POINT"
		coordinateset = Coordinate{X: coords[0], Y: coords[1]}
	}

	coordinates := CoordinateSet{Coordinates: []Coordinate{coordinateset}}
	return wkttype, coordinates
}

//LINESTRING (30 10, 10 30, 40 40)
func Line(WKTString string) (string, CoordinateSet) {

	coordinateset := CoordinateSet{}
	var wkttype string

	if strings.HasPrefix(WKTString, "LINESTRING Z") && !strings.HasPrefix(WKTString, "LINESTRING ZM") {
		wkttype = "LINESTRING Z"
	} else if strings.HasPrefix(WKTString, "LINESTRING M") {
		wkttype = "LINESTRING M"
	} else if strings.HasPrefix(WKTString, "LINESTRING ZM") {
		wkttype = "LINESTRING ZM"
	} else {
		wkttype = "LINESTRING"
	}

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

			if wkttype == "LINESTRING Z" {
				coordinateset.Coordinates = append(coordinateset.Coordinates, Coordinate{X: coords[0], Y: coords[1], Z: coords[2]})
			} else if wkttype == "LINESTRING M" {
				coordinateset.Coordinates = append(coordinateset.Coordinates, Coordinate{X: coords[0], Y: coords[1], M: coords[2]})
			} else if wkttype == "LINESTRING ZM" {
				coordinateset.Coordinates = append(coordinateset.Coordinates, Coordinate{X: coords[0], Y: coords[1], Z: coords[2], M: coords[3]})
			} else if wkttype == "LINESTRING" {
				coordinateset.Coordinates = append(coordinateset.Coordinates, Coordinate{X: coords[0], Y: coords[1]})
			}
		} else {
			fmt.Print("Not enough coordinates in this line :", WKTString)
		}

	}

	return wkttype, coordinateset
}

//LINESTRING (30 10, 10 30, 40 40)
func Polygon(WKTString string) (string, CoordinateSet) {

	// coordinateset := CoordinateSet{}
	// var wkttype string
	//
	// if strings.HasPrefix(WKTString, "LINESTRING Z") && !strings.HasPrefix(WKTString, "LINESTRING ZM") {
	// 	wkttype = "LINESTRING Z"
	// } else if strings.HasPrefix(WKTString, "LINESTRING M") {
	// 	wkttype = "LINESTRING M"
	// } else if strings.HasPrefix(WKTString, "LINESTRING ZM") {
	// 	wkttype = "LINESTRING ZM"
	// } else {
	// 	wkttype = "LINESTRING"
	// }
	//
	// lineparts := strings.Split(RemoveWrappingGeom(WKTString), ",")
	//
	// // For each part of the line
	// for _, l := range lineparts {
	//
	// 	l = strings.TrimSpace(l)
	// 	line := strings.Split(l, " ")
	// 	var coords []float64
	//
	// 	// For each coordinate in the line
	// 	for _, c := range line {
	// 		c = strings.TrimSpace(c)
	// 		if c != " " || c != "" {
	// 			if f, err := strconv.ParseFloat(c, 64); err == nil {
	// 				coords = append(coords, f)
	// 			} else {
	// 				fmt.Print("Error converting ", c, "to float64", " ", reflect.TypeOf(c), "\n ")
	// 			}
	// 		}
	// 	}
	//
	// 	if (len(coords) > 0) {
	//
	// 		if wkttype == "LINESTRING Z" {
	// 			coordinateset.Coordinates = append(coordinateset.Coordinates, Coordinate{X: coords[0], Y: coords[1], Z: coords[2]})
	// 		} else if wkttype == "LINESTRING M" {
	// 			coordinateset.Coordinates = append(coordinateset.Coordinates, Coordinate{X: coords[0], Y: coords[1], M: coords[2]})
	// 		} else if wkttype == "LINESTRING ZM" {
	// 			coordinateset.Coordinates = append(coordinateset.Coordinates, Coordinate{X: coords[0], Y: coords[1], Z: coords[2], M: coords[3]})
	// 		} else if wkttype == "LINESTRING" {
	// 			coordinateset.Coordinates = append(coordinateset.Coordinates, Coordinate{X: coords[0], Y: coords[1]})
	// 		}
	// 	} else {
	// 		fmt.Print("Not enough coordinates in this line :", WKTString)
	// 	}
	//
	// }
	//
	// return wkttype, coordinateset
}


func RemoveWrappingGeom(str string) string {
	return str[strings.Index(str, "(")+1 : strings.LastIndex(str, ")")]
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
