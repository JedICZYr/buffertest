package one

import (
	"context"
	"gorm.io/gorm"
	"strings"
	"time"
	"os"

	"github.com/gulfcoastdevops/transformers/pkg/helpers"
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
}

var testPtr = false
var debugPtr = false
var company_id int64 = 0

// TableName overrides the table name used by Hardware to `hardware`
func (PcMetricCPU) TableName() string {
	return "pc_metrics_cpu"
}

// NewTimeScaleDBRepository -
func NewTimeScaleDBRepository(server string, userName string, passWord string, dbName string, port int) *TimescaleDBRepository {

	test := os.Getenv("PCMetricCPU_TEST")
	if test == "TRUE" {
		testPtr = true
	} else {
		testPtr = false
	}

	debug := os.Getenv("PCMetricCPU_DEBUG")
	if debug == "TRUE" {
		debugPtr = true
	} else {
		debugPtr = false
	}

	company := os.Getenv("COMPANY")
	if company == "" {
		log.Fatal("COMPANY not set")
	}

	var db *gorm.DB
	if !testPtr {
		// Create postgres database connection
		db = helpers.OpenConnection(server, port, userName, passWord, dbName)
		var err error
		company_id, err = helpers.GetCompany_ID(db, company)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Unable to Retrieve Company")
			panic(err)
		}
	}


	return &TimescaleDBRepository{db, server, userName, passWord, dbName, port}
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
func (t *TimescaleDBRepository) CreateCPUMetrics(ctx context.Context, data *PcMetricCPU) error {
	loc, _ := time.LoadLocation("UTC")
	t2 := time.Now().In(loc)
	ts := t2.UnixNano() 
	timestring := t2.Format(time.RFC3339)
	if testPtr {
		log.WithFields(log.Fields{
			"metric": data,
			"start": timestring,
			"timestamp": ts,
		}).Info("To Create")
		return nil
	} else if strings.Contains(data.FQDN, debug_host) {
		log.WithFields(log.Fields{
			"metric": data,
			"start": timestring,
			"timestamp": ts,
		}).Info("To Create")
	}

	// Create a continuous session with context
	tx := t.db.WithContext(ctx)

	data.Company_ID = uint(company_id)

	computer_id, err := helpers.GetComputer_ID(tx, company_id, data.FQDN)
	if err != nil {
		return err
	}

	data.Computer_ID = uint(computer_id)

	result := tx.Create(&data) // pass pointer of data to Create
	if result.Error != nil {
		log.WithFields(log.Fields{
			"company_id": company_id,
			"metric":     data,
		}).Warn("Error inserting new record")
		return result.Error
	} else {
		if debugPtr {
			log.WithFields(log.Fields{
				"company_id":   company_id,
				"metric":       data,
				"effectedRows": result.RowsAffected,
			}).Info("Inserted new record")
		} else if strings.Contains(data.FQDN, debug_host) {
			loc, _ := time.LoadLocation("UTC")
			t2 := time.Now().In(loc)
			timestring := t2.Format(time.RFC3339)
			ts := t2.UnixNano() 
			log.WithFields(log.Fields{
				"metric":       data,
				"end":			timestring,
				"timestamp":	ts,
				"effectedRows": result.RowsAffected,
			}).Info("Inserted new record")
		}
	}

	return nil
}
