package amocrm_v4

import (
	"fmt"
	"net/http"
)

type Nt struct{}

type NoteType string
type NoteEntityType string
type MessageCashierNoteStatusType string

const (
	CommonNote                 NoteType = "common"
	CallInNote                 NoteType = "call_in"
	CallOutNote                NoteType = "call_out"
	ServiceMessageNote         NoteType = "service_message"
	ExtendedServiceMessageNote NoteType = "extended_service_message"
	MessageCashierNote         NoteType = "message_cashier"
	InvoicePaidNote            NoteType = "invoice_paid"
	GeolocationNote            NoteType = "geolocation"
	SmsInNote                  NoteType = "sms_in"
	SmsOutNote                 NoteType = "sms_out"
	AttachmentNote             NoteType = "attachment"
)

const (
	MessageCashierNoteStatusCreated  MessageCashierNoteStatusType = "created"
	MessageCashierNoteStatusShown    MessageCashierNoteStatusType = "shown"
	MessageCashierNoteStatusCanceled MessageCashierNoteStatusType = "canceled"
)

const (
	NoteEntityTypeLead    NoteEntityType = "leads"
	NoteEntityTypeContact NoteEntityType = "contacts"
)

type note struct {
	Id                int            `json:"id,omitempty"`
	EntityId          int            `json:"entity_id"`
	CreatedBy         int            `json:"created_by,omitempty"`
	UpdatedBy         int            `json:"updated_by,omitempty"`
	CreatedAt         int            `json:"created_at,omitempty"`
	UpdatedAt         int            `json:"updated_at,omitempty"`
	ResponsibleUserId int            `json:"responsible_user_id"`
	GroupId           int            `json:"group_id,omitempty"`
	NoteType          string         `json:"note_type"`
	Params            noteParams     `json:"params"`
	AccountId         int            `json:"account_id,omitempty"`
	Links             links          `json:"_links,omitempty"`
	RequestId         string         `json:"request_id,omitempty"`
	EntityType        NoteEntityType `json:"-"`
}

type noteParams struct {
	Text         string                       `json:"text,omitempty"`
	Service      string                       `json:"service,omitempty"`
	Uniq         string                       `json:"uniq,omitempty"`
	Duration     int                          `json:"duration,omitempty"`
	Source       string                       `json:"source,omitempty"`
	Link         string                       `json:"link,omitempty"`
	Phone        string                       `json:"phone,omitempty"`
	Status       MessageCashierNoteStatusType `json:"status,omitempty"`
	IconUrl      string                       `json:"icon_url,omitempty"`
	Address      string                       `json:"address,omitempty"`
	Longitude    string                       `json:"longitude,omitempty"`
	Latitude     string                       `json:"latitude,omitempty"`
	OriginalName string                       `json:"original_name,omitempty"`
	Attachment   string                       `json:"attachment,omitempty"`
}

type GetNotesQueryParams struct {
	//entityType        string      `url:"-"`
	//entityId          int         `url:"-"`
	//path              string      `url:"-"`
	Page              int         `url:"page,omitempty"`
	Limit             int         `url:"limit,omitempty"`
	Filter            interface{} `url:"filter,omitempty"`
	FilterById        interface{} `url:"filter[id],omitempty"`
	FilterByNoteType  interface{} `url:"filter[note_type],omitempty"`
	FilterByUpdatedAt interface{} `url:"filter[updated_at],omitempty"`
	Order             interface{} `url:"order,omitempty"`
}

type allNotes struct {
	Page     int   `json:"_page"`
	Links    links `json:"_links"`
	Embedded struct {
		Notes []*note `json:"notes"`
	} `json:"_embedded"`
}

// Create выполняет запрос на создание заметки
func (n *note) Create() (*allNotes, error) {
	path := fmt.Sprintf("/api/v4/%s/%d/notes", n.EntityType, n.EntityId)

	ret := allNotes{}

	return &ret, httpRequest(requestOpts{
		Path:           path,
		Method:         http.MethodPost,
		DataParameters: &n,
		Ret:            &ret,
	})
}

//func noteMultiplyRequest(params GetNotesQueryParams) ([]*note, error) {
//	var notes []*note
//
//	for {
//		var tmpNotes allNotes
//
//		err := httpRequest(requestOpts{
//			Method:        http.MethodGet,
//			Path:          params.path,
//			URLParameters: &params,
//			Ret:           &tmpNotes,
//		})
//		if err != nil {
//			return nil, err
//		}
//
//		notes = append(notes, tmpNotes.Embedded.Notes...)
//
//		if len(tmpNotes.Links.Next.Href) > 0 {
//			params.page = tmpNotes.Page + 1
//		} else {
//			break
//		}
//	}
//
//	return notes, nil
//}
