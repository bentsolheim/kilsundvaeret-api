package app

import (
	"encoding/json"
	"fmt"
	"github.com/Depado/ginprom"
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Debug struct {
	SignalStrength   string
	TimeSpent        int32
	Iteration        int32
	Errors           int32
	MillisSinceStart int64
	Battery          BatteryLevel
}

type BatteryLevel struct {
	AnalogReading int32
	Voltage       float32
	Level         int32
}

type DebugResponse struct {
	Items Debug
}

func CreateGinEngine(config AppConfig) *gin.Engine {
	r := gin.Default()
	p := ginprom.New(
		ginprom.Namespace(""),
		ginprom.Subsystem(""),
		ginprom.Engine(r),
		ginprom.Path("/api/metrics"),
	)
	r.Use(p.Instrument())
	p.AddCustomGauge("logger_battery_analog", "The analog reading of the battery level", []string{"loggerId"})
	p.AddCustomGauge("logger_battery_voltage", "The analog reading of the battery level converted to a voltage level", []string{"loggerId"})
	p.AddCustomGauge("logger_battery_level", "The analog reading of the battery level converted to a percentage", []string{"loggerId"})
	p.AddCustomGauge("logger_signal_strength", "The signal strength of the GSM connection on the last upload", []string{"loggerId"})
	p.AddCustomGauge("logger_time_spent_millis", "The number of milli seconds spent on the last iteration", []string{"loggerId"})
	errorCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "logger_processing_errors", Help: "A counter of all the times sending data has failed",
	}, []string{"loggerId"})
	prometheus.MustRegister(errorCounter)
	go func() {
		for {
			if err := updateMetrics(config, p, errorCounter); err != nil {
				log.Debug(err)
			}
			time.Sleep(60 * time.Second)
		}
	}()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/current-temp", func(c *gin.Context) {

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/logger/bua/readings", config.DataLoggerUrl))
			if err != nil {
				c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
				return
			}
			c.String(http.StatusOK, string(body))
		})
		v1.GET("/current-debug", func(c *gin.Context) {

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/logger/bua/debug", config.DataLoggerUrl))
			if err != nil {
				c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
				return
			}
			response := DebugResponse{}
			json.Unmarshal(body, &response)

			c.String(http.StatusOK, string(body))
		})
	}

	http.Handle("/metrics", promhttp.Handler())
	return r
}

func updateMetrics(config AppConfig, p *ginprom.Prometheus, errorCounter *prometheus.CounterVec) error {
	debugUrl := fmt.Sprintf("%s/api/v1/logger/bua/debug", config.DataLoggerUrl)
	resp, err := http.Get(debugUrl)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed loading debug data from [%s]", debugUrl))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading debug data body")
	}
	response := DebugResponse{}
	_ = json.Unmarshal(body, &response)

	_ = p.SetGaugeValue("logger_battery_analog", []string{"bua"}, float64(response.Items.Battery.AnalogReading))
	_ = p.SetGaugeValue("logger_battery_voltage", []string{"bua"}, float64(response.Items.Battery.Voltage))
	_ = p.SetGaugeValue("logger_battery_level", []string{"bua"}, float64(response.Items.Battery.Level))
	signalStrength, _ := strconv.ParseFloat(response.Items.SignalStrength, 64)
	_ = p.SetGaugeValue("logger_signal_strength", []string{"bua"}, signalStrength)
	_ = p.SetGaugeValue("logger_time_spent_millis", []string{"bua"}, float64(response.Items.TimeSpent))

	errorCounterBua := errorCounter.WithLabelValues("bua")
	value := readCounter(errorCounterBua)
	toAdd := float64(response.Items.Errors) - value
	errorCounterBua.Add(toAdd)
	return nil
}

func readCounter(m prometheus.Counter) float64 {
	pb := &dto.Metric{}
	_ = m.Write(pb)
	return pb.GetCounter().GetValue()
}
