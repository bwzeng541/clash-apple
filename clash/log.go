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

	ch := make(chan L.Event, 1024)

	sub := L.Subscribe()
	defer L.UnSubscribe(sub)

	go func() {
		for elm := range sub {
			log := elm.(L.Event)
			select {
			case ch <- log:
			default:
			}
		}
		close(ch)
	}()

	for log := range ch {
		if log.LogLevel < L.Level() {
			continue
		}
		logger.Log(log.Type(), log.Payload)
	}
}
