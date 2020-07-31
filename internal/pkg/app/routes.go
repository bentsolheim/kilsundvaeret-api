package app

import (
	"encoding/json"
	"fmt"
	"github.com/Depado/ginprom"
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	go func() {
		for {
			if updateMetrics(config, p) {
				return
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

func updateMetrics(config AppConfig, p *ginprom.Prometheus) bool {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/logger/bua/debug", config.DataLoggerUrl))
	if err != nil {
		return true
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return true
	}
	response := DebugResponse{}
	_ = json.Unmarshal(body, &response)

	_ = p.SetGaugeValue("logger_battery_analog", []string{"bua"}, float64(response.Items.Battery.AnalogReading))
	_ = p.SetGaugeValue("logger_battery_voltage", []string{"bua"}, float64(response.Items.Battery.Voltage))
	_ = p.SetGaugeValue("logger_battery_level", []string{"bua"}, float64(response.Items.Battery.Level))
	signalStrength, _ := strconv.ParseFloat(response.Items.SignalStrength, 64)
	_ = p.SetGaugeValue("logger_signal_strength", []string{"bua"}, signalStrength)
	_ = p.SetGaugeValue("logger_time_spent_millis", []string{"bua"}, float64(response.Items.TimeSpent))
	return false
}
