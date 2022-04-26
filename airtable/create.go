package airtable

import (
	"context"
	"errors"
)

func (a *airtable) Create(ctx context.Context, entity any) error {
	records := records{Records: []record{
		{Fields: fields(entity)},
	}}

	r, err := a.client.R().
		SetBody(records).
		SetContext(ctx).
		Post(tableName(entity))
	if err != nil {
		return err
	}

	if r.IsError() {
		return errors.New(string(r.Body()))
	}

	return nil
}
