// API model

package doc

type ApiJson struct {
	Api  []*ApiModel
	Name string
	Url  string
}

type ApiModel struct {
	Group          string      `json:"group"`
	Module         string      `json:"module"`
	Name           string      `json:"name"`
	Url            string      `json:"url"`
	MethodType     string      `json:"methodType"`
	IsJsonTransfer bool        `json:"isJsonTransfer"`
	Params         []ApiParam  `json:"params"`
	Results        []ApiResult `json:"results"`
}

type ApiParam struct {
	Name string
}

type ApiResult struct {
	Name string
}

func ParsingAPIFunctionAnnotations() {

}
