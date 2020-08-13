package service

import (
	"database/sql"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewWeatherMetricsService(db *sql.DB) *WeatherMetricsService {
	return &WeatherMetricsService{db: db}
}

type WeatherMetric struct {
	CreatedDate time.Time
	Value       float32
}

type MetricsFilter struct {
	Type       *string
	LoggerName *string
}

type WeatherMetricsService struct {
	db *sql.DB
}

func (s WeatherMetricsService) FindAllFiltered(filter MetricsFilter) ([]WeatherMetric, error) {
	rows, err := s.db.Query(`select sr.sensor_read_date, sr.value
from sensor_reading sr
         left join sensor s on sr.sensor_id = s.id
         left join logger l on s.logger_id = l.id
where l.name = ?
  and s.type = ?
order by sr.created_date desc, sr.sensor_id`, filter.LoggerName, filter.Type)

	if err != nil {
		return nil, stacktrace.Propagate(err, "unable to load sensor readings from database")
	}
	defer rows.Close()

	var metrics []WeatherMetric
	metric := WeatherMetric{}
	for rows.Next() {
		err := rows.Scan(&metric.CreatedDate, &metric.Value)
		println(metric.Value)
		if err != nil {
			log.Fatal(err)
		}
		metrics = append(metrics, metric)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return metrics, nil
}
