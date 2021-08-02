package one

import (
	"context"
	"time"


	log "github.com/sirupsen/logrus"
)

var debug_host string = ""

// MetricData includes time to perform analysis
// of delay between message receive and queue/db write

type MetricData struct {
	receive time.Time
	metric string
}

// FormatAdapter wraps the usecase interface
// with a logging adapter which can be swapped out
type FormatAdapter struct {
	usecase OneService
}

// Start start the FormatAdapter
func (f *FormatAdapter) Start(ch <-chan []byte, u OneService) error {

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
			promserver.msgRecevied.Inc()

			loc, _ := time.LoadLocation("UTC")
			t := time.Now().In(loc)

			metric := &MetricData{}
			metric.metric= string(res)
			metric.receive = t


			f.ProcessMetric(metric)
		}
	}()
	return nil
}

// ProcessData - Process CPU Data
func (f *FormatAdapter) ProcessMetric(d *MetricData) {

	m := &DataStruct{}
	m.Data = d.metric
	m.ReceiveTime = d.receive

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := f.usecase.CreateMetric(ctx, m); err != nil {
		log.Print(err)
		return
	}
}
