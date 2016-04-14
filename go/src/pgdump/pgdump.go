package main

import (
    "fmt"
    "encoding/json"
    "os"
    "io"
    "log"
    "net/http"
    "database/sql"
    "strings"
    "time"
    "strconv"
    "regexp"
    "github.com/lib/pq"
    "github.com/rs/cors"
)

type Coordinate struct {
    X   float64
    Y   float64
    Z   float64
}

type CoordinateSet struct {
    Coordinates []Coordinate  //(0 0 0, 0 1 0, 1 1 0, 1 0 0, 0 0 0)
}

// Define WTK
type WTK struct {
    WTKType     string
    Geometry    []CoordinateSet
}

type ReturnJSON struct {
    WTKGeoms  []WTK
    Elapsed   time.Duration
}


func handler(w http.ResponseWriter, r *http.Request) {

    jsonenc := json.NewEncoder(w)
    w.Header().Set("Content-Type", "application/json")
    f, err := os.OpenFile("pgdump_errorlog.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)

    log.Print("Couldn't open file")
    ///log.SetOutput(f)
    defer f.Close()

    // Timing
    start := time.Now()

    // Postgres Credentials
    const (
        DB_USER     = "postgres"
        DB_PASSWORD = "james"
        DB_PORT     = "1337"
        DB_NAME     = "Lacuna"
    )

    // Postgres Connect
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
                           DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
    db, err := sql.Open("postgres", dbinfo)
    if err != nil {
        handleError(w, err.Error())
    }
    defer db.Close()

    table := r.FormValue("table")
    feature := r.FormValue("id")
    if table != "" {

        //Postgres Query
        var (
        	id int
        	geom string
        )

        table := pq.QuoteIdentifier(table)
        identifier := pq.QuoteIdentifier("ID")
        rows, err := db.Query( fmt.Sprintf("SELECT %s, ST_AsText(geom) FROM %s WHERE %s = $1", identifier, table, identifier), feature)
        if err != nil {
            handleError(w, err.Error())
            return
        }
        defer rows.Close()
        for rows.Next() {
        	err := rows.Scan(&id, &geom)
            if err != nil {
                handleError(w, err.Error())
                return
            }
        }
        err = rows.Err()
        if err != nil {
            handleError(w, err.Error() )
            return
        }
        returnjson := ReturnJSON{}

        // Maniplate Strings
        returngeom := strings.Replace(geom, "1.#QNAN", "", -1)

        if isGeometryCollection(returngeom)  {

            // Geometry Collection e.g. - GEOMETRYCOLLECTION(POINT(4 6),LINESTRING(4 6,7 10))
            s          := strings.Index(returngeom, "(") + 1
            e          := strings.LastIndex(returngeom, ")") - 1
            geomsstr   := returngeom[s:e] // From ( to )
            geomsre    := regexp.MustCompile(`GEOMETRY|POINT|MULTIPOINT|LINESTRING|MULTILINESTRING|COMPOUNDCURVE|
                                              MULTIPOLYGON|TRIANGLE|CIRCULARSTRING|CURVE|MULTICURVE|POLYGON|POLYGON Z|
                                              CURVEPOLYGON|SURFACE|MULTISURFACE|POLYHEDRALSURFACE Z|TIN|TIN Z`)
            geometries := geomsre.Split(geomsstr, -1)
            wkttypes    := geomsre.FindAllString(geomsstr, -1)

            returnjson  = WKTtoJSON(geometries, wkttypes, start)

        } else {

            s := strings.Index(returngeom, "(")
            if s != -1 {
                wkttype    := []string{strings.TrimSpace(returngeom[:s])}
                geometry   := []string{returngeom[s:]}
                returnjson  = WKTtoJSON(geometry, wkttype, start)
            } else {
                handleError(w, "No ID by that number")
                return
            }

        }

        jsonenc.Encode(returnjson)
    }
}


func WKTtoJSON (geometries []string, wkttypes []string, start time.Time ) ReturnJSON {

    wkts       := []WTK{}
    coordslice := []Coordinate{}
    geos       := []CoordinateSet{}

    for i, g := range geometries {
        //log.Print(strconv.FormatInt(int64(len(geometries)), 10))
        if i <= len(wkttypes) - 1 {
            g = removeTrails(g)
            log.Print()
            subgeoms := strings.Split(g, ")),((")
            for _, subgeom := range subgeoms {
                subgeom = removeTrails(subgeom)
                subgeom = removeBrackets(subgeom)
                coords  := strings.Split(subgeom, ",")
                if len(coords) > 0 {
                    for _, coord := range coords {
                        xyz        := strings.Split(coord, " ")
                        if len(xyz) > 2 {
                            x, err     := strconv.ParseFloat(xyz[0], 64); if err != nil { x = 0.0 }
                            y, err     := strconv.ParseFloat(xyz[1], 64); if err != nil { y = 0.0 }
                            z, err     := strconv.ParseFloat(xyz[2], 64); if err != nil { z = 0.0 }
                            coordslice = append(coordslice, Coordinate{ X: x, Y: y, Z: z}) // CoordinateSet
                        }
                    }
                }
                if len(coordslice) > 0 {
                    geos = append(geos, CoordinateSet{Coordinates : coordslice})
                }
            }
            //log.Print(geos)
            if len(geos) > 0 {
                wkts     = append(wkts, WTK{ Geometry: geos, WTKType: wkttypes[i] })
            }
        }
    }

    return ReturnJSON{WTKGeoms: wkts, Elapsed: time.Since(start)/1000000.0 }
}

func removeBrackets(str string) string {

    str = strings.Replace(str, "(", "", -1)
    str  = strings.Replace(str, ")", "", -1)
    return str
}

func removeTrails(str string) string {
    if strings.HasSuffix(str, ",") {
        restr := strings.TrimSpace(str[: strings.LastIndex(str, ",")])
        return restr
    }
    return str
}

func isGeometryCollection(inputstr string) bool {
    return strings.Contains(inputstr,"GEOMETRYCOLLECTION")
}

func handleError(w http.ResponseWriter, err string) {
    type APIError struct {
        Error string
    }
    re, _ := json.Marshal(APIError{Error: err})
    io.WriteString(w, string(re))
}

func main() {
    c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost"},
	})

    httphandler := http.HandlerFunc(handler)
    http.ListenAndServe(":8080", c.Handler(httphandler))
}
