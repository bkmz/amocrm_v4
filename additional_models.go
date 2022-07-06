package amocrm_v4

type links struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self,omitempty"`
	Next struct {
		Href string `json:"href"`
	} `json:"next,omitempty"`
	First struct {
		Href string `json:"href"`
	} `json:"first,omitempty"`
	Prev struct {
		Href string `json:"href"`
	} `json:"prev,omitempty"`
}
