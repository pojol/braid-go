package balancer

import (
	"strings"

	"github.com/pojol/braid-go/module/discover"
	"github.com/pojol/braid-go/module/logger"
)

// Builder balancer builder
type Builder interface {
	Build(logger logger.ILogger) (IBalancer, error)
	Name() string
}

// IBalancer 负载均衡
type IBalancer interface {
	// 从服务节点列表中选取一个对应的节点，
	// 节点列表可以订阅discover模块的消息进行填充或更改，
	// braid 提供默认的`平滑加权轮询算法`如果有其他的需求，用户可以选择实现自定义的Pick接口。
	Pick() (nod discover.Node, err error)

	Add(discover.Node)
	Rmv(discover.Node)
	Update(discover.Node)
}

var (
	m = make(map[string]Builder)
)

// Register 注册balancer
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
