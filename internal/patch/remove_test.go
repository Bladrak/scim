package patch

import (
	"fmt"
	"github.com/elimity-com/scim/schema"
	"testing"
)

// The following example shows how remove all members of a group.
func Example_removeAllMembers() {
	operation := `{
	"Op": "remove",
	"Path": "members"
}`
	validator, _ := NewValidator(operation, schema.CoreGroupSchema())
	fmt.Println(validator.Validate())
	// Output:
	// <nil> <nil>
}

// The following example shows how remove a value from a complex multi-valued attribute.
func Example_removeComplexMultiValuedAttributeValue() {
	operation := `{
	"Op": "remove",
	"Path": "emails[type eq \"work\" and value ew \"elimity.com\"]"
}`
	validator, _ := NewValidator(operation, schema.CoreUserSchema())
	fmt.Println(validator.Validate())
	// Output:
	// <nil> <nil>
}

// The following example shows how remove a single member from a group.
func Example_removeSingleMember() {
	operation := `{
	"Op": "remove",
	"Path": "members[value eq \"0001\"]"
}`
	validator, _ := NewValidator(operation, schema.CoreGroupSchema())
	fmt.Println(validator.Validate())
	// Output:
	// <nil> <nil>
}

// The following example shows how to replace all of the members of a group with a different members list.
func Example_replaceAllMembers() {
	operations := []string{`{
	"Op": "remove",
	"Path": "members"
}`, `{
	"Op": "add",
	"Path": "members",
	"value": [
		{
			"display": "di-wu",
			"$ref": "https://example.com/v2/Users/0001",
			"value": "0001"
		},
		{
			"display": "example",
			"$ref": "https://example.com/v2/Users/0002",
			"value": "0002"
		}
	]
}`,
	}
	for _, op := range operations {
		validator, _ := NewValidator(op, schema.CoreGroupSchema())
		fmt.Println(validator.Validate())
	}
	// Output:
	// <nil> <nil>
	// [map[$ref:https://example.com/v2/Users/0001 display:di-wu type:<nil> value:0001] map[$ref:https://example.com/v2/Users/0002 display:example type:<nil> value:0002]] <nil>
}

func TestOperationValidator_ValidateRemove(t *testing.T) {
	// The goal this test is to cover Section 3.5.2.2 of RFC7644.
	// More info: https://tools.ietf.org/html/rfc7644#section-3.5.2.2

	// Some more indirect things are covered by these tests:
	// - If the target location is a single-value attribute, the attribute and its associated value is removed.
	// - If the target location is a multi-valued attribute and no filter is specified, the attribute and all values
	//   are removed.
	// - If the target location is a multi-valued attribute and a complex filter is specified comparing a "value", the
	//   values matched by the filter are removed.
	// - If the target location is a complex multi-valued attribute and a complex filter is specified based on the
	//   attribute's sub-attributes, the matching records are removed.

	for i, test := range []struct {
		valid   string
		invalid string
	}{
		// If "Path" is unspecified, the operation fails.
		{invalid: `{"Op":"remove"}`},

		// If the target location is a single-value attribute.
		{valid: `{"Op":"remove","Path":"attr1"}`},
		// If the target location is a multi-valued attribute and no filter is specified.
		{valid: `{"Op":"remove","Path":"multiValued"}`},
		// If the target location is a multi-valued attribute and a complex filter is specified comparing a "value".
		{valid: `{"Op":"remove","Path":"multivalued[value eq \"value\"]"}`},
		// If the target location is a complex multi-valued attribute and a complex filter is specified based on the
		// attribute's sub-attributes
		{valid: `{"Op":"remove","Path":"complexMultiValued[attr1 eq \"value\"]"}`},
		{valid: `{"Op":"remove","Path":"complexMultiValued[attr1 eq \"value\"].attr1"}`},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// valid
			if op := test.valid; op != "" {
				validator, err := NewValidator(op, patchSchema, patchSchemaExtension)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := validator.Validate(); err != nil {
					t.Errorf("The following operatation should be an VALID add operation:\n(case %d): %s\n%v", i, op, err)
				}
			}
			// invalid
			if op := test.invalid; op != "" {
				validator, err := NewValidator(op, patchSchema, patchSchemaExtension)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := validator.Validate(); err == nil {
					t.Errorf("The following operatation should be an INVALID add operation:\n(case %d): %s", i, op)
				}
			}
		})
	}
}