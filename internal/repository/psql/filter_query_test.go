package psql

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/tahmooress/store/internal/model"
)

var testCases = []struct {
	input         *model.Filter
	expectedQuery string
	expectedArgs  []interface{}
}{
	{
		input: model.NewFilter("123456").SetName("someName").SetTags([]string{"important", "myTag"}),
		expectedQuery: `
		SELECT owner_id, file_object.name, type, storage_location, created_at FROM
		file_object INNER JOIN  
		file_tag ft ON  file_object.id = file_tag.file_object_id 
		INNER JOIN tag ON file_tag.tag_id = tag.id 
		WHERE owner_id = $1 AND (file_object.name = $2 OR tag.name IN ($3,$4))
	`,
		expectedArgs: []interface{}{"123456", "someName", "important", "myTag"},
	},
	{
		input: model.NewFilter("123456").SetTags([]string{"important", "myTag"}),
		expectedQuery: `
		SELECT owner_id, file_object.name, type, storage_location, created_at FROM
		file_object INNER JOIN  
		file_tag ft ON  file_object.id = file_tag.file_object_id 
		INNER JOIN tag ON file_tag.tag_id = tag.id 
		WHERE owner_id = $1 AND tag.name IN ($2,$3)
		`,
		expectedArgs: []interface{}{"123456", "important", "myTag"},
	},
	{
		input: model.NewFilter("123456").SetName("myName"),
		expectedQuery: `
		SELECT owner_id, file_object.name, type, storage_location, created_at FROM
		file_object INNER JOIN  
		file_tag ft ON  file_object.id = file_tag.file_object_id 
		INNER JOIN tag ON file_tag.tag_id = tag.id
		WHERE owner_id = $1 AND file_object.name = $2
		`,
		expectedArgs: []interface{}{"123456", "myName"},
	},
	{
		input: model.NewFilter("123456"),
		expectedQuery: `
		SELECT owner_id, file_object.name, type, storage_location, created_at FROM
		file_object INNER JOIN  
		file_tag ft ON  file_object.id = file_tag.file_object_id 
		INNER JOIN tag ON file_tag.tag_id = tag.id
		WHERE owner_id = $1
		`,
		expectedArgs: []interface{}{"123456"},
	},
}

func TestFetchObjectfilterQuery(t *testing.T) {
	for _, tc := range testCases {
		query, args, err := tc.input.Query(fetchObjectfilterQuery)
		if err != nil {
			t.Errorf("query fail for case: %v", *tc.input)
		}

		fmt.Println(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(query, " ", ""), "\n", ""), "\t", ""))

		if strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					query, " ", ""), "\n", ""), "\t", "") !=
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(
						tc.expectedQuery, " ", ""), "\n", ""), "\t", "") {
			t.Errorf("query are not equal got: %s expected: %s", query, tc.expectedQuery)
		}

		for i, arg := range args {
			if arg != tc.expectedArgs[i] {
				t.Errorf("args are not equal got: %v expected: %v", arg, tc.expectedArgs[i])
			}
		}
	}
}

var ParseParenthesesTestCases = []struct {
	input    []string
	expected error
}{
	{
		input:    []string{"(", "(", ")", ")"},
		expected: nil,
	},
	{
		input:    []string{"(", "(", "(", ")", ")"},
		expected: ErrInvalidParentheses,
	},
	{
		input:    []string{"(", ")", ")", ")"},
		expected: ErrInvalidParentheses,
	},
	{
		input:    []string{")", "(", ")", "("},
		expected: ErrInvalidParentheses,
	},
	{
		input:    []string{"(", "(", ")", ")", "(", ")"},
		expected: nil,
	},
	{
		input:    []string{"(", ")", "(", ")", "(", ")"},
		expected: nil,
	},
}

func TestParseParentheses(t *testing.T) {
	for _, v := range ParseParenthesesTestCases {
		err := parseParentheses(v.input)
		if !errors.Is(err, v.expected) {
			t.Errorf("parseParentheses got: %s expected: %s", err, v.expected)
		}
	}
}
