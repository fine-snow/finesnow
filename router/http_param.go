// Http request partial fixed parameters

package router

type httpContentType string

var (
	textPlain       = new(httpContentType)
	applicationJson = new(httpContentType)
)

func init() {
	*textPlain = "text/plain"
	*applicationJson = "application/json"
}
