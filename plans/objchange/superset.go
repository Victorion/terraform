package objchange

import (
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/hashicorp/terraform/configs/configschema"
)

// AssertObjectSuperset checks whether the given "actual" value is a valid
// completion of the possibly-partially-unknown "planned" value.
//
// This means that any known leaf value in "planned" must be equal to the
// corresponding value in "actual", and various other similar constraints.
//
// Any inconsistencies are reported by returning a non-zero number of errors.
// These errors are usually (but not necessarily) cty.PathError values
// referring to a particular nested value within the "actual" value.
//
// The two values must have types that conform to the given schema's implied
// type, or this function will panic.
func AssertObjectSuperset(schema *configschema.Block, planned, actual cty.Value) []error {
	return assertObjectSuperset(schema, planned, actual, nil)
}

func assertObjectSuperset(schema *configschema.Block, planned, actual cty.Value, path cty.Path) []error {
	var errs []error
	for name := range schema.Attributes {
		plannedV := planned.GetAttr(name)
		actualV := actual.GetAttr(name)

		moreErrs := assertValueSuperset(plannedV, actualV, path)
		errs = append(errs, moreErrs...)
	}
	// TODO: Check the nested blocks too
	/*
		for name, blockS := range schema.BlockTypes {
		}
	*/
	return errs
}

func assertValueSuperset(planned, actual cty.Value, path cty.Path) []error {
	var errs []error
	if problems := planned.Type().TestConformance(actual.Type()); len(problems) > 0 {
		errs = append(errs, path.NewErrorf("wrong final value type: %s", convert.MismatchMessage(actual.Type(), planned.Type())))
		// If the types don't match then we can't do any other comparisons,
		// so we bail early.
		return errs
	}

	// TODO: Check the two values for conformance

	return errs
}
