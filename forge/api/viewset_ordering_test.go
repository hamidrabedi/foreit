package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/forgego/forge/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type orderingTestModel struct {
	schema.BaseSchema
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	metaOrder []string
}

func (orderingTestModel) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64Field("id", schema.Primary()),
		schema.StringField("name"),
	}
}

func (m orderingTestModel) Meta() schema.Meta {
	return schema.Meta{TableName: "ordering_test_models", OrderBy: m.metaOrder}
}

type orderingTestQuerySet struct {
	fields []interface{}
}

func (*orderingTestQuerySet) Count(context.Context) (int64, error) { return 0, nil }
func (qs *orderingTestQuerySet) Offset(int) interface{}            { return qs }
func (qs *orderingTestQuerySet) Limit(int) interface{}             { return qs }
func (qs *orderingTestQuerySet) Filter(interface{}) interface{}    { return qs }
func (qs *orderingTestQuerySet) All(context.Context) ([]*orderingTestModel, error) {
	return nil, nil
}

func (qs *orderingTestQuerySet) OrderBy(fields ...interface{}) interface{} {
	qs.fields = append([]interface{}(nil), fields...)
	return qs
}

func TestBaseViewSet_ListAppliesStableDefaultOrdering(t *testing.T) {
	tests := []struct {
		name      string
		model     *orderingTestModel
		target    string
		wantField string
	}{
		{name: "model Meta ordering", model: &orderingTestModel{metaOrder: []string{"-name"}}, target: "/items/", wantField: "-name"},
		{name: "primary key fallback", model: &orderingTestModel{}, target: "/items/", wantField: "id"},
		{name: "explicit ordering wins", model: &orderingTestModel{metaOrder: []string{"-name"}}, target: "/items/?ordering=name", wantField: "name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := &orderingTestQuerySet{}
			vs := NewBaseViewSet(func() Serializer { return NewBaseSerializer(nil) }, qs, tt.model)
			recorder := httptest.NewRecorder()
			vs.List(recorder, httptest.NewRequest(http.MethodGet, tt.target, nil))

			require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
			assert.Equal(t, []interface{}{tt.wantField}, qs.fields)
		})
	}
}
