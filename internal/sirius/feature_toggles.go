package sirius

type FeatureToggles map[string]bool
type Config struct {
	FeatureToggles FeatureToggles `json:"featureToggles"`
}

func (c *Client) FeatureToggles(ctx Context) (FeatureToggles, error) {
	var v Config
	err := c.get(ctx, "/config", &v)
	if err != nil {
		return nil, err
	}

	return v.FeatureToggles, err
}

func (f FeatureToggles) Enabled(name string) bool {
	enabled, ok := f[name]
	return ok && enabled
}
