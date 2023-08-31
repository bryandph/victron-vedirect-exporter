package main

import (
	"fmt"
	"net/http"

	vedirect_device "github.com/bryandph/victron-vedirect/vedirect-device"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial/enumerator"
)

func getFTDIPort() string {
	ports, err := enumerator.GetDetailedPortsList()
	if (err != nil) || (len(ports) == 0) {
		log.Fatal("No serial ports found!")
	}
	var portname string
	for _, port := range ports {
		if port.IsUSB {
			if port.VID == "0403" {
				portname = port.Name
			}
		}
	}
	if portname != "" {
		log.WithFields(log.Fields{"serialport": portname}).Debug("Using found FTDI device")
	}
	return portname
}

type vedirectCollector struct {
	vMetric    *prometheus.Desc
	iMetric    *prometheus.Desc
	ilMetric   *prometheus.Desc
	ppvMetric  *prometheus.Desc
	vpvMetric  *prometheus.Desc
	loadMetric *prometheus.Desc
	errMetric  *prometheus.Desc
	csMetric   *prometheus.Desc
	mpptMetric *prometheus.Desc
	device     *vedirect_device.VEDirectDevice
}

func newVedirectCollector() *vedirectCollector {
	dev, err := vedirect_device.NewVEDirectDevice(getFTDIPort())
	if err != nil {
		log.Fatal(err)
	}
	var block = dev.GetBlock()
	return &vedirectCollector{
		vMetric:    prometheus.NewDesc(fmt.Sprintf("v_%s", block.Label["V"].Unit), block.Label["V"].Description, nil, nil),
		iMetric:    prometheus.NewDesc(fmt.Sprintf("i_%s", block.Label["I"].Unit), block.Label["I"].Description, nil, nil),
		ilMetric:   prometheus.NewDesc(fmt.Sprintf("il_%s", block.Label["IL"].Unit), block.Label["IL"].Description, nil, nil),
		ppvMetric:  prometheus.NewDesc(fmt.Sprintf("ppv_%s", block.Label["PPV"].Unit), block.Label["PPV"].Description, nil, nil),
		vpvMetric:  prometheus.NewDesc(fmt.Sprintf("vpv_%s", block.Label["VPV"].Unit), block.Label["VPV"].Description, nil, nil),
		loadMetric: prometheus.NewDesc("load", block.Label["LOAD"].Description, nil, nil),
		errMetric:  prometheus.NewDesc("err", block.Label["ERR"].Description, nil, nil),
		csMetric:   prometheus.NewDesc("cs", block.Label["CS"].Description, nil, nil),
		mpptMetric: prometheus.NewDesc("mppt", block.Label["MPPT"].Description, nil, nil),
		device:     dev,
	}
}

func (collector *vedirectCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.vMetric
	ch <- collector.iMetric
	ch <- collector.ilMetric
	ch <- collector.ppvMetric
	ch <- collector.vpvMetric
	ch <- collector.loadMetric
	ch <- collector.errMetric
	ch <- collector.csMetric
	ch <- collector.mpptMetric
}

func (collector *vedirectCollector) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var block = collector.device.GetBlock()

	v := prometheus.MustNewConstMetric(collector.vMetric, prometheus.GaugeValue, block.Label["V"].AsFloat())
	i := prometheus.MustNewConstMetric(collector.iMetric, prometheus.GaugeValue, block.Label["I"].AsFloat())
	il := prometheus.MustNewConstMetric(collector.ilMetric, prometheus.GaugeValue, block.Label["IL"].AsFloat())
	ppv := prometheus.MustNewConstMetric(collector.ppvMetric, prometheus.GaugeValue, block.Label["PPV"].AsFloat())
	vpv := prometheus.MustNewConstMetric(collector.vpvMetric, prometheus.GaugeValue, block.Label["VPV"].AsFloat())
	load := prometheus.MustNewConstMetric(collector.loadMetric, prometheus.UntypedValue, block.Label["LOAD"].AsFloat())
	err := prometheus.MustNewConstMetric(collector.errMetric, prometheus.UntypedValue, block.Label["ERR"].AsFloat())
	cs := prometheus.MustNewConstMetric(collector.csMetric, prometheus.UntypedValue, block.Label["CS"].AsFloat())
	mppt := prometheus.MustNewConstMetric(collector.mpptMetric, prometheus.UntypedValue, block.Label["MPPT"].AsFloat())
	ch <- v
	ch <- i
	ch <- il
	ch <- ppv
	ch <- vpv
	ch <- load
	ch <- err
	ch <- cs
	ch <- mppt
}

func main() {
	log.SetLevel(log.WarnLevel)
	vedirect := newVedirectCollector()
	r := prometheus.NewRegistry()
	r.MustRegister(vedirect)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(":9101", nil))
}
