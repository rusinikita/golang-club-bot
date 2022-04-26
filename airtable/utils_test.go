package airtable

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Test struct {
	RecordID      string
	Name          string
	Note          string
	Condition     string `json:",omitempty"`
	Status        string
	Created       time.Time `json:",omitempty"`
	StatusUpdated time.Time `json:",omitempty"`
}

type Test2 struct {
	ID string
}

func (t Test2) TableName() string {
	return "Test"
}

func TestTableName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "Tests", tableName(&Test{}))
	assert.Equal(t, "Tests", tableName([]Test{}))
	assert.Equal(t, "Tests", tableName(&[]Test{}))
	assert.Equal(t, "Test", tableName(&[]Test2{}))
}

func TestFields(t *testing.T) {
	t.Parallel()

	entity := Test{
		RecordID: "1",
		Name:     "bla",
		Note:     "bla",
		Status:   "todo",
	}

	// no Condition, Created, StatusUpdated field cause omitempty
	// no CreatedField cause computed in
	expected := map[string]any{
		"Name":   entity.Name,
		"Note":   entity.Note,
		"Status": entity.Status,
	}

	assert.Equal(t, expected, fields(entity))
}
