package exporter

import (
	"IG-Parser/parser"
	"IG-Parser/tree"
	"fmt"
	"io/ioutil"
	"testing"
)

/*
Tests the header generation function for tabular output, here as variant that does not assume component aggregation.
Note that this test is implicitly contained in IGStatementParser_test.go
*/
func TestHeaderRowGenerationWithoutComponentAggregation(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// This is tested in IGStatementParser_test.go as well
	nodeArray, componentIdx := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	if nodeArray == nil || componentIdx == nil {
		t.Fatal("Generated array or component header array should not be empty.")
	}

	fmt.Println(componentIdx)

	if componentIdx[tree.ATTRIBUTES] != 1 || componentIdx[tree.DIRECT_OBJECT] != 1 ||
		componentIdx[tree.EXECUTION_CONSTRAINT] != 2 || componentIdx[tree.DEONTIC] != 1 ||
		componentIdx[tree.AIM] != 2 {
		t.Fatal("Component element count is incorrect.")
	}

	if componentIdx[tree.CONSTITUTED_ENTITY] != 0 || componentIdx[tree.CONSTITUTED_ENTITY_PROPERTY] != 0 ||
		componentIdx[tree.CONSTITUTIVE_FUNCTION] != 0 || componentIdx[tree.CONSTITUTING_PROPERTIES] != 0 ||
		componentIdx[tree.CONSTITUTING_PROPERTIES_PROPERTY] != 0 {
		t.Fatal("Component element count is not empty for some elements.")
	}

	output := "Init, "
	output, outArr, _, err := generateHeaderRow(output, componentIdx, ";")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during header generation:", err)
	}

	// Unnecessary check is intentional to allow for flexible row specification change
	if output == "" || len(outArr) == 0 {
		t.Fatal("Generate header row did not return filled data structures")
	}

	if output != "Init, Attributes;Deontic;Aim_1;Aim_2;Direct Object;Execution Constraint_1;Execution Constraint_2" {
		t.Fatal("Generated header row is wrong. Output: " + output)
	}

	if fmt.Sprint(outArr) != "[A D I_1 I_2 Bdir Cex_1 Cex_2]" {
		t.Fatal("Generated component array is wrong: " + fmt.Sprint(outArr))
	}

}

/*
Tests the header generation function for tabular output, here as variant that applies component aggregation.
Note that this test is implicitly contained in IGStatementParser_test.go
*/
func TestHeaderRowGenerationWithComponentAggregation(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// OVERRIDE dynamic output setting
	tree.AGGREGATE_IMPLICIT_LINKAGES = true

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// This is tested in IGStatementParser_test.go as well
	nodeArray, componentIdx := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	if nodeArray == nil || componentIdx == nil {
		t.Fatal("Generated array or component header array should not be empty.")
	}

	fmt.Println(componentIdx)

	if componentIdx[tree.ATTRIBUTES] != 1 || componentIdx[tree.DIRECT_OBJECT] != 1 ||
		componentIdx[tree.EXECUTION_CONSTRAINT] != 1 || componentIdx[tree.DEONTIC] != 1 ||
		componentIdx[tree.AIM] != 1 {
		t.Fatal("Component element count is incorrect.")
	}

	if componentIdx[tree.CONSTITUTED_ENTITY] != 0 || componentIdx[tree.CONSTITUTED_ENTITY_PROPERTY] != 0 ||
		componentIdx[tree.CONSTITUTIVE_FUNCTION] != 0 || componentIdx[tree.CONSTITUTING_PROPERTIES] != 0 ||
		componentIdx[tree.CONSTITUTING_PROPERTIES_PROPERTY] != 0 {
		t.Fatal("Component element count is not empty for some elements.")
	}

	output := "Init, "
	output, outArr, _, err := generateHeaderRow(output, componentIdx, ";")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during header generation:", err)
	}

	// Unnecessary check is intentional to allow for flexible row specification change
	if output == "" || len(outArr) == 0 {
		t.Fatal("Generate header row did not return filled data structures")
	}

	if output != "Init, Attributes;Deontic;Aim;Direct Object;Execution Constraint" {
		t.Fatal("Generated header row is wrong. Output: " + output)
	}

	if fmt.Sprint(outArr) != "[A D I Bdir Cex]" {
		t.Fatal("Generated component array is wrong: " + fmt.Sprint(outArr))
	}

}

/*
Tests simple tabular output without any combinations or nesting.
*/
func TestSimpleTabularOutput(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect)" +
		"Bdir(certified production facilities) "

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false
	// OVERRIDE dynamic output setting
	tree.AGGREGATE_IMPLICIT_LINKAGES = true

	// Override cell separator symbol
	CellSeparator = ";"

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 5 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputSimpleTabularOutput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Take separator for Google Sheets output
	separator := ";"

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting, but component-level combinations; no implicit combinations
*/
func TestBasicTabularOutputCombinations(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) "

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 5 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputSimpleNoNestingWithCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting, but implicitly sAND-linked components
*/
func TestBasicTabularOutputImplicitAnd(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(review) I(sustain) " +
		"Bdir(certified production facilities) "

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 6 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputSimpleNoNestingImplicitAnd.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting, but component-level combinations and implicitly linked components
*/
func TestTabularOutputCombinationsImplicitAnd(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputSimpleNoNestingWithCombinationsImplicitAnd.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting - only component-level combinations,
but including shared left elements in output.
*/
func TestTabularOutputWithSharedLeftElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputSimpleNoNestingWithSharedLeftElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting - only component-level combinations,
but including shared right elements in output.
*/
func TestTabularOutputWithSharedRightElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I((review [AND] (refresh [AND] drink)) rightShared) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputSimpleNoNestingWithSharedRightElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting - only component-level combinations,
but including shared left and right elements in output.
*/
func TestTabularOutputWithSharedLeftAndRightElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(leftShared (review [AND] (refresh [AND] revise)) rightShared) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputSimpleNoNestingWithSharedLeftAndRightElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests multi-level nesting on statements, i.e., activation with own activation condition
*/
func TestTabularOutputWithTwoLevelNestedComponent(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// This is the tricky lines, specifically the second Cac{}
		"Cac{E(Program Manager) F(is) P((approved [AND] committed)) Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputTwoLevelComplexNesting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests tabular output with combination of two-level statement-level nested component and
simple activation condition (to be linked by implicit AND).
*/
func TestTabularOutputWithCombinationOfSimpleAndTwoLevelNestedComponent(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Simple activation condition
		"Cac(Upon approval)" +
		// Complex activation condition, including two-level nesting (Cac{Cac{}})
		"Cac{E(Program Manager) F(is) P((approved [AND] committed)) Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 9 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputSimpleAndTwoLevelComplexNesting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests combination of two nested activation conditions (single level)
*/
func TestTabularOutputWithCombinationOfTwoNestedComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputTwoNestedComplexComponents.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests combination of three nested activation conditions (single level), including embedded component-level nesting
*/
func TestTabularOutputWithCombinationOfThreeNestedComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}" +
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputThreeNestedComplexComponents.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests partial AND-linked statement-level combinations, expects the inference of implicit AND to non-linked complex component
*/
func TestTabularOutputWithNestedStatementCombinationsImplicitAnd(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Proper combination
		"{Cac{E(Program Manager) F(is) P(approved)} [OR] " +
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Implicitly linked nested statement
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputNestedComplexCombinationsImplicitAnd.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests flattened representation (IG Core) of partial AND-linked statement-level combinations, expects the inference of implicit AND to non-linked complex component
*/
func TestTabularOutputWithNestedStatementCombinationsImplicitAndIGCore(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Proper combination
		"{Cac{E(Program Manager) F(is) P(approved)} [OR] " +
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Implicitly linked nested statement
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Core output
	SetProduceIGExtendedOutput(false)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputNestedComplexCombinationsImplicitAndIGCore.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests partial XOR-linked statement-level combinations, expects the inference of implicit AND to non-linked complex component
*/
func TestTabularOutputWithNestedStatementCombinationsXOR(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Complex expression with XOR linkage
		"{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// should be automatically linked using AND
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputNestedComplexCombinationsXor.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement-level combinations, alongside embedded component-level combinations to ensure the
filtering of within-statement component-level combinations are filtered prior to statement assembly.
*/
func TestTabularOutputWithNestedStatementCombinationsAndComponentCombinations(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputNestedStatementCombinationsAndComponentCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement-level combinations, alongside embedded component-level combinations to ensure the
filtering of within-statement component-level combinations are filtered prior to statement assembly.
Includes generation of output with shared elements.
*/
func TestTabularOutputWithNestedStatementCombinationsAndComponentCombinationsWithSharedElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputNestedStatementCombinationsAndComponentCombinationsWithSharedElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement-level combinations on multiple components, alongside selected embedded component-level,
alongside combination with non-nested components
*/
func TestTabularOutputWithMultipleNestedStatementsAndSimpleComponentsAcrossDifferentComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		// This is another nested component - should have implicit link to regular Bdir
		"Bdir{A(farmers) that I((apply [OR] plan to apply)) for Bdir(organic farming status)}" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Generated Component References: ", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 9 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputMultipleNestedStatementsAndSimpleComponentsAcrossDifferentComponents.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement with multiple levels of component-level nesting, alongside component combination.
Note: Generally initially identified as combination, but recategorized as regular nested statement.
Exemplifies exception handling of parser.
*/
func TestStaticTabularOutputWithMultiLevelNestingAndComponentLevelCombinations(t *testing.T) {

	text := "Cac{When A(Program Manager) I(reveals)	Bdir{A,p(accredited) A(certifying agent) I([is not in compliance]) " +
		// Parser miscategorizes as combination based on nested logical operators, but should recover
		"with the Bdir((Act [OR] regulations in this part))} " +
		"Cac{When A(Program Manager) I((([inspects] [OR] [reviews]) [OR] [investigates])) " +
		"Bind,p(accredited) Bind(certifying agent)}} " +
		"A([Program Manager]) D(shall) I([send]) Bdir(notification) Bdir,p(of non-compliance) to the " +
		"Bind,p(accredited) Bind(certifying agent)."

	// Dynamic output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := "|"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Generated Component References: ", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaMultiLevelNestingAndComponentLevelCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests or else statements with statement-level annotations.
*/
func TestStaticTabularOutputOrElseAnnotations(t *testing.T) {

	text := "A(farmer) D([must]) I(submit) Bdir,p(an organic systems) Bdir(plan) Cex(by the end of the " +
		"calendar year) O[consequence]{the A(certifier) D(may) I(suspend) the Bdir,p(farmer’s) Bdir(operating license)}"

	// Dynamic output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Activate annotations
	SetIncludeAnnotations(true)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Take separator for Google Sheets output
	separator := "|"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Generated Component References: ", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticOrElseAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement-level combinations with incompatible component symbols
*/
func TestTabularOutputWithStatementCombinationsOfIncompatibleComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This combination mixes the previous Cac with Cex - and should fail
		"Cex{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	_, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_INVALID_TYPES_IN_NESTED_STATEMENT_COMBINATION {
		t.Fatal("Parser should have picked up on invalid component combinations.")
	}

}

/*
Tests combination of two nested activation conditions (single level) for static output
*/
func TestStaticTabularOutputBasicStatement(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(approved) Bdir,p(certified) Bdir((production [operations] [AND] handling operations)) " + //Bdir,p1(accredited) Bdir1(certifying agents) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaBasicStatement.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests nested combinations with focus on proper resolution of references.
*/
func TestStaticTabularOutputNestedCombinations(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(approved) Bdir,p(certified) Bdir((production [operations] [AND] handling operations)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"{Cac{E(Program Manager) F(is) P((approved [AND] committed))} [XOR] " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaNestedCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests nested combinations and implicit AND-linked statement with focus on proper resolution of references.
*/
func TestStaticTabularOutputNestedCombinationsImplicitAnd(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(approved) Bdir,p(certified) Bdir((production [operations] [AND] handling operations)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"{Cac{E(Program Manager) F(is) P((approved [AND] committed))} [XOR] " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Activation condition 3 (to be linked implicitly)
		"Cac{A(Another official) I(does) Bdir(something else)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaNestedCombinationsImplicitAnd.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests flat printing (IG Core) of nested combinations and implicit AND-linked statement with focus on proper resolution of references.
*/
func TestStaticTabularOutputNestedCombinationsImplicitAndIGCore(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(approved) Bdir,p(certified) Bdir((production [operations] [AND] handling operations)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"{Cac{E(Program Manager) F(is) P((approved [AND] committed))} [XOR] " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Activation condition 3 (to be linked implicitly)
		"Cac{A(Another official) I(does) Bdir(something else)}"

	// Static output
	SetDynamicOutput(false)
	// IG Core output (no component-level nesting)
	SetProduceIGExtendedOutput(false)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaNestedCombinationsImplicitAndIGCore.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests printing of nested properties and proper resolution of references.
*/
func TestStaticTabularOutputNestedProperties(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(approved) Bdir,p(certified) Bdir,p{E(operations) that F(have experience) with P(farming)} " +
		"Bdir((production [operations] [AND] handling operations)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"{Cac{E(Program Manager) F(is) P((approved [AND] committed))} [XOR] " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Activation condition 3 (to be linked implicitly)
		"Cac{A(Another official) I(does) Bdir(something else)}"

	// Static output
	SetDynamicOutput(false)
	// IG Core output (no component-level nesting)
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 9 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaNestedProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Correct parsing of shared elements left and right (on Aim and Cex).
*/
func TestStaticTabularOutputBasicStatementSharedLeftAndRightElements(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink)) rightShared) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part) and beyond) "

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 5 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaBasicStatementLeftAndRightElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Correct parsing of private properties (mix of shared and private) for static output
*/
func TestStaticTabularOutputBasicStatementSharedAndPrivateProperties(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p1(accredited) Bdir1(certifying agents) Bdir(other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaBasicStatementPrivateProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests combination of two nested activation conditions (single level) for static output and tests for a mix
of shared and private properties (on top level) and private properties only on nested level
*/
func TestStaticTabularOutputBasicStatementMixSharedPrivateAndNestedPrivatePropertiesOnly(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p1(accredited) Bdir1(certifying agents) Bdir(other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2 -- note that associated is wrongly annotated, leading to linkage to both Bdirs
		"Cac{A(NOP Official) I(recognizes) Bdir1,p1(responsible) Bdir1(Program Manager) and Bdir,p2(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaStatementSharedAndPrivateOnlyProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests presence of shared and private properties on indexed components (e.g., A1(content)), as opposed
to indexed properties only (e.g., A,p1(content)).
*/
func TestStaticTabularOutputBasicStatementComponentLevelIndexedProperties(t *testing.T) {

	text := "A1,p(National Organic Program's) A1(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p(accredited) Bdir1(certifying agents) Bdir(other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2 - shared properties implied in this activation condition due to wrong syntax
		"Cac{A(NOP Official) I(recognizes) Bdir,p1(responsible) Bdir1(Program Manager) and Bdir,p2(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)

	// Take separator for Google Sheets output
	separator := ";"

	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaStatementComponentLevelIndexedProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests presence of shared and private properties on indexed components (e.g., A1(content)), as opposed
to indexed properties only (e.g., A,p1(content)).
*/
func TestStaticTabularOutputBasicStatementComponentLevelIndexedPropertiesAnnotations(t *testing.T) {

	text := "A1,p[prop=qualitative](National Organic Program's) A1[type=animate](Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p(accredited) Bdir1(certifying agents) Bdir(other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2 - shared properties implied in this activation condition
		"Cac{A(NOP Official) I(recognizes) Bdir,p1(responsible) Bdir1(Program Manager) and Bdir,p2(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaStatementComponentLevelIndexedProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests combination of two nested activation conditions (single level) , tests for a mix
of shared and private properties (on top level) and private properties only on nested level,
and includes *deactivated* annotations on various components.
*/
func TestStaticTabularOutputBasicStatementMixedPropertiesAnnotationsDeactivated(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p1[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		// Activation condition 2 - misspelt second property is intentional (to make it shared property)
		"Cac{A(NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1(Program Manager) and Bdir,p1(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)

	// Take separator for Google Sheets output
	separator := ";"

	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaStatementSharedAndPrivateOnlyProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests combination of two nested activation conditions (single level) , tests for a mix
of shared and private properties (on top level) and private properties only on nested level,
and includes *activated* annotations on various components.
*/
func TestStaticTabularOutputBasicStatementMixedPropertiesAnnotationsActivated(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A[type=enforcer](NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1[type=main object](Program Manager) and Bdir2,p2[type=third party](associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)

	// Take separator for Google Sheets output
	separator := ";"

	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaStatementMixedPropertiesAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests for proper replacement of special symbols for preprocessing Google Sheets output. Includes
complexity of previous tests.
*/
func TestStaticTabularOutputBasicStatementSpecialSymbols(t *testing.T) {

	text := "A,p(National Organic Program's) A(\"Program Manager\"), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p(accredited) Bdir1(\"certifying agents) Bdir(\"other agents\")" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A('NOP Official') I(recognizes) Bdir1,p(responsible) Bdir1(Program Manager) and Bdir2,p(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaStatementSpecialSymbolsGoogleSheets.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper export of annotations on nested components.
*/
func TestStaticTabularOutputNestedStatementsAnnotations(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac1[ctx=stAte]{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		// Activation condition 2
		"CacB[annotation2]{A[type=enforcer](NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1[type=main object](Program Manager) and Bdir2,p2[type=third party](associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)

	// Take separator for Google Sheets output
	separator := ";"

	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaNestedStatementsAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper export of annotations on nested component combinations.
*/
func TestStaticTabularOutputNestedStatementCombinationAnnotations(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"{Cac1[ctx=stAte]{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		"[XOR] " +
		// Activation condition 2
		"CacB[annotation2]{A[type=enforcer](NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1[type=main object](Program Manager) and Bdir2,p2[type=third party](associated) Bdir2(inspectors)}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)

	// Take separator for Google Sheets output
	separator := ";"

	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaNestedStatementCombinationsAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper export of annotations on nested component combinations and individual nested statements.
*/
func TestStaticTabularOutputNestedStatementsAndCombinationMixAnnotations(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		// Implicitly linked nested direct object with annotations and invalid structure (no aim and context)
		"Bdir{A[type=animate](another actor) A,p[prop=qualitative](who does not comply)} " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"{Cac1[ctx=stAte]{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		"[XOR] " +
		// Activation condition 2
		"CacB[annotation2]{A[type=enforcer](NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1[type=main object](Program Manager) and Bdir2,p2[type=third party](associated) Bdir2(inspectors)}} " +
		// Implicitly linked activation condition with diverse annotations
		"CacC[ABdir]{A[type=animate](further entity) I[act=violate](violates) Bdir[entity=law](part of provisions)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)

	// Take separator for Google Sheets output
	separator := ";"

	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 9 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaNestedStatementsAndCombinationMixAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper parsing and export of within-component combination linkages (wAND).
*/
func TestStaticTabularOutputParsingOfWithinComponentLinkages(t *testing.T) {

	// Simple output with multiple combinations within Cex
	text := "Cex[exampleConstraint](for compliance with (left [AND] right) as well as (left1 [XOR] right1) shared) Cex(outlier)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)

	// Take separator for Google Sheets output
	separator := "|"

	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 1 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaBasicWithinComponentCombinationLinkages.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests the generation of basic tree output for visual output.
*/
func TestVisualOutputBasic(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Deactivate annotations
	SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualTreeOutputBasic.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the generation of tree output with nested properties for visual output.
*/
func TestVisualOutputNestedProperties(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Bdir,p{E(operation) F(has been vetted) Cex(before certification)} " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Deactivate annotations
	SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualTreeOutputNestedProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the generation of tree output for visual output, including annotations.
*/
func TestVisualOutputAnnotations(t *testing.T) {

	text := "A[gov=enforcer,anim=animate](National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=monitor](inspect and), I[act=enforce](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir[gov=monitored,anim=animate](approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex[ref=(Act,part)](for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E[gov=enforcer,anim=animate](Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A[gov=monitor,anim=animate](NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Activate annotations
	SetIncludeAnnotations(true)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualTreeOutputAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests proper output of proper linkages of complex private properties alongside shared properties for visual output.
Tests flat output for properties (labels, as opposed to tree structure).
*/
func TestVisualOutputPropertyNodesFlatPrinting(t *testing.T) {
	text := "A(General Manager) A,p(shared quality) A1(Region Manager) A1,p(left quality) A1,p(right quality) A1,p(third quality)"

	// Activate annotations
	SetIncludeAnnotations(true)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualTreeComplexPrivateNodesFlat.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper output of proper linkages of complex private properties alongside shared properties for visual output.
Tests tree structure output for properties.
*/
func TestVisualOutputPropertyNodesTreePrinting(t *testing.T) {
	text := "A(General Manager) A,p(shared quality) A1(Region Manager) A1,p(left quality) A1,p(right quality) A1,p(third quality)"

	// Activate annotations
	SetIncludeAnnotations(true)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualTreeComplexPrivateNodesTree.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests the generation of basic tree output for visual output, but as non-binary tree (i.e., collapsing entries associated with same logical operator for given component).
*/
func TestVisualOutputBasicNonBinaryTree(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Deactivate annotations
	SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualTreeOutputBasicNonBinary.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the generation of complex tree output for visual output as non-binary tree (i.e., collapsing entries associated with same logical operator for given component).
Does not decompose property trees.
*/
func TestVisualOutputComplexNonBinaryTree(t *testing.T) {

	// Complex entry
	text := " The Congress finds and declares that it is the E(national policy) F([is] to (encourage [AND] assist)) the P(states) Cex{ A(states) I(to exercise) Cex(effectively) their Bdir(responsibilities) Bdir,p(in the coastal zone) Cex(through the (development [AND] implementation) of management programs to achieve wise use of the (land [AND] water) resources of the coastal zone, giving full consideration to (ecological [AND] cultural [AND] historic [AND] esthetic) values as well as the needs for compatible economic development), Cex{which E(programs) M(should) Cex(at least) F(provide for)— (A) the P1(protection) P,p1(of natural resources, including (wetlands [AND] floodplains [AND] estuaries [AND] beaches [AND] dunes [AND] barrier islands [AND] coral reefs [AND] fish and wildlife and their habitat) within the coastal zone), the P2(management) P,p2((of coastal development to minimize the loss of (life [AND] property) caused by improper development in (flood-prone [AND] storm surge [AND] geological hazard [AND] erosion-prone) areas [AND] in areas likely to be (affected by [OR] vulnerable to) (sea level rise [AND] land subsidence [AND] saltwater intrusion) [AND] by the destruction of natural protective features such as (beaches [AND] dunes [AND] wetlands [AND] barrier islands))), (C) the P3(management) P3,p(of coastal development to (improve [AND] safeguard [AND] restore) the quality of coastal waters, [AND] to protect (natural resources [AND] existing uses of those waters)), (D) P4,p1(priority) P4(consideration) P4,p2(being given to (coastal-dependent (uses [AND] orderly processes) for siting major facilities related to (national defense [AND] energy [AND] fisheries development [AND] recreation [AND] ports [AND] transportation), [AND] the location to the maximum extent practicable of new (commercial [AND] industrial) developments (in [XOR] adjacent) to areas where such development already exists)), (E) P5,p1(public) P5(access) P5,p2(to the coasts for recreation purposes), (F) P6(assistance) P6,p(in the redevelopment of (deteriorating urban (waterfronts [AND] ports) [AND] sensitive (preservation [AND] restoration) of (historic [AND] cultural [AND] esthetic) coastal features)), (G) P7(the (coordination [AND] simplification) of procedures) P7,p1(in order to ensure expedited governmental decision making for the management of coastal resources), (H) P8((continued (consultation [AND] coordination) with, [AND] the giving of adequate consideration to the views of affected Federal agencies)), (I) P9(the giving of ((timely [AND] effective) notification of , [AND] opportunities for (public [AND] local) government participation in coastal management decision making)), (J) P10(assistance) P10,p1(to support comprehensive (planning [AND] conservation [AND] management) for living marine resources) P10,p1,p1(including planning for (the siting of (pollution control [AND] aquaculture facilities) within the coastal zone [AND] improved coordination between ((State [AND] Federal) coastal zone management agencies [AND] (State [AND] wildlife) agencies))) }}"

	// Deactivate annotations
	SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualTreeOutputComplexNonBinary.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests multiple combinations embedded in second level phrase/side of first-order statement.
*/
func TestMultiLevelEmbeddedCombinations(t *testing.T) {

	// Entry containing multiple combinations on right side of first-order combination
	text := "Cex(( left1 [XOR] shared (left [AND] right) via (left2 [XOR] right2)))"

	// Deactivate annotations
	SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualMultiCombinationsPhrase.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests tabular output for the default example statement that showcases most IG features.
*/
func TestTabularOutputDefaultExample(t *testing.T) {

	// Default example
	text := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that {Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate](were (non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)}}."

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)

	// Take separator for Google Sheets output
	separator := ";"

	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputTabularDefaultExample.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for the default example statement that showcases most IG features.
*/
func TestVisualOutputDefaultExample(t *testing.T) {

	// Default example
	text := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that {Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate](were (non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)}}."

	// Activate annotations
	SetIncludeAnnotations(true)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualDefaultExample.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for the default example statement that showcases most IG features, but moving activation conditions to the top.
*/
func TestVisualOutputDefaultExampleActivationConditionsFirst(t *testing.T) {

	// Default example
	text := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that {Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate](were (non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)}}."

	// Activate annotations
	SetIncludeAnnotations(true)
	// Deactivate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Move activation conditions to the top
	tree.SetMoveActivationConditionsToFront(true)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualDefaultExampleActivationConditionsFirst.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests escaping of symbols (e.g., quotation marks) and internal parentheses for visual output.
*/
func TestVisualOutputEscapingSymbols(t *testing.T) {

	// Statement with quotation marks and internal parentheses
	text := "The E(corporation) M(shall) F(be) P(a \"Type B\" corporation) Cex(pursuant to Section 201(b) of the New York State Not-for-Profit Corporation Law.)"

	// Activate annotations
	SetIncludeAnnotations(true)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualEscapingSymbols.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests linear multi-level nesting in visual output (i.e., Cac{Cac{Cac{}}}).
*/
func TestVisualOutputLinearMultiLevelNesting(t *testing.T) {

	// Statement with quotation marks and internal parentheses
	text := "A,p(First) A(Actor) I(action1) I(action2) Bdir{A(actor2) I(actionLevel2) Cac{A(actor3) I(actionLevel3) Bdir(some object)}}"

	// Deactivate annotations
	SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualLinearMultilevelNesting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests complex multi-level nesting in visual output (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestVisualOutputComplexMultiLevelNesting(t *testing.T) {

	// Statement with multi-level nesting withe embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) {   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Deactivate annotations
	SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)

	// Parse statement
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), IncludeAnnotations(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}
	fmt.Println("Generated output: " + output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputVisualComplexMultilevelNesting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests complex multi-level nesting in visual output (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNesting(t *testing.T) {

	// Statement with multi-level nesting withe embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) {   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)

	// Take separator for Google Sheets output
	separator := ";"

	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputStaticSchemaComplexMultilevelNesting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, nil, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, separator, "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := WriteToFile("errorOutput.error", output)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

// test with invalid statement and empty input nodes, unbalanced parentheses, missing ID
