package one

import (
	"context"
	"gorm.io/gorm"
	"time"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// InfluxDBRepository -
type TimescaleDBRepository struct {
	db       *gorm.DB
	server   string
	userName string
	passWord string
	dbName   string
	port     int
	pool []*DataStruct
	poolsize int
	poolindex int
	maxdelay uint
}

var testPtr = false
var debugPtr = false
var company_id int64 = 0


// NewTimeScaleDBRepository -
func NewTimeScaleDBRepository(server string, userName string, passWord string, dbName string, port int) *TimescaleDBRepository {

	var db *gorm.DB

	poolindex := 0
	maxdelay := uint(0)
	poolsize := 25
	poolstr := os.Getenv("BufferTest_PoolSize")
	if i, err := strconv.Atoi(poolstr); err == nil {
		poolsize = i
	}

	pool := make([]*DataStruct, poolsize)
	
	return &TimescaleDBRepository{db, server, userName, passWord, dbName, port, pool, poolsize, poolindex, maxdelay}
}

func (t *TimescaleDBRepository) CloseRepository() error {
	sqlDB, err := t.db.DB()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Unable to Retrieve DB")
	}
	err = sqlDB.Close()
	return err
}

//  CreateInstallerData Record
func (t *TimescaleDBRepository) CreateMetric(ctx context.Context, data *DataStruct) error {

	if t.poolindex < t.poolsize  {
		t.pool[t.poolindex] = data
		t.poolindex++
	}

	if t.poolindex == t.poolsize {
		loc, _ := time.LoadLocation("UTC")
		t2 := time.Now().In(loc)
		ts := t2.UnixNano() 
		timestring := t2.Format(time.RFC3339)
	
		log.WithFields(log.Fields{
			"pool": t.pool,
			"index": t.poolindex,
			"size": t.poolsize,
			"time": timestring,
			"timestamp": ts,
		}).Info("To Create")
		promserver.msgWritten.Inc()
		delay := ts - data.ReceiveTime.UnixNano()
		if delay > int64(t.maxdelay) {
			t.maxdelay = uint(delay)
			promserver.opsDelay.Set(float64(t.maxdelay))
		}
		for index := range t.pool {
			t.pool[index] = nil
		}
		t.poolindex = 0
	}
	return nil
}