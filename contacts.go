package amocrm_v4

import "net/http"

type Cntct struct{}

type GetContactsOpts struct {
	With   []string    `url:"with,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
	Query  interface{} `url:"query,omitempty"`
	Filter interface{} `url:"filter,omitempty"`
	Order  interface{} `url:"order,omitempty"`
}

type contact struct {
	Id                 int         `json:"id"`
	Name               string      `json:"name"`
	FirstName          string      `json:"first_name"`
	LastName           string      `json:"last_name"`
	ResponsibleUserId  int         `json:"responsible_user_id"`
	GroupId            int         `json:"group_id"`
	CreatedBy          int         `json:"created_by"`
	UpdatedBy          int         `json:"updated_by"`
	CreatedAt          int         `json:"created_at"`
	UpdatedAt          int         `json:"updated_at"`
	ClosestTaskAt      interface{} `json:"closest_task_at"`
	CustomFieldsValues interface{} `json:"custom_fields_values"`
	AccountId          int         `json:"account_id"`
	Links              struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Embedded struct {
		Tags      []interface{} `json:"tags"`
		Companies []interface{} `json:"companies"`
	} `json:"_embedded"`
}

type allContacts struct {
	Page  int `json:"_page"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
	} `json:"_links"`
	Embedded struct {
		Contacts []*contact `json:"contacts"`
	} `json:"_embedded"`
}

// Create Method creates empty struct
func (c Cntct) Create() *contact {
	return &contact{}
}

func (c Cntct) All() ([]*contact, error) {
	contacts, err := c.multiplyRequest(GetContactsOpts{})
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (c Cntct) Query() (*contact, error) {
	return nil, nil
}

func (c Cntct) multiplyRequest(opts GetContactsOpts) ([]*contact, error) {
	var contacts []*contact

	path := "/api/v4/contacts"

	for {
		var tmpContacts allContacts

		err := httpRequest(requestOpts{
			Method:        http.MethodGet,
			Path:          path,
			URLParameters: &opts,
			Ret:           &tmpContacts,
		})
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, tmpContacts.Embedded.Contacts...)

		if len(tmpContacts.Links.Next.Href) > 0 {
			path = tmpContacts.Links.Next.Href
		} else {
			break
		}
	}

	return contacts, nil
}
