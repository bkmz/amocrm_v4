package amocrm_v4

import (
	"fmt"
	"net/http"
)

type Ld struct{}
type LeadWithType string

const (
	// LeadWithCatalogElements Добавляет в ответ связанные со сделками элементы списков
	LeadWithCatalogElements LeadWithType = "catalog_elements"

	// LeadWithIsPriceModifiedByRobot Добавляет в ответ свойство, показывающее,
	// изменен ли в последний раз бюджет сделки роботом
	LeadWithIsPriceModifiedByRobot LeadWithType = "is_price_modified_by_robot"

	// LeadWithLossReason Добавляет в ответ расширенную информацию по причине отказа
	LeadWithLossReason LeadWithType = "loss_reason"

	// LeadWithContacts Добавляет в ответ информацию о связанных со сделкой контактах
	LeadWithContacts LeadWithType = "contacts"

	// LeadWithOnlyDeleted Если передать данный параметр, то в ответе на запрос метода,
	// вернутся удаленные сделки, которые еще находятся в корзине. В ответ вы получите модель сделки,
	// у которой доступны дату изменения, ID пользователя сделавшего последнее изменение,
	// её ID и параметр is_deleted = true.
	LeadWithOnlyDeleted LeadWithType = "only_deleted"

	// LeadWithSourceID Добавляет в ответ ID источника
	LeadWithSourceID LeadWithType = "source_id"
)

type GetLeadsQueryParams struct {
	With   []string    `url:"with,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
	Query  interface{} `url:"query,omitempty"`
	Filter interface{} `url:"filter,omitempty"`
	Order  interface{} `url:"order,omitempty"`
}

type lead struct {
	Id                     int           `json:"id,omitempty"`                         //ID сделки
	Name                   string        `json:"name,omitempty"`                       //Название сделки
	Price                  int           `json:"price,omitempty"`                      //Бюджет сделки
	ResponsibleUserId      int           `json:"responsible_user_id,omitempty"`        //ID пользователя, ответственного за сделку
	GroupId                int           `json:"group_id,omitempty"`                   //ID группы, в которой состоит ответственны пользователь за сделку
	StatusId               int           `json:"status_id,omitempty"`                  //ID статуса, в который добавляется сделка, по-умолчанию – первый этап главной воронки
	PipelineId             int           `json:"pipeline_id,omitempty"`                //ID воронки, в которую добавляется сделка
	LossReasonId           interface{}   `json:"loss_reason_id,omitempty"`             //ID причины отказа
	SourceId               interface{}   `json:"source_id,omitempty"`                  //Требуется GET параметр with. ID источника сделки
	CreatedBy              int           `json:"created_by,omitempty"`                 //ID пользователя, создающий сделку
	UpdatedBy              int           `json:"updated_by,omitempty"`                 //ID пользователя, изменяющий сделку
	CreatedAt              int           `json:"created_at,omitempty"`                 //Дата создания сделки, передается в Unix Timestamp
	UpdatedAt              int           `json:"updated_at,omitempty"`                 //Дата изменения сделки, передается в Unix Timestamp
	ClosedAt               int           `json:"closed_at,omitempty"`                  //Дата закрытия сделки, передается в Unix Timestamp
	ClosestTaskAt          interface{}   `json:"closest_task_at,omitempty"`            //Дата ближайшей задачи к выполнению, передается в Unix Timestamp
	IsDeleted              bool          `json:"is_deleted,omitempty"`                 //Удалена ли сделка
	CustomFieldsValues     []CustomField `json:"custom_fields_values,omitempty"`       //Массив, содержащий информацию по значениям дополнительных полей, заданных для данной сделки
	Score                  interface{}   `json:"score,omitempty"`                      //Скоринг сделки
	AccountId              int           `json:"account_id,omitempty"`                 //ID аккаунта, в котором находится сделка
	IsPriceModifiedByRobot bool          `json:"is_price_modified_by_robot,omitempty"` //Требуется GET параметр with. Изменен ли в последний раз бюджет сделки роботом
	Embedded               struct {
		Tags     []Tag `json:"tags,omitempty"`
		Contacts []struct {
			Id     int  `json:"id,omitempty"`
			IsMain bool `json:"is_main,omitempty"`
		} `json:"contacts,omitempty"`
		Companies []struct {
			Id int `json:"id,omitempty"`
		} `json:"companies,omitempty"`
		CatalogElements []struct {
			Id        int         `json:"id,omitempty"`
			Metadata  interface{} `json:"metadata,omitempty"`
			Quantity  int         `json:"quantity,omitempty"`
			CatalogId int         `json:"catalog_id,omitempty"`
		} `json:"catalog_elements,omitempty"`
	} `json:"_embedded"`
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"_links,omitempty"`
}

type Leads []*lead

type allLeads struct {
	Page     int   `json:"_page"`
	Links    links `json:"_links"`
	Embedded struct {
		Leads []*lead `json:"leads"`
	} `json:"_embedded"`
}

func (l Ld) New() *lead {
	return &lead{}
}

func (l *lead) NewTask() *task {
	return &task{
		EntityType: TaskForLead,
		EntityId:   l.Id,
	}
}

func (l Ld) Create(leads Leads) (*allLeads, error) {
	ret := allLeads{}

	return &ret, httpRequest(requestOpts{
		Method:         http.MethodPost,
		Path:           "/api/v4/leads",
		DataParameters: &leads,
		Ret:            &ret,
	})
}

func (l Ld) Update(leads Leads) (*allLeads, error) {
	ret := allLeads{}

	return &ret, httpRequest(requestOpts{
		Method:         http.MethodPatch,
		Path:           "/api/v4/leads",
		DataParameters: &leads,
		Ret:            &ret,
	})
}

func (l Ld) All() ([]*lead, error) {
	leads, err := l.multiplyRequest(&GetLeadsQueryParams{
		Limit: 250,
	})
	if err != nil {
		return nil, err
	}

	return leads, nil
}

func (l Ld) Query(params *GetLeadsQueryParams) ([]*lead, error) {
	if params.Limit == 0 {
		params.Limit = 250
	}

	leads, err := l.multiplyRequest(params)
	if err != nil {
		return nil, err
	}

	return leads, nil
}

func (l Ld) ByID(id int) (*lead, error) {
	var ld *lead

	err := httpRequest(requestOpts{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/api/v4/leads/%d", id),
		Ret:    &ld,
	})
	if err != nil {
		return nil, err
	}

	return ld, nil
}

func (l Ld) multiplyRequest(params *GetLeadsQueryParams) ([]*lead, error) {
	var leads []*lead

	path := "/api/v4/leads"

	for {
		var tmpLeads allLeads

		err := httpRequest(requestOpts{
			Method:        http.MethodGet,
			Path:          path,
			URLParameters: &params,
			Ret:           &tmpLeads,
		})
		if err != nil {
			return nil, err
		}

		leads = append(leads, tmpLeads.Embedded.Leads...)

		if len(tmpLeads.Links.Next.Href) > 0 {
			params.Page = tmpLeads.Page + 1
		} else {
			break
		}
	}

	return leads, nil
}
