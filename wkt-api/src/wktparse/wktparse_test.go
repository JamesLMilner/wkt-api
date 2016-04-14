package wktparse

import "testing"
import "strconv"
import "strings"
//import "log"

func TestRemoveWrappingGeom(t *testing.T) {
	var wktone string
	var wkttwo string

	wktone = RemoveWrappingGeom("GEOMETRYCOLLECTION(POINT(4 6),LINESTRING(4 6,7 10))")
	wkttwo = RemoveWrappingGeom("TRIANGLE((0 0 0,0 1 0,1 1 0,0 0 0))")

	if wktone != "POINT(4 6),LINESTRING(4 6,7 10)" {
		t.Error("Expected POINT(4 6),LINESTRING(4 6,7 10), got ", wktone)
	}
	if wkttwo != "(0 0 0,0 1 0,1 1 0,0 0 0)" {
		t.Error("Expected (0 0 0,0 1 0,1 1 0,0 0 0), got ", wkttwo)
	}
}

func TestRemoveAllBrackets(t *testing.T) {
	var wktone string
	var wkttwo string

	wktone = RemoveAllBrackets("GEOMETRYCOLLECTION(POINT(4 6),LINESTRING(4 6,7 10))")
	wkttwo = RemoveAllBrackets("TRIANGLE((0 0 0,0 1 0,1 1 0,0 0 0))")

	if wktone != "GEOMETRYCOLLECTIONPOINT4 6,LINESTRING4 6,7 10" {
		t.Error("Expected GEOMETRYCOLLECTIOPOINT4 6,LINESTRING4 6,7 10, got ", wktone)
	}
	if wkttwo != "TRIANGLE0 0 0,0 1 0,1 1 0,0 0 0" {
		t.Error("Expected TRIANGLE0 0 0,0 1 0,1 1 0,0 0 0, got ", wkttwo)
	}
}

func TestRemoveAllAlphabet(t *testing.T) {
	var wktone string
	var wkttwo string

	wktone = RemoveAllAlphabet("GEOMETRYCOLLECTION(POINT(4 6),LINESTRING(4 6,7 10))")
	wkttwo = RemoveAllAlphabet("TRIANGLE((0 0 0,0 1 0,1 1 0,0 0 0))")

	if wktone != "((4 6),(4 6,7 10))" {
		t.Error("Expected ((4 6),(4 6,7 10)) got ", wktone)
	}
	if wkttwo != "((0 0 0,0 1 0,1 1 0,0 0 0))" {
		t.Error("Expected ((0 0 0,0 1 0,1 1 0,0 0 0)) got ", wkttwo)
	}
}

// Test Point , Point Z, Point M, Point ZM
func TestPoint(t *testing.T) {
	var point string = "POINT(123.45 543.21)"
	pointtype, pointwkt := ParseGeometry(point)

	if strings.HasPrefix(point, "POINT") == false {
		t.Error("String was not prefixed with POINT")
	}
	if pointtype != "POINT" {
		t.Error("Expected POINT got", pointtype)
	}

	if len(pointwkt.Coordinates) != 1 {
		t.Error("Point does not have at least 1 pair of coordinates, has ", strconv.Itoa(len(pointwkt.Coordinates)))
	} else {
		if pointwkt.Coordinates[0].X != 123.45 {
			t.Error("Expected X to be 123.45 got ", pointwkt.Coordinates[0].X)
		}
		if pointwkt.Coordinates[0].Y != 543.21 {
			t.Error("Expected Y to be 543.21 got ", pointwkt.Coordinates[0].Y)
		}
	}
}

func TestPointZ(t *testing.T) {
	var pointz string = "POINT Z(123.45 543.21 65.6)"
	pointztype, pointzwkt := ParseGeometry(pointz)

	if strings.HasPrefix(pointz, "POINT Z") == false {
		t.Error("String was not prefixed with POINT Z")
	}
	if pointztype != "POINT Z" {
		t.Error("Expected POINT Z got", pointztype)
	}

	if len(pointzwkt.Coordinates) != 1 {
		t.Error("Point does not have at least 1 pair of coordinates, has ", strconv.Itoa(len(pointzwkt.Coordinates)))
	} else {
		if pointzwkt.Coordinates[0].X != 123.45 {
			t.Error("Expected X to be 123.45 got ", pointzwkt.Coordinates[0].X)
		}
		if pointzwkt.Coordinates[0].Y != 543.21 {
			t.Error("Expected Y to be 543.21 got ", pointzwkt.Coordinates[0].Y)
		}
		if pointzwkt.Coordinates[0].Z != 65.6 {
			t.Error("Expected Y to be 65.6 got ", pointzwkt.Coordinates[0].Z)
		}
	}

}

func TestPointZM(t *testing.T) {
	var pointzm string = "POINT ZM(123.45 543.21 65.6 100.0)"
	pointzmtype, pointzmwkt := ParseGeometry(pointzm)

	if strings.HasPrefix(pointzm, "POINT ZM") == false {
		t.Error("String was not prefixed with POINT ZM")
	}
	if pointzmtype != "POINT ZM" {
		t.Error("Expected POINT ZM got", pointzmtype)
	}

	if len(pointzmwkt.Coordinates) != 1 {
		t.Error("Point does not have at least 1 pair of coordinates, has ", strconv.Itoa(len(pointzmwkt.Coordinates)))
	} else {
		if pointzmwkt.Coordinates[0].X != 123.45 {
			t.Error("Expected X to be 123.45 got ", pointzmwkt.Coordinates[0].X)
		}
		if pointzmwkt.Coordinates[0].Y != 543.21 {
			t.Error("Expected Y to be 543.21 got ", pointzmwkt.Coordinates[0].Y)
		}
		if pointzmwkt.Coordinates[0].Z != 65.6 {
			t.Error("Expected Z to be 65.6 got ", pointzmwkt.Coordinates[0].Z)
		}
		if pointzmwkt.Coordinates[0].M != 100 {
			t.Error("Expected M to be 100.0 got ", pointzmwkt.Coordinates[0].M)
		}
	}

}

func TestLine(t *testing.T) {

	var line string = "LINESTRING (30 10, 10 30, 40 40)"
	linetype, linewkt := ParseGeometry(line)

	if strings.HasPrefix(line, "LINESTRING") == false {
		t.Error("String was not prefixed with LINESTRING")
	}
	if linetype != "LINESTRING" {
		t.Error("Expected LINESTRING got", linetype)
	}

	if len(linewkt.Coordinates) < 1 {
		t.Error("Line does not have at least 1 pair of coordinates, has ", strconv.Itoa(len(linewkt.Coordinates)))
	} else {
		if linewkt.Coordinates[0].X != 30.0 {
			t.Error("Expected X to be 30.0 got ", linewkt.Coordinates[0].X)
		}
		if linewkt.Coordinates[0].Y != 10.0 {
			t.Error("Expected Y to be 5.0 got ", linewkt.Coordinates[0].Y)
		}
	}

}

func TestLineZ(t *testing.T) {

	var line string = "LINESTRING Z (30 10 5, 10 30 5, 40 40 5)"
	linetype, linewkt := ParseGeometry(line)

	if strings.HasPrefix(line, "LINESTRING Z") == false {
		t.Error("String was not prefixed with LINESTRING Z")
	}
	if linetype != "LINESTRING Z" {
		t.Error("Expected LINESTRING got Z", linetype)
	}

	if len(linewkt.Coordinates) < 1 {
		t.Error("Line does not have at least 1 pair of coordinates, has ", strconv.Itoa(len(linewkt.Coordinates)))
	} else {
		if linewkt.Coordinates[0].X != 30.0 {
			t.Error("Expected X to be 30.0 got ", linewkt.Coordinates[0].X)
		}
		if linewkt.Coordinates[0].Y != 10.0 {
			t.Error("Expected Y to be 10.0 got ", linewkt.Coordinates[0].Y)
		}
		if linewkt.Coordinates[0].Z != 5 {
			t.Error("Expected X to be 5.0 got ", linewkt.Coordinates[0].Z)
		}
	}

}

func TestLineM(t *testing.T) {

	var line string = "LINESTRING M (30 10 10, 10 30 9, 40 40 8)"
	linetype, linewkt := ParseGeometry(line)

	if strings.HasPrefix(line, "LINESTRING M") == false {
		t.Error("String was not prefixed with LINESTRING M")
	}
	if linetype != "LINESTRING M" {
		t.Error("Expected LINESTRING M got ", linetype)
	}

	if len(linewkt.Coordinates) < 1 {
		t.Error("Line does not have at least 1 pair of coordinates, has ", strconv.Itoa(len(linewkt.Coordinates)))
	} else {
		if linewkt.Coordinates[0].X != 30.0 {
			t.Error("Expected X to be 30.0 got ", linewkt.Coordinates[0].X)
		}
		if linewkt.Coordinates[0].Y != 10.0 {
			t.Error("Expected Y to be 10.0 got ", linewkt.Coordinates[0].Y)
		}
		if linewkt.Coordinates[0].Z != 0.0 {
			t.Error("Expected Z to be 0.0 got ", linewkt.Coordinates[0].Z)
		}
		if linewkt.Coordinates[0].M != 10.0 {
			t.Error("Expected M to be 10.0 got ", linewkt.Coordinates[0].M)
		}
	}

}

func TestLineZM(t *testing.T) {

	var line string = "LINESTRING ZM (30 10 5 10, 10 30 5 9, 40 40 5 8)"
	linetype, linewkt := ParseGeometry(line)

	if strings.HasPrefix(line, "LINESTRING ZM") == false {
		t.Error("String was not prefixed with LINESTRING ZM")
	}
	if linetype != "LINESTRING ZM" {
		t.Error("Expected LINESTRING got ZM", linetype)
	}

	if len(linewkt.Coordinates) < 1 {
		t.Error("Line does not have at least 1 pair of coordinates, has ", strconv.Itoa(len(linewkt.Coordinates)))
	} else {
		if linewkt.Coordinates[0].X != 30.0 {
			t.Error("Expected X to be 30.0 got ", linewkt.Coordinates[0].X)
		}
		if linewkt.Coordinates[0].Y != 10.0 {
			t.Error("Expected Y to be 10.0 got ", linewkt.Coordinates[0].Y)
		}
		if linewkt.Coordinates[0].Z != 5.0 {
			t.Error("Expected X to be 5.0 got ", linewkt.Coordinates[0].Z)
		}
		if linewkt.Coordinates[0].M != 10.0 {
			t.Error("Expected X to be 5.0 got ", linewkt.Coordinates[0].M)
		}
	}

}
