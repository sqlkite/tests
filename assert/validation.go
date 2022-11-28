package assert

/*
Helper to validation a src.sqlkite.com/utils/validation.Result
but using reflection so we don't create a cyclical dependency
(ya, that's normal...)
*/

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type v struct {
	t      *testing.T
	json   []byte
	errors []map[string]any
}

func Validation(t *testing.T, result any) *v {
	e1 := reflect.ValueOf(result).MethodByName("Errors").Call(nil)[0]
	data, err := json.MarshalIndent(e1.Interface(), "", " ")
	if err != nil {
		panic(err)
	}

	var e2 []map[string]any
	if err := json.Unmarshal(data, &e2); err != nil {
		panic(err)
	}

	return &v{
		t:      t,
		json:   data,
		errors: e2,
	}
}

func (v *v) Fieldless(meta any, data ...map[string]any) *v {
	return v.Field("", meta, data...)
}

func (v *v) Field(expectedField string, meta any, data ...map[string]any) *v {
	t := v.t
	t.Helper()

	var expectedData map[string]any
	if len(data) == 1 {
		// pass through json so we have the same type as what we'll
		// be comparing too
		bytes, err := json.Marshal(data[0])
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal(bytes, &expectedData); err != nil {
			panic(err)
		}
	}
	expectedCode := int(reflect.ValueOf(meta).FieldByName("Code").Uint())

	for _, error := range v.errors {
		data := error["data"]
		field := error["field"]
		if field == nil && expectedField != "" && field.(string) != expectedField {
			continue
		}

		if int(error["code"].(float64)) != expectedCode {
			continue
		}

		if (data == nil && expectedData != nil) || (data != nil && expectedData == nil) {
			continue
		}

		if data != nil && expectedData != nil && !reflect.DeepEqual(data.(map[string]any), expectedData) {
			continue
		}

		return v
	}

	err := "\nexpected validation error:\n"
	if expectedField != "" {
		err += fmt.Sprintf("  field=%s\n", expectedField)
	}
	err += fmt.Sprintf("  code=%d\n", expectedCode)
	err += fmt.Sprintf("  data=%v\n\n", expectedData)
	err += fmt.Sprintf("got: %s", string(v.json))
	t.Error(err)
	t.FailNow()
	return v
}

func (v *v) FieldsHaveNoErrors(noFields ...string) *v {
	t := v.t
	t.Helper()

	for _, error := range v.errors {
		field := error["field"]
		if field == nil {
			continue
		}
		for _, noField := range noFields {
			if field == noField {
				t.Errorf("Expected no error for field '%s', but got:\n%v", field, error)
				t.FailNow()
			}
		}
	}
	return v
}
