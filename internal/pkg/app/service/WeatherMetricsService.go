package service

import (
	"database/sql"
	"fmt"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app/utils"
	"github.com/palantir/stacktrace"
	"time"
)

func NewWeatherMetricsService(db *sql.DB) *WeatherMetricsService {
	return &WeatherMetricsService{db: db}
}

type MetricType struct {
	Name string
	Unit string
}

type MetricSource struct {
}

type WeatherMetric struct {
	Type        MetricType
	Source      *MetricSource
	CreatedDate time.Time
	Value       float64
}

type WeatherReport struct {
	Metrics []WeatherMetric
}

type MetricsFilter struct {
	Type       *string
	LoggerName *string
}

type WeatherMetricsService struct {
	db *sql.DB
}

func (s WeatherMetricsService) FindAllFiltered(filter MetricsFilter) ([]WeatherMetric, error) {
	query := `select s.type, s.unit, sr.value, sr.sensor_read_date
from sensor_reading sr
         left join sensor s on sr.sensor_id = s.id
         left join logger l on s.logger_id = l.id
where l.name = ?
  and s.type = ?
order by sr.created_date desc, sr.sensor_id`

	if metrics, err := s.readWeatherMetricsRows(s.db.Query(query, filter.LoggerName, filter.Type)); err != nil {
		return nil, err
	} else {
		return metrics, nil
	}
}

func (s WeatherMetricsService) CurrentWeather() (*WeatherReport, error) {
	query := `SELECT s.type, s.unit, sr.value, sr.sensor_read_date
FROM sensor_reading sr
INNER JOIN (
    SELECT sensor_id, max(sensor_read_date) as maxdate
    FROM sensor_reading
    GROUP BY sensor_id
) AS maxdates ON (sr.sensor_id = maxdates.sensor_id) AND (sr.sensor_read_date = maxdate)
left join sensor s on sr.sensor_id = s.id
left join logger l on s.logger_id = l.id
where l.name='met' or s.name in ('vann-temp')`

	if metrics, err := s.readWeatherMetricsRows(s.db.Query(query)); err != nil {
		println(fmt.Sprintf("%v", err))
		return nil, err
	} else {
		weatherReport := WeatherReport{Metrics: metrics}
		return &weatherReport, nil
	}
}

func (s WeatherMetricsService) readWeatherMetricsRows(rows *sql.Rows, err error) ([]WeatherMetric, error) {
	if err != nil {
		return nil, stacktrace.Propagate(err, "error when executing sql")
	}
	defer utils.CloseSilently(rows)

	var metrics []WeatherMetric
	metric := WeatherMetric{}
	for rows.Next() {
		if err := rows.Scan(&metric.Type.Name, &metric.Type.Unit, &metric.Value, &metric.CreatedDate); err != nil {
			return nil, stacktrace.Propagate(err, "error while reading result row")
		}
		metrics = append(metrics, metric)
	}
	if err := rows.Err(); err != nil {
		return nil, stacktrace.Propagate(err, "unexpected error while reading query result")
	}
	return metrics, nil
}
