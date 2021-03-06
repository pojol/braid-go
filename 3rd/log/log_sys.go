package log

import (
	"runtime/debug"

	"go.uber.org/zap"
)

/*
	诊断日志
	包含如
	1. 系统&进程状态日志
	2. 定期任务状态时间日志
	3. 程序追踪日志（警告，错误
	4. request慢日志
*/

// 日志类型
const (
	LogElector = "LogElector"
)

// SysError 系统错误日志
func SysError(module string, function string, desc string) {

	logPtr.SysLog.Error(desc, // err msg
		zap.String("stack", string(debug.Stack())),
		zap.String("module", module), // 模块
		zap.String("func", function), // 函数
	)
}

// SysSlow 慢日志
func SysSlow(apiName string, requestid string, et int, desc string) {

	logPtr.SysLog.Warn(desc,
		zap.String("api", apiName),         // trace name
		zap.String("requestID", requestid), // trace id
		zap.Int("executionTime", et),       // 总计执行时间
	)

}

// SysRoutingError 路由警告日志
func SysRoutingError(serviceName string, desc string) {
	logPtr.SysLog.Warn(desc,
		zap.String("service", serviceName),
	)
}

// SysCompose 输出启动节点
func SysCompose(nods []string, desc string) {
	logPtr.SysLog.Info(desc,
		zap.Strings("nods", nods),
	)
}

// SysElection 选举日志
func SysElection(nod string, session string) {
	logPtr.SysLog.Info("current master nod",
		zap.String("logt", LogElector),
		zap.String("node", nod),
		zap.String("session", session),
	)
}

// SysWelcome 欢迎日志
func SysWelcome(nodeName string, mode string, ty string, info string) {
	logPtr.SysLog.Info(info,
		zap.String("node", nodeName),
		zap.String("mode", mode),
		zap.String("type", ty),
	)
}
