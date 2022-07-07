package amocrm_v4

type links struct {
	Self struct {
		Href string `json:"href,omitempty"`
	} `json:"self,omitempty"`
	Next struct {
		Href string `json:"href,omitempty"`
	} `json:"next,omitempty"`
	First struct {
		Href string `json:"href,omitempty"`
	} `json:"first,omitempty"`
	Prev struct {
		Href string `json:"href,omitempty"`
	} `json:"prev,omitempty"`
}
