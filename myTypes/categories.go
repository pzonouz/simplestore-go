package myTypes

import "github.com/jackc/pgx/v5/pgtype"

type ListCategoriesTreeRowDTO struct {
	ID       pgtype.UUID                `json:"id"`
	Name     pgtype.Text                `json:"name"`
	Slug     pgtype.Text                `json:"slug"`
	Priority pgtype.Numeric             `json:"priority"`
	Children []ListCategoriesTreeRowDTO `json:"children"`
}
