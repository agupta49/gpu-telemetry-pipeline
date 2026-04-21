package collector

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/agupta49/gpu-telemetry-pipeline/pkg/pb"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(dsn string) (*Repo, error) {
	if dsn == "" {
		return nil, fmt.Errorf("empty dsn")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS telemetry (
		id SERIAL PRIMARY KEY,
		gpu_id TEXT NOT NULL,
		metric_name TEXT NOT NULL,
		value DOUBLE PRECISION NOT NULL,
		timestamp TIMESTAMPTZ NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_telemetry_gpu_time ON telemetry(gpu_id, timestamp);`)
	if err != nil {
		return nil, err
	}
	return &Repo{db: db}, nil
}

func (r *Repo) Close() error {
	return r.db.Close()
}

func (r *Repo) Insert(p *pb.TelemetryPoint) error {
	_, err := r.db.Exec(`INSERT INTO telemetry (gpu_id, metric_name, value, timestamp) VALUES ($1, $2, $3, to_timestamp($4))`,
		p.GpuId, p.MetricName, p.Value, p.Timestamp)
	return err
}
