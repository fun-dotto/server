package gtfs

import (
	"archive/zip"
	"bytes"
	"fmt"
	"path"
	"strings"
)

// Archive は GTFS ZIP を開いた読み取りハンドル。データの集約ではなく、
// 各ファイルを必要なときに CSV モデルとして読み出す入口。
type Archive struct {
	files map[string]*zip.File
}

// Open は GTFS ZIP のバイト列を読み取りハンドルにする。
func Open(zipData []byte) (*Archive, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	files := make(map[string]*zip.File, len(reader.File))
	for _, f := range reader.File {
		name := path.Base(f.Name)
		if strings.HasSuffix(name, "/") {
			continue
		}
		files[name] = f
	}
	return &Archive{files: files}, nil
}

func (a *Archive) Stops() ([]Stop, error)                 { return parseStops(a.files) }
func (a *Archive) Routes() ([]Route, error)               { return parseRoutes(a.files) }
func (a *Archive) Calendars() ([]Calendar, error)         { return parseCalendars(a.files) }
func (a *Archive) CalendarDates() ([]CalendarDate, error) { return parseCalendarDates(a.files) }
func (a *Archive) Trips() ([]Trip, error)                 { return parseTrips(a.files) }
func (a *Archive) StopTimes() ([]StopTime, error)         { return parseStopTimes(a.files) }

func (a *Archive) FareTables() ([]FareRule, []FareAttribute, error) {
	return parseFareTables(a.files)
}
