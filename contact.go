package amocrm_v4

import (
	"fmt"
	"net/http"
)

type Ct struct{}
type contactNote note

type ContactWithType string

const (
	ContactWithLeads           ContactWithType = "leads"
	ContactWithCustomers       ContactWithType = "customers"
	ContactWithCatalogElements ContactWithType = "catalog_elements"
)

type GetContactsQueryParams struct {
	With   []ContactWithType `url:"with,omitempty"`
	Limit  int               `url:"limit,omitempty"`
	Page   int               `url:"page,omitempty"`
	Query  interface{}       `url:"query,omitempty"`
	Filter interface{}       `url:"filter,omitempty"`
	Order  interface{}       `url:"order,omitempty"`
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
	Links              links       `json:"_links"`
	Embedded           struct {
		Customers       []interface{} `json:"customers"`
		Leads           []*lead       `json:"leads"`
		CatalogElements []interface{} `json:"catalog_elements"`
		Tags            []interface{} `json:"tags"`
		Companies       []interface{} `json:"companies"`
	} `json:"_embedded"`
}

type allContacts struct {
	Page     int   `json:"_page"`
	Links    links `json:"_links"`
	Embedded struct {
		Contacts []*contact `json:"contacts"`
	} `json:"_embedded"`
}

type allContactNotes struct {
	Page     int   `json:"_page"`
	Links    links `json:"_links"`
	Embedded struct {
		Notes []*contactNote `json:"notes"`
	} `json:"_embedded"`
}

// New Method creates empty struct
func (c Ct) New() *contact {
	return &contact{}
}

func (ct *contact) NewNote() *note {
	return &note{
		EntityId:   ct.Id,
		EntityType: NoteEntityTypeContact,
	}
}

func (c Ct) All() ([]*contact, error) {
	contacts, err := c.multiplyRequest(&GetContactsQueryParams{
		Limit: 250,
	})
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (c Ct) Query(params *GetContactsQueryParams) ([]*contact, error) {
	if params.Limit == 0 {
		params.Limit = 250
	}

	contacts, err := c.multiplyRequest(params)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (c Ct) ByID(id int, with []ContactWithType) (*contact, error) {
	var ct *contact

	opts := GetContactsQueryParams{
		With: with,
	}

	err := httpRequest(requestOpts{
		Method:        http.MethodGet,
		Path:          fmt.Sprintf("/api/v4/contacts/%d", id),
		URLParameters: &opts,
		Ret:           &ct,
	})
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (ct *contact) Notes(params *GetNotesQueryParams) ([]*contactNote, error) {
	notes, err := ct.noteMultiplyRequest(params)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (ct *contact) noteMultiplyRequest(opts *GetNotesQueryParams) ([]*contactNote, error) {
	var notes []*contactNote

	if opts.Limit == 0 {
		opts.Limit = 250
	}

	for {
		var tmpNotes allContactNotes

		path := fmt.Sprintf("/api/v4/contacts/%d/notes", ct.Id)
		err := httpRequest(requestOpts{
			Method:        http.MethodGet,
			Path:          path,
			URLParameters: &opts,
			Ret:           &tmpNotes,
		})
		if err != nil {
			return nil, fmt.Errorf("ошибка обработки запроса %s: %s", path, err)
		}

		notes = append(notes, tmpNotes.Embedded.Notes...)

		if len(tmpNotes.Links.Next.Href) > 0 {
			opts.Page = opts.Page + 1
		} else {
			break
		}
	}

	return notes, nil
}

func (c Ct) multiplyRequest(opts *GetContactsQueryParams) ([]*contact, error) {
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
			opts.Page = tmpContacts.Page + 1
		} else {
			break
		}
	}

	return contacts, nil
}
