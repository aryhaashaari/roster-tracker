package validator

import validation "github.com/go-ozzo/ozzo-validation/v4"

func ExtractMessageToMap(err error) map[string]string {
	result := map[string]string{}
	vErr, ok := err.(validation.Errors)
	if !ok {
		return result
	}

	for k, m := range vErr {
		result[k] = m.Error()
	}

	return result
}

func ExtractMessageToSliceMap(err error) []map[string]string {
	result := []map[string]string{}
	vErr, ok := err.(validation.Errors)
	if !ok {
		return result
	}

	for k, m := range vErr {
		result = append(result, map[string]string{
			k: m.Error(),
		})
	}

	return result
}
