// Http request partial fixed parameters

package router

type httpMethod string
type httpContentType string

var (
	HttpMethodGet     = new(httpMethod)
	HttpMethodPost    = new(httpMethod)
	HttpMethodOptions = "OPTIONS"
	textPlain         = new(httpContentType)
	applicationJson   = new(httpContentType)
)

func init() {
	*HttpMethodGet = "GET"
	*HttpMethodPost = "POST"
	*textPlain = "text/plain"
	*applicationJson = "application/json"
}
