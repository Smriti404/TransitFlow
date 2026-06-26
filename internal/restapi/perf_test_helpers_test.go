//go:build perftest

package restapi

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"transitflow/internal/app"
	"transitflow/internal/appconf"
	"transitflow/internal/clock"
	"transitflow/internal/gtfs"
	"transitflow/internal/models"
)

// createLargeAgencyApi builds a RestAPI backed by the TriMet GTFS feed in testdata/perf/trimet.zip.
// Used only by perftest-tagged tests. Run scripts/download-perf-data.sh first.
func createLargeAgencyApi(tb testing.TB) *RestAPI {
	tb.Helper()
	zipPath := models.GetFixturePath(tb, "perf/trimet.zip")
	if _, err := os.Stat(zipPath); err != nil {
		tb.Fatalf("perf GTFS not found at %s: run scripts/download-perf-data.sh first: %v", zipPath, err)
	}
	dbPath := filepath.Join(filepath.Dir(zipPath), "trimet-perf.db")
	cfg := gtfs.Config{
		GtfsURL:      zipPath,
		GTFSDataPath: dbPath,
	}
	mgr, err := gtfs.InitGTFSManager(cfg)
	if err != nil {
		tb.Fatalf("init GTFS manager for perf: %v", err)
	}
	application := &app.Application{
		Config: appconf.Config{
			Env:       appconf.EnvFlagToEnvironment("test"),
			ApiKeys:   []string{"TEST"},
			RateLimit: 100,
		},
		GtfsConfig:  cfg,
		GtfsManager: mgr,
		Clock:       clock.RealClock{},
	}
	api := NewRestAPI(application)
	api.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	return api
}
