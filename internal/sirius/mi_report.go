package sirius

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type miReportData struct {
	Data *MiReportResponse `json:"data"`
}

type MiReportResponse struct {
	ResultCount       int    `json:"result_count"`
	ReportType        string `json:"report_type"`
	ReportDescription string `json:"report_description"`
}

type miReportError struct {
	Detail string `json:"detail"`
}

func (c *Client) MiReport(ctx Context, params url.Values) (*MiReportResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/reporting/applications?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var v miReportError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return nil, err
		}

		return nil, errors.New(v.Detail)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v miReportData
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v.Data, err
}
