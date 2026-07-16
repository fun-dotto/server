package service

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/fun-dotto/server/internal/modules/bus/gtfs"
)

// FetchArchive は URL から GTFS ZIP を取得し、読み取りハンドルを返す。
func (s *ScheduleService) FetchArchive(ctx context.Context, url string) (*gtfs.Archive, error) {
	zipData, err := s.download(ctx, url)
	if err != nil {
		return nil, err
	}
	archive, err := gtfs.Open(zipData)
	if err != nil {
		return nil, fmt.Errorf("open gtfs archive: %w", err)
	}
	return archive, nil
}

func (s *ScheduleService) download(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download zip: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download zip: unexpected status %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read zip body: %w", err)
	}
	return data, nil
}
