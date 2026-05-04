package shared

import "errors"

func GetFromQueryParams(name string, params map[string][]string) (string, error) {
	for key, param := range params {
		if key == name {
			if len(param) == 0 {
				continue
			}
			return param[0], nil
		}
	}

	return "", errors.New("Failed to find the target parameter " + name + " in provided query parameters")
}
