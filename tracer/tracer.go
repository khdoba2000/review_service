package tracer

import (
	"fmt"
	"monorepo/src/libs/log"
	"monorepo/src/libs/tracer"
	"monorepo/src/review_service/configs"

	"github.com/opentracing/opentracing-go"
	jexpvar "github.com/uber/jaeger-lib/metrics/expvar"
)

// Load ...
func Load(config *configs.Configuration, logger log.Factory) opentracing.Tracer {

	fmt.Println("tracer")
	metricsFactory := jexpvar.NewFactory(10) // 10 buckets for histograms
	tracer := tracer.Init(config.ServiceName, metricsFactory, logger)

	return tracer
}
