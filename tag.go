package amocrm_v4

type Tg struct{}

type allTags struct {
	Page     int   `json:"_page"`
	Links    links `json:"_links"`
	Embedded struct {
		Tags []Tag `json:"tags"`
	} `json:"_embedded"`
}

type Tag struct {
	Id        int         `json:"id,omitempty"`
	Name      string      `json:"name"`
	Color     interface{} `json:"color"`
	RequestID string      `json:"request_id"`
}
