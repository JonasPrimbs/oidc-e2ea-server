/*
 * OIDC² - Identity Certification Token Endpoint
 *
 * Endpoint for OpenID Connect's Identity Certification Token endpoint.
 *
 * API version: 0.5.0
 */
package ict

import (
	"errors"
	"math/big"
	"reflect"
	"strconv"
)

func Int64FromJson(json map[string]interface{}, attributeName string) (int64, error) {
	// Get attribute from json object
	attributeInterface, ok := json[attributeName]
	if !ok {
		return 0, errors.New("attribute '" + attributeName + "' not found")
	}

	// Ensure that attribute value is of type string
	switch attributeType := attributeInterface.(type) {
	case string:
		valueInt, err := strconv.ParseInt(attributeInterface.(string), 10, 64)
		if err != nil {
			return 0, err
		}
		return valueInt, nil
	case int:
		return int64(attributeInterface.(int)), nil
	case int16:
		return int64(attributeInterface.(int16)), nil
	case int32:
		return int64(attributeInterface.(int32)), nil
	case int64:
		return attributeInterface.(int64), nil
	case float32:
		return int64(attributeInterface.(float32)), nil
	case float64:
		return int64(attributeInterface.(float64)), nil
	default:
		return 0, errors.New("attribute '" + attributeName + "' is of type '" + reflect.TypeOf(attributeType).Name() + "' but expected type 'string'")
	}
}

func StringFromJson(json map[string]interface{}, attributeName string) (string, error) {
	// Get attribute from json object
	attributeInterface, ok := json[attributeName]
	if !ok {
		return "", errors.New("attribute '" + attributeName + "' not found")
	}

	// Ensure that attribute value is of type string
	switch attributeType := attributeInterface.(type) {
	case string:
	default:
		return "", errors.New("attribute '" + attributeName + "' is of type '" + reflect.TypeOf(attributeType).Name() + "' but expected type 'string'")
	}

	// Cast attribute value to string and return it
	return attributeInterface.(string), nil
}

func BigIntFromJsonBase64(json map[string]interface{}, attributeName string) (*big.Int, error) {
	// Read string value
	valueString, err := StringFromJson(json, attributeName)
	if err != nil {
		return nil, err
	}

	// Decode base64 string to big integer
	value, err := Base64ToBigInt(valueString)
	if err != nil {
		return nil, errors.New("failed to decode big integer: " + err.Error())
	}

	return value, nil
}

func ByteArrayFromJsonBase64(json map[string]interface{}, attributeName string) ([]byte, error) {
	// Read string value
	valueString, err := StringFromJson(json, attributeName)
	if err != nil {
		return nil, err
	}

	// Decode base64 string to byte array
	value, err := Base64ToByteArray(valueString)
	if err != nil {
		return nil, errors.New("failed to decode byte array: " + err.Error())
	}

	return value, nil
}

func IntFromJsonBase64(json map[string]interface{}, attributeName string) (int, error) {
	// Read string value
	valueString, err := StringFromJson(json, attributeName)
	if err != nil {
		return 0, err
	}

	// Decode base64 string to integer
	value, err := Base64ToInt(valueString)
	if err != nil {
		return 0, errors.New("failed to decode integer: " + err.Error())
	}

	return value, nil
}

func JsonFromJson(json map[string]interface{}, attributeName string) (map[string]interface{}, error) {
	// Get attribute from json object
	attributeInterface, ok := json[attributeName]
	if !ok {
		return nil, errors.New("attribute '" + attributeName + "' not found")
	}

	// Ensure that attribute value is of type json object
	switch attributeType := attributeInterface.(type) {
	case map[string]interface{}:
	default:
		return nil, errors.New("attribute '" + attributeName + "' is of type '" + reflect.TypeOf(attributeType).Name() + "' but expected type 'object'")
	}

	// Cast attribute value to json object and return it
	return attributeInterface.(map[string]interface{}), nil
}
