package airtable

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Options struct {
	View   string // view for records filter
	Filter any    // entity obj with not zero values filter fields
}

func View(view string) Options {
	return Options{View: view}
}

func Filter(entity any) Options {
	return Options{Filter: entity}
}

func (a *airtable) List(ctx context.Context, list any, options ...Options) error {
	records := records{}

	request := a.client.R().
		SetContext(ctx).
		SetResult(&records)

	if len(options) > 0 {
		o := options[0]

		if o.View != "" {
			request = request.SetQueryParam("view", o.View)
		}

		if o.Filter != nil {
			request = request.SetQueryParam("filterByFormula", filterString(o.Filter))
		}
	}

	r, err := request.Get(tableName(list))
	if err != nil {
		return err
	}

	if r.IsError() {
		return errors.New(string(r.Body()))
	}

	// mapping
	maps := make([]map[string]any, len(records.Records))
	for i, r := range records.Records {
		r.Fields["RecordID"] = r.ID

		maps[i] = r.Fields
	}

	return decode(maps, list)
}

func filterString(filter any) string {
	if s, ok := filter.(string); ok {
		return s
	}

	var fieldFilters []string //nolint:prealloc // can't preallocate

	for key, value := range nonZeroFields(filter) {
		// potential bug with '' wrapping
		fieldFilter := fmt.Sprintf("{%s} = '%v'", key, value)

		fieldFilters = append(fieldFilters, fieldFilter)
	}

	return fmt.Sprintf("AND(%s)", strings.Join(fieldFilters, ", "))
}

func decode(data any, result any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  result,
		TagName: "json",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			intToDuration,
		),
	})
	if err != nil {
		return err
	}

	return decoder.Decode(data)
}

func intToDuration(_, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeOf(time.Duration(0)) {
		return data, nil
	}

	return time.ParseDuration(fmt.Sprint(data) + "s")
}
