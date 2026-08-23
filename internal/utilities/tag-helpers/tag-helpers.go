package taghelpers

import (
	"errors"
	"reflect"
)

// Parses all fields in an entity looking for a specific tag and then returns a slice of the tag values from all fields
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

//Looks for a field with the sourceTag set to the sourceValue, then looks for the targetTag in that field, and finally returns its value
func GetTagValueBySourceTag[T any](sourceTag, sourceValue, targetTag string) (string, error) {
	var instance T

	fields := reflect.VisibleFields(reflect.TypeOf(instance))
	for i := range fields {
		tagValue, found := fields[i].Tag.Lookup(sourceTag)
		if found && tagValue == sourceValue {
			targetTagValue, found := fields[i].Tag.Lookup(targetTag)
			if !found {
				return "", errors.New("Failed to find the value for the tag " + targetTag + " in struct " + reflect.TypeOf(instance).Name())
			}
			return targetTagValue, nil
		}
	}

	return "", errors.New("Failed to find the target tag value")
}
