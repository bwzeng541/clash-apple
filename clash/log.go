package clash

import L "github.com/Dreamacro/clash/log"

var (
	logger NativeLogger
)

type NativeLogger interface {
	Log(level string, payload string)
}

func SetNativeLogger(l NativeLogger) {
	logger = l
}

func SetLogLevel(level string) {
	if basic == nil {
		return
	}
	L.SetLevel(L.LogLevelMapping[level])
}

func fetchLogs() {
	sub := L.Subscribe()
	defer L.UnSubscribe(sub)
	for elm := range sub {
		if logger == nil {
			continue
		}
		log := elm.(*L.Event)
		if log.LogLevel < L.Level() {
			continue
		}
		logger.Log(log.Type(), log.Payload)
	}
}
