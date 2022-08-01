package sirius

import (
	"encoding/json"
	"fmt"
)

type miConfig struct {
	Data struct {
		Items []struct {
			Properties map[string]MiConfigProperty `json:"properties"`
		} `json:"items"`
	} `json:"data"`
}

type MiConfigProperty struct {
	Description   string            `json:"description"`
	Type          string            `json:"type"`
	Required      bool              `json:"required"`
	Enum          []MiConfigEnum    `json:"enum"`
	DependsOn     MiConfigDependsOn `json:"dependsOn"`
	Format        string            `json:"format"`
	FormatMaximum string            `json:"formatMaximum"`
}

type MiConfigDependsOn struct {
	ReportType []MiConfigReportType `json:"reportType"`
}

type MiConfigReportType struct {
	Name string `json:"name"`
}

type MiConfigEnum struct {
	Name        string
	Description string
}

type rawMiConfigEnum struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Label       string `json:"label"`
}

func (e *MiConfigEnum) UnmarshalJSON(text []byte) error {
	var v rawMiConfigEnum
	if err := json.Unmarshal(text, &v); err == nil {
		if v.Name != "" {
			e.Name = v.Name
			e.Description = v.Description
		}
		if v.Value != "" {
			e.Name = v.Value
			e.Description = v.Label
		}
	}

	var s string
	if err := json.Unmarshal(text, &s); err == nil {
		e.Name = s
		e.Description = s
	}

	if e.Name == "" {
		return fmt.Errorf("could not unmarshal '%s' to MiConfigEnum", text)
	}

	return nil
}

func (c *Client) MiConfig(ctx Context) (map[string]MiConfigProperty, error) {
	var v miConfig
	if err := c.get(ctx, "/api/reporting/config", &v); err != nil {
		return nil, err
	}

	return v.Data.Items[0].Properties, nil
}
