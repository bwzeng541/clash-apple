package clash

import L "github.com/Dreamacro/clash/log"

var (
	logger RealTimeLogger
)

type RealTimeLogger interface {
	Log(level string, payload string)
}

func SetRealTimeLogger(l RealTimeLogger) {
	logger = l
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

func SetLogLevel(level string) {
	if basic == nil {
		return
	}
	L.SetLevel(L.LogLevelMapping[level])
}
