package clash

import "runtime"

func ForceGC() {
	runtime.GC()
}
