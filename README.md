# FIO exporter

Simple Prometheus exporter written in Golang. Periodically executes an FIO benchmark and exports disk IO read and write bandwidth.

Example:

```bash
$ docker run --rm -it -p 9997:9997 neoaggelos/fio-exporter &
$ curl localhost:9997/metrics
# HELP fio_benchmark_duration_ms Duration of last successful benchmark (in ms)
# TYPE fio_benchmark_duration_ms gauge
fio_benchmark_duration_ms 50613
# HELP fio_benchmark_success 1 if last benchmark was successful, 0 otherwise
# TYPE fio_benchmark_success gauge
fio_benchmark_success 1
# HELP fio_read_bandwidth_kbps Read bandwidth measured by FIO (in KiB/s)
# TYPE fio_read_bandwidth_kbps gauge
fio_read_bandwidth_kbps 15715
# HELP fio_write_bandwidth_kbps Write bandwidth measured by FIO (in KiB/s)
# TYPE fio_write_bandwidth_kbps gauge
fio_write_bandwidth_kbps 6745
```

The fio benchmark command is currently hard-coded and executed every 30 minutes.
