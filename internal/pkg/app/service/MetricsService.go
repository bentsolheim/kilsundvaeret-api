package service

import (
	"fmt"
	"github.com/bentsolheim/go-app-utils/utils"
	"github.com/palantir/stacktrace"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"log"
	"strconv"
	"time"
)

func NewMetricsService(dataLoggerUrl string, registerer prometheus.Registerer) MetricsService {
	newGauge := func(name string, help string) *prometheus.GaugeVec {
		gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, []string{"loggerId"})
		registerer.MustRegister(gauge)
		return gauge
	}
	newCounter := func(name string, help string) *prometheus.CounterVec {
		counter := prometheus.NewCounterVec(prometheus.CounterOpts{Name: name, Help: help}, []string{"loggerId"})
		registerer.MustRegister(counter)
		return counter
	}
	s := MetricsService{
		dataLoggerUrl:          dataLoggerUrl,
		batteryAnalog:          newGauge("logger_battery_analog", "The analog reading of the battery level"),
		batteryVoltage:         newGauge("logger_battery_voltage", "The analog reading of the battery level converted to a voltage level"),
		batteryLevel:           newGauge("logger_battery_level", "The analog reading of the battery level converted to a percentage"),
		signalStrength:         newGauge("logger_signal_strength", "The signal strength of the GSM connection on the last upload"),
		timeSpent:              newGauge("logger_time_spent_millis", "The number of milli seconds spent on the last iteration"),
		connectionErrorCounter: newCounter("logger_connection_errors", "A counter of all the times sending data has failed"),
		sensorErrorCounter:     newCounter("logger_sensor_errors", "A counter of all the times reading sensors has failed"),
	}
	return s
}

type MetricsService struct {
	dataLoggerUrl          string
	batteryAnalog          *prometheus.GaugeVec
	batteryVoltage         *prometheus.GaugeVec
	batteryLevel           *prometheus.GaugeVec
	signalStrength         *prometheus.GaugeVec
	timeSpent              *prometheus.GaugeVec
	connectionErrorCounter *prometheus.CounterVec
	sensorErrorCounter     *prometheus.CounterVec
}

func (s MetricsService) UpdateMetrics(loggerId string) error {
	url := fmt.Sprintf("%s/api/v1/logger/%s/debug", s.dataLoggerUrl, loggerId)
	response := DebugResponse{}
	if err := utils.HttpGetJson(url, response); err != nil {
		return stacktrace.Propagate(err, "error getting debug data from data receiver")
	}

	gauge := func(gv *prometheus.GaugeVec) prometheus.Gauge { return gv.WithLabelValues(loggerId) }
	setCounterValue := func(cv *prometheus.CounterVec, newValue float64) {
		counter := cv.WithLabelValues(loggerId)
		diff := newValue - readCounter(counter)
		if diff < 0 {
			cv.Reset()
			counter = cv.WithLabelValues(loggerId)
			diff = newValue
		}
		counter.Add(diff)
	}

	debug := response.Items
	gauge(s.batteryAnalog).Set(float64(debug.Battery.AnalogReading))
	gauge(s.batteryVoltage).Set(float64(debug.Battery.Voltage))
	gauge(s.batteryLevel).Set(float64(debug.Battery.Level))
	signalStrength, _ := strconv.ParseFloat(debug.SignalStrength, 64)
	gauge(s.signalStrength).Set(signalStrength)
	gauge(s.timeSpent).Set(float64(debug.TimeSpent))

	setCounterValue(s.connectionErrorCounter, float64(debug.ConnectionErrors))
	setCounterValue(s.sensorErrorCounter, float64(debug.SensorErrors))

	return nil
}

func (s MetricsService) UpdateMetricsForever(loggerId string, sleepSeconds time.Duration) {
	for {
		if err := s.UpdateMetrics(loggerId); err != nil {
			log.Print(err)
		}
		time.Sleep(sleepSeconds * time.Second)
	}
}

func readCounter(m prometheus.Counter) float64 {
	pb := &dto.Metric{}
	_ = m.Write(pb)
	return pb.GetCounter().GetValue()
}
