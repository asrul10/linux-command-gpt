package gpt

import (
	"testing"
)

func TestApiKey(t *testing.T) {
	gpt3 := Gpt3{
		ApiKeyFile: ".openai_api_key_test",
	}

	tests := []struct {
		homeDir        string
		apiKey         string
		expected       bool
		expectedApiKey string
	}{
		{".", "", false, ""},
		{"./", "", false, ""},
		{".", "the key 123", true, "the key 123"},
		{".", "the key 123\n", true, "the key 123"},
		{".", "    the key 123    ", true, "the key 123"},
		{".", " \n\n   the key 123    \n\n", true, "the key 123"},
	}
	defer gpt3.deleteApiKey()

	for _, test := range tests {
		gpt3.HomeDir = test.homeDir
		gpt3.storeApiKey(test.apiKey)
		load := gpt3.loadApiKey()
		gpt3.deleteApiKey()
		if load != test.expected {
			t.Error("Expected load to be", test.expected, "got", load)
		}
		if gpt3.ApiKey != test.expectedApiKey {
			t.Error("Expected ApiKey to be", test.expectedApiKey, "got", gpt3.ApiKey)
		}
	}

	// Test update api key
	gpt3.HomeDir = "."
	gpt3.storeApiKey("test")
	updateTests := []struct {
		apiKey         string
		expectedApiKey string
	}{
		{"the key 123", "the key 123"},
		{"the key 123\n", "the key 123"},
		{"    the key 123    ", "the key 123"},
		{" \n\n   the key 123    \n\n", "the key 123"},
	}
	for _, test := range updateTests {
		gpt3.updateApiKey(test.apiKey)
		if gpt3.ApiKey != test.expectedApiKey {
			t.Error("Expected ApiKey to be", test.expectedApiKey, "got", gpt3.ApiKey)
		}
	}
}
