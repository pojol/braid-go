package jaegertracing

import (
	"errors"
	"fmt"
	"io"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pojol/braid-go/module/tracer"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
	"github.com/uber/jaeger-lib/metrics"
)

const (
	// Name module name
	Name = "JaegerTracing"
)

var (
	// ErrFactoryNotExist factory not exist
	ErrFactoryNotExist = errors.New("factory not exist")
)

type jaegerTracingBuilder struct {
	factory map[string]tracer.SpanFactory
	opts    []interface{}
}

func newJaegerTracingBuilder() tracer.Builder {
	jtb := &jaegerTracingBuilder{
		factory: make(map[string]tracer.SpanFactory),
	}

	jtb.factory[EchoSpan] = createEchoTraceSpan()
	jtb.factory[RedisSpan] = createRedisSpanFactory()

	return jtb
}

func (jtb *jaegerTracingBuilder) Name() string {
	return Name
}

func (jtb *jaegerTracingBuilder) AddOption(opt interface{}) {
	jtb.opts = append(jtb.opts, opt)
}

func (jtb *jaegerTracingBuilder) AddFactory(strategy string, factory tracer.SpanFactory) {

	if _, ok := jtb.factory[strategy]; !ok {
		jtb.factory[strategy] = factory
	}

}

func newTransport(rc *jaegerCfg.ReporterConfig) (jaeger.Transport, error) {
	switch {
	case rc.CollectorEndpoint != "":
		httpOptions := []transport.HTTPOption{transport.HTTPBatchSize(1), transport.HTTPHeaders(rc.HTTPHeaders)}
		if rc.User != "" && rc.Password != "" {
			httpOptions = append(httpOptions, transport.HTTPBasicAuth(rc.User, rc.Password))
		}
		return transport.NewHTTPTransport(rc.CollectorEndpoint, httpOptions...), nil
	default:
		return jaeger.NewUDPTransport(rc.LocalAgentHostPort, 0)
	}
}

func (jtb *jaegerTracingBuilder) Build(name string) (tracer.ITracer, error) {

	p := Parm{
		Probabilistic: 1,
		SlowRequest:   time.Millisecond * 200,
		SlowSpan:      time.Millisecond * 50,
	}

	for _, opt := range jtb.opts {
		opt.(Option)(&p)
	}

	jcfg := jaegerCfg.Configuration{
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:           true,
			CollectorEndpoint:  p.CollectorEndpoint, //with http
			LocalAgentHostPort: p.LocalAgentHostPort,
		},
		ServiceName: name,
	}

	jt := &jaegerTracing{
		parm:        p,
		serviceName: name,
		jcfg:        jcfg,
		factory:     make(map[string]tracer.SpanFactory, len(jtb.factory)),
	}

	return jt, nil
}

func (jt *jaegerTracing) Init() error {
	sender, err := newTransport(jt.jcfg.Reporter)
	if err != nil {
		return fmt.Errorf("%v Dependency check error %v [%v]", jt.serviceName, "jaegertracing", err.Error())
	}

	r := jaegerCfg.Reporter(NewSlowReporter(sender, nil, jt.parm.Probabilistic))
	m := jaegerCfg.Metrics(metrics.NullFactory)

	jtracing, closer, err := jt.jcfg.NewTracer(r, m)
	if err != nil {
		return fmt.Errorf("%v Dependency check error %v [%v]", jt.serviceName, "jaegertracing", err.Error())
	}

	jt.tracing = jtracing
	jt.closer = closer

	for k, v := range jt.factory {
		jt.factory[k] = v
	}

	return nil
}

type jaegerTracing struct {
	parm        Parm
	serviceName string
	jcfg        jaegerCfg.Configuration

	closer  io.Closer
	tracing opentracing.Tracer

	factory map[string]tracer.SpanFactory
}

func (jt *jaegerTracing) Run() {

}

func (jt *jaegerTracing) GetSpan(strategy string) (tracer.ISpan, error) {

	spanfactory, ok := jt.factory[strategy]
	if !ok {
		return nil, ErrFactoryNotExist
	}

	span, err := spanfactory(jt.tracing)
	if err != nil {
		return nil, err
	}

	return span, nil
}

func (jt *jaegerTracing) GetTracing() interface{} {
	return jt.tracing
}

func (jt *jaegerTracing) Close() {
	jt.closer.Close()
}

func init() {
	tracer.Register(newJaegerTracingBuilder())
}
