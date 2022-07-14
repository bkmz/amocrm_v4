package amocrm_v4

type CustomField struct {
	FieldId   int     `json:"field_id"`
	FieldName string  `json:"field_name"`
	FieldCode *string `json:"field_code"`
	FieldType string  `json:"field_type"`
	Values    []struct {
		Value    string `json:"value"`
		EnumId   int    `json:"enum_id,omitempty"`
		EnumCode string `json:"enum_code,omitempty"`
	} `json:"values"`
}
