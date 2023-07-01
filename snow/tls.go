// Server TLS

package snow

var (
	certFile string
	keyFile  string
)

func SetCertFile(url string) {
	certFile = url
}

func SetKeyFile(url string) {
	keyFile = url
}
