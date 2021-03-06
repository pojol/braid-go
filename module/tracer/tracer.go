package tracer

import (
	"strings"

	"github.com/pojol/braid-go/module"
)

// SpanFactory span 工厂
type SpanFactory func(interface{}) (ISpan, error)

// Builder tracer build
type Builder interface {
	Build(name string) (ITracer, error)
	Name() string
	AddOption(opt interface{})
	AddFactory(strategy string, factory SpanFactory)
}

// ISpan span interface
type ISpan interface {
	Begin(ctx interface{})
	End(ctx interface{})
}

// ITracer tracer interface
type ITracer interface {
	module.IModule

	GetSpan(strategy string) (ISpan, error)

	GetTracing() interface{}
}

var (
	m = make(map[string]Builder)
)

// Register 注册tracer
func Register(b Builder) {
	m[strings.ToLower(b.Name())] = b
}

// GetBuilder 获取构建器
func GetBuilder(name string) Builder {
	if b, ok := m[strings.ToLower(name)]; ok {
		return b
	}
	return nil
}
