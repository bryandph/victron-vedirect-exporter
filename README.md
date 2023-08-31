# Victron VEDirect Prometheus Exporter

This app implements a Prometheus Exporter for an attached VEDirect device. The exporter looks for an FDTI adapter connected to the machine and hosts metrics at `localhost:9101/metrics`. Only a few metrics are currently supported:

```
# HELP cs State of operation
# TYPE cs untyped
cs 0
# HELP err Error code
# TYPE err untyped
err 0
# HELP i_mA Main or channel 1 battery current
# TYPE i_mA gauge
i_mA -330
# HELP il_mA Load current
# TYPE il_mA gauge
il_mA 300
# HELP load Load output state (ON/OFF) 
# TYPE load untyped
load 1
# HELP ppv_W Panel power
# TYPE ppv_W gauge
ppv_W 0
# HELP v_mV Main or channel 1 (battery) voltage
# TYPE v_mV gauge
v_mV 13720
# HELP vpv_mV Panel voltage
# TYPE vpv_mV gauge
vpv_mV 13710
```