package sirius

import "errors"

type EntityType string

const (
	EntityTypeLpa    = EntityType("lpa")
	EntityTypeEpa    = EntityType("epa")
	EntityTypePerson = EntityType("person")
)

func ParseEntityType(s string) (EntityType, error) {
	switch s {
	case "lpa":
		return EntityTypeLpa, nil
	case "epa":
		return EntityTypeEpa, nil
	case "person":
		return EntityTypePerson, nil
	}

	return EntityType(""), errors.New("could not parse entity type")
}
