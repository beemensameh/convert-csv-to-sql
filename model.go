package main

import (
	"encoding/json"
)

type RequirementStatusType string

const (
	Filled   RequirementStatusType = "filled"
	Unfilled RequirementStatusType = "unfilled"
)

type Requirement struct {
	ID              int              `json:"id"`
	UserUUID        string           `json:"user_uuid"`
	Value           *json.RawMessage `json:"value"`
	RequirementType string           `json:"requirement_type"`
}

type RequirementProviderType string

const (
	Rain       RequirementProviderType = "rain"
	Onfido     RequirementProviderType = "onfido"
	Compliance RequirementProviderType = "compliance"
	S3         RequirementProviderType = "s3"
)

type NationalitySchema struct {
	ResidenceCountryCode string   `json:"residence_country_code"`
	NationalityCountries []string `json:"nationality_countries"`
}

type CSVSheet struct {
	UUID      string `csv:"uid"`
	Country   string `csv:"country"`
	Residence string `csv:"nationality"`
}
