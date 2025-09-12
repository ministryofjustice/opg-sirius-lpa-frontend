package server

import "net/http"

func validateRemoveAttorneyPage(r *http.Request, data *removeAnAttorneyData) {
	if data.Form.RemovedAttorneyUid == "" {
		data.Error.Field["removeAttorney"] = map[string]string{
			"reason": "Please select an attorney for removal",
		}
	}

	if data.Form.RemovedReason == "" {
		data.Error.Field["removedReason"] = map[string]string{
			"reason": "Please select a reason for removal",
		}
	}

	if len(data.Form.EnabledAttorneyUids) > 0 && postFormCheckboxChecked(r, "skipEnableAttorney", "yes") {
		data.Error.Field["enableAttorney"] = map[string]string{
			"reason": "Please do not select both a replacement attorney and the option to skip",
		}
	}

	if len(data.Form.EnabledAttorneyUids) == 0 && !postFormCheckboxChecked(r, "skipEnableAttorney", "yes") {
		data.Error.Field["enableAttorney"] = map[string]string{
			"reason": "Please select either the attorneys that can be enabled or skip the replacement of the attorneys",
		}
	}
}

func validateManageAttorneysPage(r *http.Request, data *removeAnAttorneyData) {
	if (len(data.Form.DecisionAttorneysUids) == 0 && !postFormCheckboxChecked(r, "skipDecisionAttorney", "yes")) ||
		(len(data.Form.DecisionAttorneysUids) > 0 && postFormCheckboxChecked(r, "skipDecisionAttorney", "yes")) {
		data.Error.Field["decisionAttorney"] = map[string]string{
			"reason": "Select who cannot make joint decisions, or select 'Joint decisions can be made by all attorneys'",
		}
	}
}
