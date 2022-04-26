package airtable

import (
	"context"
	"errors"
)

func (a *airtable) Delete(ctx context.Context, entity any) error {
	r, err := a.client.R().Delete(tableName(entity) + "/" + id(entity))
	if err != nil {
		return err
	}

	if r.IsError() {
		return errors.New(string(r.Body()))
	}

	return nil
}
