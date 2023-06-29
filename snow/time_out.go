// Timeout function

package snow

import "time"

var (
	readTimeout       = 3 * time.Second
	readHeaderTimeout = 3 * time.Second
	writeTimeout      = 3 * time.Second
	idleTimeout       = time.Minute
)

func SetReadTimeout(t time.Duration) {
	readTimeout = t
}

func SetReadHeaderTimeout(t time.Duration) {
	readHeaderTimeout = t
}

func SetWriteTimeout(t time.Duration) {
	writeTimeout = t
}

func SetIdleTimeout(t time.Duration) {
	idleTimeout = t
}
