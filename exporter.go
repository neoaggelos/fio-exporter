package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr       = ":9997"
	fioCommand = "fio --randrepeat=1 --ioengine=libaio --direct=1 --gtod_reduce=1 --name=test --filename=/tmp/fio-exporter.test --bs=4k --iodepth=32 --size=1G --readwrite=randrw --rwmixread=70"
	interval   = 30 * time.Minute
)

var (
	promRegistry = prometheus.NewRegistry()
	fioReadBW    = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fio_read_bandwidth_kbps",
		Help: "Read bandwidth measured by FIO (in KiB/s)",
	})
	fioWriteBW = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fio_write_bandwidth_kbps",
		Help: "Write bandwidth measured by FIO (in KiB/s)",
	})
	fioBenchmarkDurationMS = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fio_benchmark_duration_ms",
		Help: "Duration of last successful benchmark (in ms)",
	})
	fioBenchmarkSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fio_benchmark_success",
		Help: "1 if last benchmark was successful, 0 otherwise",
	})
)

func init() {
	promRegistry.MustRegister(
		fioReadBW,
		fioWriteBW,
		fioBenchmarkDurationMS,
		fioBenchmarkSuccess,
	)
}

func handleError(err error) {
	log.Printf("Error: %s\n", err)
	fioBenchmarkSuccess.Set(0)
}

func handleOK(readBW, writeBW float64, duration time.Duration) {
	fioBenchmarkSuccess.Set(1)
	fioReadBW.Set(readBW)
	fioWriteBW.Set(writeBW)
	fioBenchmarkDurationMS.Set(float64(duration / time.Millisecond))
}

func main() {
	go func() {
		ch := make(chan struct{}, 1)
		ch <- struct{}{}
		for {
			<-ch
			time.AfterFunc(interval, func() { ch <- struct{}{} })

			log.Printf("Running fio: %s", fioCommand)

			start := time.Now()
			cmd := fmt.Sprintf("%s --minimal", fioCommand)
			cmdParts := strings.Split(cmd, " ")
			output, err := exec.Command(cmdParts[0], cmdParts[1:]...).Output()
			if err != nil {
				handleError(err)
				continue
			}

			fioBenchmarkSuccess.Set(1)
			s := strings.TrimSuffix(string(output), "\n")
			parts := strings.Split(s, ";")

			d := time.Now().Sub(start)
			log.Printf("Benchmark finished after %v: %v", d, parts)

			readBW, err := strconv.ParseFloat(parts[6], 64)
			if err != nil {
				handleError(err)
				continue
			}

			writeBW, err := strconv.ParseFloat(parts[47], 64)
			if err != nil {
				handleError(err)
				continue
			}

			handleOK(readBW, writeBW, d)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(
		promRegistry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
