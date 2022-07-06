package amocrm_v4

import "net/http"

type NoteType string
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

type note struct {
	Id                int        `json:"id"`
	EntityId          int        `json:"entity_id"`
	CreatedBy         int        `json:"created_by"`
	UpdatedBy         int        `json:"updated_by"`
	CreatedAt         int        `json:"created_at"`
	UpdatedAt         int        `json:"updated_at"`
	ResponsibleUserId int        `json:"responsible_user_id"`
	GroupId           int        `json:"group_id"`
	NoteType          string     `json:"note_type"`
	Params            noteParams `json:"params"`
	AccountId         int        `json:"account_id"`
	Links             links      `json:"_links"`
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

type getNotesQueryParams struct {
	entityType        string      `url:"-"`
	entityId          int         `url:"-"`
	path              string      `url:"-"`
	page              int         `url:"page,omitempty"`
	limit             int         `url:"limit,omitempty"`
	filter            interface{} `url:"filter,omitempty"`
	filterById        interface{} `url:"filter[id],omitempty"`
	filterByNoteType  interface{} `url:"filter[note_type],omitempty"`
	filterByUpdatedAt interface{} `url:"filter[updated_at],omitempty"`
	order             interface{} `url:"order,omitempty"`
}

type allNotes struct {
	Page     int   `json:"_page"`
	Links    links `json:"_links"`
	Embedded struct {
		Notes []*note `json:"notes"`
	} `json:"_embedded"`
}

func noteMultiplyRequest(params getNotesQueryParams) ([]*note, error) {
	var notes []*note

	for {
		var tmpNotes allNotes

		err := httpRequest(requestOpts{
			Method:        http.MethodGet,
			Path:          params.path,
			URLParameters: &params,
			Ret:           &tmpNotes,
		})
		if err != nil {
			return nil, err
		}

		notes = append(notes, tmpNotes.Embedded.Notes...)

		if len(tmpNotes.Links.Next.Href) > 0 {
			params.page = tmpNotes.Page + 1
		} else {
			break
		}
	}

	return notes, nil
}

func (nt *note) Delete() error {
	return httpRequest(requestOpts{
		Method:        http.MethodDelete,
		Path:          "/api/v4/notes/" + string(nt.Id),
		URLParameters: nil,
		Ret:           nil,
	})
}
