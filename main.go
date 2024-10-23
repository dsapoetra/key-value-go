package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// AttributeType represents the possible data types for attribute values
type AttributeType int

const (
	StringType AttributeType = iota
	FloatType
	BoolType
)

// AttributeMetadata stores the data type for an attribute
type AttributeMetadata struct {
	dataType AttributeType
}

// Store represents the thread-safe key-value store
type Store struct {
	data           map[string]map[string]interface{}
	attributeTypes map[string]AttributeMetadata
	mutex          sync.RWMutex
}

// NewStore creates a new instance of the key-value store
func NewStore() *Store {
	return &Store{
		data:           make(map[string]map[string]interface{}),
		attributeTypes: make(map[string]AttributeMetadata),
	}
}

// determineType returns the AttributeType for a given string value
func determineType(value string) (AttributeType, interface{}, error) {
	// Try boolean first
	if value == "true" || value == "false" {
		return BoolType, value == "true", nil
	}

	// Try float (will also handle integers)
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return FloatType, floatVal, nil
	}

	// Default to string
	return StringType, value, nil
}

// Put adds or updates a key-value pair in the store
func (s *Store) Put(key string, attributes [][]string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newData := make(map[string]interface{})

	for _, attr := range attributes {
		attrKey := attr[0]
		attrValue := attr[1]

		valueType, parsedValue, err := determineType(attrValue)
		if err != nil {
			return err
		}

		if metadata, exists := s.attributeTypes[attrKey]; exists {
			if metadata.dataType != valueType {
				return errors.New("Data Type Error")
			}
		} else {
			s.attributeTypes[attrKey] = AttributeMetadata{dataType: valueType}
		}

		newData[attrKey] = parsedValue
	}

	s.data[key] = newData
	return nil
}

// Get retrieves a value from the store
func (s *Store) Get(key string) map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if value, exists := s.data[key]; exists {
		return value
	}
	return nil
}

// Delete removes a key-value pair from the store
func (s *Store) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, key)
}

// Search finds all keys that have the given attribute key-value pair
func (s *Store) Search(attrKey, attrValue string) []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var results []string
	_, expectedValue, _ := determineType(attrValue)

	for key, attributes := range s.data {
		if value, exists := attributes[attrKey]; exists {
			if fmt.Sprintf("%v", value) == fmt.Sprintf("%v", expectedValue) {
				results = append(results, key)
			}
		}
	}

	sort.Strings(results)
	return results
}

// Keys returns all keys in the store
func (s *Store) Keys() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func formatValue(value interface{}) string {
	if floatVal, ok := value.(float64); ok {
		// For float values, check if they're whole numbers
		if floatVal == float64(int(floatVal)) {
			return fmt.Sprintf("%.1f", floatVal) // Always show one decimal place
		}
		return fmt.Sprintf("%.2f", floatVal) // Show two decimal places
	}
	return fmt.Sprintf("%v", value)
}

func main() {
	store := NewStore()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please input command and param")
	for scanner.Scan() {

		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "put":
			if len(parts) < 4 || len(parts)%2 != 0 {
				fmt.Println("Number of params is incorrect")
				continue
			}
			key := parts[1]
			var attributes [][]string
			for i := 2; i < len(parts); i += 2 {
				attributes = append(attributes, []string{parts[i], parts[i+1]})
			}
			if err := store.Put(key, attributes); err != nil {
				fmt.Println("Data Type Error")
			}
			fmt.Println("Put is done")

		case "get":
			if len(parts) != 2 {
				continue
			}
			key := parts[1]
			value := store.Get(key)
			if value == nil {
				fmt.Printf("No entry found for %s\n", key)
				continue
			}

			keys := make([]string, 0, len(value))
			for k := range value {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			var output []string
			for _, k := range keys {
				output = append(output, fmt.Sprintf("%s: %v", k, formatValue(value[k])))
			}
			fmt.Println(strings.Join(output, ", "))

		case "delete":
			if len(parts) != 2 {
				continue
			}
			key := parts[1]
			store.Delete(key)
			fmt.Println("Delete is done")

		case "search":
			if len(parts) != 3 {
				continue
			}
			attrKey, attrValue := parts[1], parts[2]
			results := store.Search(attrKey, attrValue)
			if len(results) > 0 {
				fmt.Println(strings.Join(results, ","))
			} else {
				fmt.Println()
			}

		case "keys":
			keys := store.Keys()
			if len(keys) > 0 {
				fmt.Println(strings.Join(keys, ","))
			} else {
				fmt.Println()
			}

		case "exit":
			return
		}
	}
}
