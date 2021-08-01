package timer

import (

	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var output chan<- []byte

// Start is to start the client
func Start(wg *sync.WaitGroup, ch chan<- []byte, topic string) error {
	output = ch

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		select {
		case s := <-sig:
			log.Println("Recevied Signal:", s)
			close(ch)
			wg.Done()
		default:
			output <- []byte("Message")
			time.Sleep(10 * time.Millisecond)     // sleep for 10 Milliseconds
		}
	}()
	return nil
}