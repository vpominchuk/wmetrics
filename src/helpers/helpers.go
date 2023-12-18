package helpers

import "regexp"

func RegexpStringMatch(pattern string, subject string) (bool, error) {
	re, err := regexp.Compile(pattern)

	if err != nil {
		return false, err
	}

	return re.MatchString(subject), nil
}
