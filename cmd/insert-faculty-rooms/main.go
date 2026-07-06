package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/repository"
	"github.com/fun-dotto/server/internal/modules/batch-jobs/service"
	"github.com/fun-dotto/server/internal/shared/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	csvPathByYear := map[int]string{}
	flag.Func("faculties", "年度と CSV パスのペア（例: --faculties 2026=faculties_2026.csv）。複数年度は本オプションを複数回指定する。CSV 必須カラム: email, room_name", func(value string) error {
		yearStr, path, ok := strings.Cut(value, "=")
		if !ok || path == "" {
			return fmt.Errorf("YEAR=PATH 形式で指定してください: %q", value)
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return fmt.Errorf("年度は整数で指定してください: %q", yearStr)
		}
		if _, exists := csvPathByYear[year]; exists {
			return fmt.Errorf("年度 %d が複数回指定されています", year)
		}
		csvPathByYear[year] = path
		return nil
	})
	flag.Parse()
	if len(csvPathByYear) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	years := make([]int, 0, len(csvPathByYear))
	for year := range csvPathByYear {
		years = append(years, year)
	}
	sort.Ints(years)

	var rows []domain.FacultyRoomImportRow
	for _, year := range years {
		path := csvPathByYear[year]
		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("Failed to open CSV %s: %v", path, err)
		}
		parsed, err := service.ParseFacultyRoomCSV(year, f)
		_ = f.Close()
		if err != nil {
			log.Fatalf("Failed to parse CSV %s: %v", path, err)
		}
		rows = append(rows, parsed...)
	}

	conn, err := db.ConnectWithConnectorIAMAuthN()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(conn); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	svc := service.NewFacultyRoomImportService(
		repository.NewFacultyRepository(conn),
		repository.NewRoomRepository(conn),
		repository.NewFacultyRoomRepository(conn),
	)

	summary, err := svc.Import(context.Background(), rows)
	if err != nil {
		log.Fatalf("Failed to import faculty rooms: %v (unmatched_emails=%v unmatched_rooms=%v)", err, summary.UnmatchedEmails, summary.UnmatchedRooms)
	}
	log.Printf("import summary: inserted=%d skipped_blank=%d", summary.Inserted, summary.SkippedBlank)
}
