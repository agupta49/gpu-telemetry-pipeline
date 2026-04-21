package api

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

type TelemetryPoint struct {
	GPUID      string    `json:"gpu_id"`
	MetricName string    `json:"metric_name"`
	Value      float64   `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
}

func NewRepo(dsn string) (*Repo, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Repo{db: db}, db.Ping()
}

func (r *Repo) Close() error {
	return r.db.Close()
}

func (r *Repo) ListGPUs(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT DISTINCT gpu_id FROM telemetry ORDER BY gpu_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func (r *Repo) GetTelemetry(ctx context.Context, gpuID string, start, end *time.Time) ([]TelemetryPoint, error) {
	query := `SELECT gpu_id, metric_name, value, timestamp FROM telemetry WHERE gpu_id = $1`
	args := []interface{}{gpuID}
	if start != nil {
		query += ` AND timestamp >= $2`
		args = append(args, *start)
	}
	if end != nil {
		if start != nil {
			query += ` AND timestamp <= $3`
		} else {
			query += ` AND timestamp <= $2`
		}
		args = append(args, *end)
	}
	query += ` ORDER BY timestamp ASC`
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var res []TelemetryPoint
	for rows.Next() {
		var p TelemetryPoint
		if err := rows.Scan(&p.GPUID, &p.MetricName, &p.Value, &p.Timestamp); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}
