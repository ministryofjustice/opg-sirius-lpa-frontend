package sirius

type MiReportConfig struct {
	Class       string `json:"class"`
	Description string `json:"description"`
	Fields      []struct {
		Name     string `json:"name"`
		Optional bool   `json:"optional"`
	} `json:"fields"`
}

type MiConfigField struct {
	Type    string         `json:"type"`
	MaxDate string         `json:"maxDate"`
	Options []MiConfigEnum `json:"options"`
}

type MiConfigEnum struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type MiConfig struct {
	Reports map[string]MiReportConfig `json:"reports"`

	Fields map[string]MiConfigField `json:"fields"`
}

func (c *Client) MiConfig(ctx Context) (MiConfig, error) {
	var v MiConfig
	if err := c.get(ctx, "/api/reporting/config", &v); err != nil {
		return MiConfig{}, err
	}

	return v, nil
}
