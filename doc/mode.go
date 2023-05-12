// API model

package doc

type apiModel struct {
	Group          string      `json:"group"`
	Module         string      `json:"module"`
	Name           string      `json:"name"`
	Url            string      `json:"url"`
	MethodType     string      `json:"methodType"`
	IsJsonTransfer bool        `json:"isJsonTransfer"`
	Params         []apiParam  `json:"params"`
	Results        []apiResult `json:"results"`
}

type apiParam struct {
	name string
}

type apiResult struct {
	name string
}
