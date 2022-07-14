package amocrm_v4

import (
	"fmt"
	"net/http"
)

type (
	Tsk                     struct{}
	TaskTypeIdType          int
	TaskEntityType          string
	FilterByIsCompletedType int
	OrderType               string
	OrderDirectionType      string
	Tasks                   []*task
)

const (
	// TaskCall задача – Звонок
	TaskCall TaskTypeIdType = iota + 1
	// TaskMeeting задача – Встреча
	TaskMeeting

	TaskForLead       TaskEntityType = "leads"
	TaskForContact    TaskEntityType = "contacts"
	TasksForCompany   TaskEntityType = "companies"
	TasksForCustomers TaskEntityType = "customers"

	FilterByIsCompletedTrue  FilterByIsCompletedType = 1
	FilterByIsCompletedFalse FilterByIsCompletedType = 0

	OrderByCreatedAt    OrderType = "created_at"
	OrderByCompleteTill OrderType = "complete_till"
	OrderById           OrderType = "id"

	OrderAsc  OrderDirectionType = "asc"
	OrderDesc OrderDirectionType = "desc"
)

type task struct {
	Id                int            `json:"id,omitempty"`                  // Id задачи
	CreatedBy         int            `json:"created_by,omitempty"`          // ID пользователя, создавшего задачу
	UpdatedBy         int            `json:"updated_by,omitempty"`          // ID пользователя, изменившего задачу
	CreatedAt         int            `json:"created_at,omitempty"`          // Дата создания задачи, передается в UNIX Timestamp
	UpdatedAt         int            `json:"updated_at,omitempty"`          // Дата изменения задачи, передается в UNIX Timestamp
	ResponsibleUserId int            `json:"responsible_user_id,omitempty"` // ID пользователя, ответственного за задачу
	GroupID           int            `json:"group_id,omitempty"`            // ID группы, в которой состоит пользователь ответственный за задачу
	EntityId          int            `json:"entity_id,omitempty"`           // ID сущности, к которой привязана задача
	EntityType        TaskEntityType `json:"entity_type,omitempty"`         // Тип сущности, к которой привязана задача
	IsCompleted       bool           `json:"is_completed,omitempty"`        // Флаг завершенности задачи
	TaskTypeId        TaskTypeIdType `json:"task_type_id,omitempty"`        // Тип задачи
	Text              string         `json:"text,omitempty"`                // Описание задачи
	Duration          int            `json:"duration,omitempty"`            // Длительность задачи в секундах
	CompleteTill      int            `json:"complete_till,omitempty"`       // Дата завершения задачи, передается в UNIX Timestamp
	RequestId         string         `json:"request_id,omitempty"`          // Идентификатор запроса, никак не обрабатывается в API и возвращается неизмененным
	AccountID         int            `json:"account_id,omitempty"`          // ID аккаунта, в котором создана задача
	Result            struct {
		Text string `json:"text,omitempty"` // Текст результата задачи
	} `json:"result,omitempty"` // Результат выполнения задачи
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"_links,omitempty"`
}

type allTasks struct {
	Page     int   `json:"_page,omitempty"`
	Links    links `json:"_links"`
	Embedded struct {
		Tasks []*task `json:"tasks"`
	} `json:"_embedded"`
}

type GetTaskQueryParams struct {
	Page                      int                     `url:"page,omitempty"`                          // Страница выборки
	Limit                     int                     `url:"limit,omitempty"`                         // Количество возвращаемых сущностей за один запрос (Максимум – 250)
	FilterByResponsibleUserId []int                   `url:"filter[responsible_user_id][],omitempty"` // Фильтр по ID ответственного за задачу пользователя. Можно передать как один ID, так и массив из нескольких ID
	FilterByIsCompleted       FilterByIsCompletedType `url:"filter[is_completed][],omitempty"`        // Фильтр по статусу задачи.
	FilterByTaskType          []TaskEntityType        `url:"filter[task_type][],omitempty"`           // Фильтр по ID типа задачи. Можно передать как один ID, так и массив из нескольких ID
	FilterByEntityType        TaskEntityType          `url:"filter[entity_type][],omitempty"`         // Фильтр по типу привязанной к задаче сущности. Возможные значения: leads, contacts, companies, customers
	FilterByEntityId          []int                   `url:"filter[entity_id][],omitempty"`           // Фильтр по ID, привязанной к задаче, сущности. Для его использования необходимо передать значение в filter[entity_type]. Можно передать как один ID, так и массив из нескольких ID
	FilterById                []int                   `url:"filter[id][],omitempty"`                  // Фильтр по ID задачи. Можно передать как один ID, так и массив из нескольких ID
	FilterByUpdatedAt         int                     `url:"filter[updated_at][],omitempty"`          // Фильтр по дате последнего изменения задачи. Можно передать timestamp, в таком случае будут возвращены задачи, которые были изменены после переданного значения. Также можно передать массив вида filter[updated_at][from]=… и filter[updated_at][to]=…, для фильтрации по значениям ОТ и ДО.
	FilterByUpdatedAtFrom     int                     `url:"filter[updated_at][from],omitempty"`
	FilterByUpdatedAtTo       int                     `url:"filter[updated_at][to],omitempty"`
	OrderByCreatedAt          OrderDirectionType      `url:"order[created_at],omitempty"`
	OrderBuCompleteTill       OrderDirectionType      `url:"order[complete_till],omitempty"`
	OrderById                 OrderDirectionType      `url:"order[id],omitempty"`
}

func (t Tsk) New() *task {
	return &task{}
}

// Create Создает новую задачу.
// Для создания задачи нужно передать 2 обязательных параметра:
// text и complete_till.
func (t Tsk) Create(tsk Tasks) (*allTasks, error) {
	ret := allTasks{}

	return &ret, httpRequest(requestOpts{
		Method:         http.MethodPost,
		Path:           "/api/v4/tasks",
		DataParameters: &tsk,
		Ret:            &ret,
	})
}

// Update Обновляет задачу. Данный метод может использоваться для пакетного обновления задач.
func (t Tsk) Update(tsk Tasks) (*allTasks, error) {
	ret := allTasks{}

	return &ret, httpRequest(requestOpts{
		Method:         http.MethodPatch,
		Path:           "/api/v4/tasks",
		DataParameters: &tsk,
		Ret:            &ret,
	})
}

// Update Обновляет задачу. Данный метод используется для индивидуального обновления задачи.
func (t *task) Update() (*allTasks, error) {
	ret := allTasks{}

	return &ret, httpRequest(requestOpts{
		Method:         http.MethodPatch,
		Path:           fmt.Sprintf("/api/v4/tasks/%d", t.Id),
		DataParameters: &t,
		Ret:            &ret,
	})
}

//// Complete Обновляет задаче статус выполнения. Данный метод может использоваться для пакетного обновления задач.
//func (t Tsk) Complete(tsk Tasks) (*allTasks, error) {
//	ret := allTasks{}
//
//	return &ret, httpRequest(requestOpts{
//		Method:         http.MethodPatch,
//		Path:           "/api/v4/tasks",
//		DataParameters: &tsk,
//		Ret:            &ret,
//	})
//}

// Complete Обновляет задаче статус выполнения. Данный метод используется для индивидуального обновления задачи.
func (t *task) Complete(result string) (*task, error) {
	t.IsCompleted = true
	t.Result.Text = result

	return t, httpRequest(requestOpts{
		Method:         http.MethodPatch,
		Path:           fmt.Sprintf("/api/v4/tasks/%d", t.Id),
		DataParameters: &t,
		Ret:            &t,
	})
}

// All Возвращает список всех задач.
func (t Tsk) All() (Tasks, error) {
	return t.multiplyRequest(&GetTaskQueryParams{
		Limit: 250,
	})
}

// Query Возвращает список задач по заданным параметрам.
func (t Tsk) Query(params GetTaskQueryParams) (Tasks, error) {
	if params.Limit == 0 {
		params.Limit = 250
	}

	return t.multiplyRequest(&GetTaskQueryParams{
		Limit: 250,
	})
}

func (t Tsk) ByID(id int) (*task, error) {
	ret := task{}

	return &ret, httpRequest(requestOpts{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/api/v4/tasks/%d", id),
		Ret:    &ret,
	})
}

func (t Tsk) multiplyRequest(params *GetTaskQueryParams) (Tasks, error) {
	var tasks Tasks

	path := "/api/v4/tasks"

	for {
		var tmpTasks allTasks

		err := httpRequest(requestOpts{
			Method:        http.MethodGet,
			Path:          path,
			URLParameters: &params,
			Ret:           &tmpTasks,
		})
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, tmpTasks.Embedded.Tasks...)

		if len(tmpTasks.Links.Next.Href) > 0 {
			params.Page = tmpTasks.Page + 1
		} else {
			break
		}
	}

	return tasks, nil
}
