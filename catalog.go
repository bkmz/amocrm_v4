package amocrm_v4

import (
	"fmt"
	"net/http"
	"strconv"
)

type (
	Ctg             struct{}
	CatalogType     string
	Catalogs        []*catalog
	Elements        []*element
	ElementWithType string
)

type GetCatalogsQueryParams struct {
	Page  int `url:"page"`
	Limit int `url:"limit"`
}

type GetCatalogElementsQueryParams struct {
	Page       int             `url:"page"`         //Страница выборки
	Limit      int             `url:"limit"`        //Количество возвращаемых сущностей за один запрос (Максимум – 250)
	Query      string          `url:"query"`        //Поисковый запрос (Осуществляет поиск по заполненным полям сущности)
	FilterByID string          `url:"filter_by_id"` //Фильтр по ID элемента. Можно передать как один ID, так и массив из нескольких ID
	With       ElementWithType `url:"with"`         //Поля для выборки
}

const (
	CatalogRegular  CatalogType = "regular"
	CatalogInvoices CatalogType = "invoices"
	CatalogProducts CatalogType = "products"

	ElementWithInvoiceLink ElementWithType = "invoice_link" // При передаче данного параметра, вернется дополнительное свойство invoice_link, содержащие ссылку на печатную форму счета. Если передать этот параметр с отличным от списка Счетов списком, то вернется null.
)

type catalog struct {
	Id              int         `json:"id,omitempty"`                // ID списка
	Name            string      `json:"name,omitempty"`              // Название списка
	CreatedBy       int         `json:"created_by,omitempty"`        // ID пользователя, создавший список
	UpdatedBy       int         `json:"updated_by,omitempty"`        // ID пользователя, изменивший список последним
	CreatedAt       int         `json:"created_at,omitempty"`        // Дата создания списка, передается в Unix Timestamp
	UpdatedAt       int         `json:"updated_at,omitempty"`        // Дата изменения списка, передается в Unix Timestamp
	Sort            int         `json:"sort,omitempty"`              // Сортировка списка
	Type            CatalogType `json:"type,omitempty"`              // Тип списка
	CanAddElements  bool        `json:"can_add_elements,omitempty"`  // Можно ли добавлять элементы списка из интерфейса (Применяется только для списка счетов)
	CanShowInCards  bool        `json:"can_show_in_cards,omitempty"` // Должна ли добавляться вкладка со списком в карточку сделки/покупателя (Применяется только для списка счетов)
	CanLinkMultiple bool        `json:"can_link_multiple,omitempty"` // Если ли возможность привязывать один элемент данного списка к нескольким сделкам/покупателям
	CanBeDeleted    bool        `json:"can_be_deleted,omitempty"`    // Может ли список быть удален через интерфейс
	SdkWidgetCode   interface{} `json:"sdk_widget_code,omitempty"`   // Код виджета, который управляет списком и может отобразить своё собственное окно редактирования элемента (Применяется только для списка счетов)
	AccountId       int         `json:"account_id,omitempty"`        // ID аккаунта, в котором находится список
	RequestId       string      `json:"request_id,omitempty"`        // Поле, которое вернется вам в ответе без изменений и не будет сохранено. Необязательный параметр
	Links           struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"_links,omitempty"`
}

type allCatalogs struct {
	Page     int   `json:"_page"`
	Links    links `json:"_links"`
	Embedded struct {
		Catalogs Catalogs `json:"catalogs"`
	} `json:"_embedded"`
}

type element struct {
	Id                 int           `json:"id,omitempty"`         //ID элемента списка
	CatalogId          int           `json:"catalog_id,omitempty"` //ID списка
	Name               string        `json:"name,omitempty"`       //Название элемента
	CreatedBy          int           `json:"created_by,omitempty"` //ID пользователя, создавший элемент
	UpdatedBy          int           `json:"updated_by,omitempty"` //ID пользователя, изменивший элемент последним
	CreatedAt          int           `json:"created_at,omitempty"` //Дата создания элемента, передается в Unix Timestamp
	UpdatedAt          int           `json:"updated_at,omitempty"` //Дата изменения элемента, передается в Unix Timestamp
	IsDeleted          bool          `json:"is_deleted,omitempty"` //Удален ли элемент
	CustomFieldsValues []CustomField `json:"custom_fields_values"`
	AccountId          int           `json:"account_id"`
	Links              struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

type allCatalogElements struct {
	Page     int   `json:"_page"`
	Links    links `json:"_links"`
	Embedded struct {
		Elements []*element `json:"elements"`
	} `json:"_embedded"`
}

//All Метод позволяет получить доступные списки в аккаунте.
func (c Ctg) All() (*Catalogs, error) {
	req := GetCatalogsQueryParams{
		Limit: 250,
	}
	ret := allCatalogs{}

	return &ret.Embedded.Catalogs, httpRequest(requestOpts{
		Method:        http.MethodGet,
		Path:          "/api/v4/catalogs",
		URLParameters: req,
		Ret:           &ret,
	})
}

// ByID Метод позволяет получить данные конкретного списка по ID.
func (c Ctg) ByID(id int) (*catalog, error) {
	ret := catalog{}

	return &ret, httpRequest(requestOpts{
		Method: http.MethodGet,
		Path:   "/api/v4/catalogs/" + strconv.Itoa(id),
		Ret:    &ret,
	})
}

func (c Ctg) New() *catalog {
	return &catalog{}
}

// Create Метод позволяет добавлять списки в аккаунт пакетно.
func (c Ctg) Create(catalogs Catalogs) (*allCatalogs, error) {
	ret := allCatalogs{}

	return &ret, httpRequest(requestOpts{
		Method:         http.MethodPost,
		Path:           "/api/v4/catalogs",
		DataParameters: &catalogs,
		Ret:            &ret,
	})
}

//TODO: PATCH /api/v4/catalogs
//TODO: PATCH /api/v4/catalogs/{id}

func (c *catalog) AllElements() (Elements, error) {
	return c.multiplyRequest(&GetCatalogElementsQueryParams{
		Limit: 250,
	})
}

func (c *catalog) QueryElements(opts *GetCatalogElementsQueryParams) (Elements, error) {
	if opts.Limit == 0 {
		opts.Limit = 250
	}

	return c.multiplyRequest(opts)
}

func (c *catalog) multiplyRequest(opts *GetCatalogElementsQueryParams) (Elements, error) {
	var elements Elements

	if opts.Limit == 0 {
		opts.Limit = 250
	}

	for {
		var tmpElements allCatalogElements

		path := fmt.Sprintf("/api/v4/catalogs/%d/elements", c.Id)
		err := httpRequest(requestOpts{
			Method:        http.MethodGet,
			Path:          path,
			URLParameters: &opts,
			Ret:           &tmpElements,
		})
		if err != nil {
			return nil, fmt.Errorf("ошибка обработки запроса %s: %s", path, err)
		}

		elements = append(elements, tmpElements.Embedded.Elements...)

		if len(tmpElements.Links.Next.Href) > 0 {
			opts.Page = opts.Page + 1
		} else {
			break
		}
	}

	return elements, nil
}

func GetAllProducts() {

}
