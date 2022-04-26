package airtable

import (
	"context"
	"errors"
)

func (a *airtable) Patch(ctx context.Context, entity any) error {
	r, err := a.client.R().
		SetContext(ctx).
		SetBody(records{Records: []record{
			{
				ID:     id(entity),
				Fields: nonZeroFields(entity),
			},
		}}).
		Patch(tableName(entity))
	if err != nil {
		return err
	}

	if r.IsError() {
		return errors.New(string(r.Body()))
	}

	return nil
}
