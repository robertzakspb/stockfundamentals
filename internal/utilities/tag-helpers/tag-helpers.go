package taghelpers

import (
	"errors"
	"reflect"
)

func GetEntityTagValues[T any](tag string) ([]string, error) {
	tags := []string{}
	var instance T

	fields := reflect.VisibleFields(reflect.TypeOf(instance))
	for i := range fields {
		tagValue, found := fields[i].Tag.Lookup(tag)
		if found {
			tags = append(tags, tagValue)
		}
	}

	if len(tags) == 0 {
		return tags, errors.New("Found zero values for tag " + tag + " in entity " + reflect.TypeOf(instance).Name())
	}

	return tags, nil
}
