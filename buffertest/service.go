package one

import (
	"context"
	"os"
	"strconv"

	"github.com/aws/aws-xray-sdk-go/xray"
	"go.uber.org/zap"

	log "github.com/sirupsen/logrus"
)

var promserver *PrometheusServer

// CMDBPCService is the top level signature of this service
type OneService interface {
	CreateMetric(ctx context.Context, data *DataStruct) error
	CloseRepository() error
}

// Init sets up an instance of this domains
// usecase, pre-configured with the dependencies.
func Init(integration bool) (OneService, error) {

	if !integration {
		err := xray.Configure(xray.Config{LogLevel: "trace"})
		if err != nil {
			log.Fatal("Error configuring x ray. ", err)
			return nil, err
		}
		// xray.AWS(ddb.Client)
	}

	logger, _ := zap.NewProduction()

	t_server, ok := os.LookupEnv("TIMESCALE_SERVER")
	if !ok {
		t_server = ""
	}

	t_user, ok := os.LookupEnv("TIMESCALE_USER")
	if !ok {
		t_user = ""
	}

	t_pass, ok := os.LookupEnv("TIMESCALE_PASS")
	if !ok {
		t_pass = ""
	}

	t_db, ok := os.LookupEnv("TIMESCALE_DB")
	if !ok {
		t_db = ""
	}

	var t_port int = 0
	env_t_port, ok := os.LookupEnv("TIMESCALE_PORT")
	if ok {
		t_port, _ = strconv.Atoi(env_t_port)
	}

	promserver= NewPrometheusServer()
	repository := NewTimeScaleDBRepository(t_server, t_user, t_pass, t_db, t_port)

	usecase := &LoggerAdapter{
		Logger:  logger,
		Usecase: &Usecase{Repository: repository},
	}
	return usecase, nil
}
