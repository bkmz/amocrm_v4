package amocrm_v4

import (
	"fmt"
	"net/http"
)

type Ld struct{}
type leadNote note

type getLeadsQueryParams struct {
	With   []string    `url:"with,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
	Query  interface{} `url:"query,omitempty"`
	Filter interface{} `url:"filter,omitempty"`
	Order  interface{} `url:"order,omitempty"`
}

type lead struct {
	Id                 int         `json:"id"`
	Name               string      `json:"name"`
	Price              int         `json:"price"`
	ResponsibleUserId  int         `json:"responsible_user_id"`
	GroupId            int         `json:"group_id"`
	StatusId           int         `json:"status_id"`
	PipelineId         int         `json:"pipeline_id"`
	LossReasonId       interface{} `json:"loss_reason_id"`
	SourceId           interface{} `json:"source_id"`
	CreatedBy          int         `json:"created_by"`
	UpdatedBy          int         `json:"updated_by"`
	CreatedAt          int         `json:"created_at"`
	UpdatedAt          int         `json:"updated_at"`
	ClosedAt           int         `json:"closed_at"`
	ClosestTaskAt      interface{} `json:"closest_task_at"`
	IsDeleted          bool        `json:"is_deleted"`
	CustomFieldsValues interface{} `json:"custom_fields_values"`
	Score              interface{} `json:"score"`
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

type allLeads struct {
	Page  int `json:"_page"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
		First struct {
			Href string `json:"href"`
		} `json:"first"`
		Prev struct {
			Href string `json:"href"`
		} `json:"prev"`
	} `json:"_links"`
	Embedded struct {
		Leads []*lead `json:"leads"`
	} `json:"_embedded"`
}

func (l Ld) New() *lead {
	return &lead{}
}

func (l Ld) All() ([]*lead, error) {
	leads, err := l.multiplyRequest(getContactsQueryParams{
		Limit: 250,
	})
	if err != nil {
		return nil, err
	}

	return leads, nil
}

func (l Ld) multiplyRequest(opts getContactsQueryParams) ([]*lead, error) {
	var leads []*lead

	path := "/api/v4/leads"

	for {
		var tmpLeads allLeads

		err := httpRequest(requestOpts{
			Method:        http.MethodGet,
			Path:          path,
			URLParameters: &opts,
			Ret:           &tmpLeads,
		})
		if err != nil {
			return nil, err
		}

		leads = append(leads, tmpLeads.Embedded.Leads...)

		if len(tmpLeads.Links.Next.Href) > 0 {
			opts.Page = tmpLeads.Page + 1
		} else {
			break
		}
	}

	return leads, nil
}

func (ld *lead) Notes(params *getNotesQueryParams) ([]*leadNote, error) {
	path := fmt.Sprintf("/api/v4/leads/%d/notes", ld.Id)

	notes, err := ld.noteMultiplyRequest(getNotesQueryParams{
		limit: 250,
		path:  path,
	})

	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (ld *lead) noteMultiplyRequest(opts getNotesQueryParams) ([]*leadNote, error) {
	var notes []*leadNote

	for {
		var tmpNotes []*leadNote

		err := httpRequest(requestOpts{
			Method:        http.MethodGet,
			Path:          opts.path,
			URLParameters: &opts,
			Ret:           &tmpNotes,
		})
		if err != nil {
			return nil, err
		}

		notes = append(notes, tmpNotes...)

		if len(tmpNotes[0].Links.Next.Href) > 0 {
			opts.page = opts.page + 1
		} else {
			break
		}
	}

	return notes, nil
}

func (ldn *leadNote) New() *leadNote {
	return &leadNote{}
}

func (ldn *leadNote) Delete() error {
	return httpRequest(requestOpts{
		Method:        http.MethodDelete,
		Path:          fmt.Sprintf("/api/v4/leads/%d/notes/%d", ldn.EntityId, ldn.Id),
		URLParameters: nil,
		Ret:           nil,
	})
}
