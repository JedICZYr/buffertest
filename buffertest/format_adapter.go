package one

import (
	"context"
	"strconv"
	"os"
	"strings"
	"time"


	log "github.com/sirupsen/logrus"
)

var debug_host string = ""

// FormatAdapter wraps the usecase interface
// with a logging adapter which can be swapped out
type FormatAdapter struct {
	usecase OneService
}

// Start start the FormatAdapter
func (f *FormatAdapter) Start(ch <-chan []byte, u OneService) error {
	h := os.Getenv("PCMetricCPU_HOST")

	f.usecase = u
	go func() {
		for {
			// prior to reading the next message we have to determine if 
			// the pool is full using select and send a persist message if it is.
			res, ok := <-ch
			if !ok {
				// Channel close
				return
			}

			// channel open new message
			m, err := helpers.ExtractMetric(string(res))
			if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Warn("Parsing Metric")
				continue
			}

			if h != "" {
				debug_host = h
				fqdn, ok := m.GetTag("host")
				if ok {
					if strings.Contains(fqdn, h) {
						loc, _ := time.LoadLocation("UTC")
						t2 := time.Now().In(loc)
						ts := t2.UnixNano() 
						timestring := t2.Format(time.RFC3339)
						log.WithFields(log.Fields{
							"metric": m,
							"time": timestring,
							"timestamp": ts,
						}).Info("Received Metric")
					}
				}

			}
			f.ProcessMetric(m)
		}
	}()
	return nil
}

// ProcessData - Process CPU Data
func (f *FormatAdapter) ProcessMetric(m telegraf.Metric) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	metric := &PcMetricCPU{}

	// FQDN
	if fqdn, ok := m.GetTag("host"); ok {
		metric.FQDN = fqdn
	}

	if i, ok := m.GetTag("instance"); ok {
		if d, err := strconv.Atoi(i); err == nil {
			metric.Instance = d
		} else {
			log.WithFields(log.Fields{
				"err":  err,
				"data": i,
			}).Warn("Failure converting Instance Atoi")
		}
	}

	if r, err := helpers.GetFloatField(m, "Percent_Idle_Time"); err == nil {
		metric.Utilization = r
	} else {
		log.WithFields(log.Fields{
			"err": err,
		}).Warn("Failure converting Percent_Idle_Time Float")
	}

	// Metric Timestamp
	t2 := m.Time()
	ts := t2.Format(time.RFC3339)

	timestring, err := helpers.ConvertTime(strings.Replace(ts, "\n", "", -1))
	if err != nil {
		return
	}
	metric.Time = timestring

	if err := f.usecase.CreateCPUMetrics(ctx, metric); err != nil {
		log.Print(err)
		return
	}
}
