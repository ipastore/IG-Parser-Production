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
	output, outArr, _ := generateHeaderRow(output, componentIdx, ";")

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
	output, outArr, _ := generateHeaderRow(output, componentIdx, ";")

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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false
	// OVERRIDE dynamic output setting
	tree.AGGREGATE_IMPLICIT_LINKAGES = true

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
but including shared elements in output.
*/
func TestTabularOutputWithSharedElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := ioutil.ReadFile("TestOutputSimpleNoNestingWithSharedElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
		"{Cac{E(Program Manager) F(is) P(approved)} [AND] " +
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Implicitly linked nested statement
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// With shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Generated Component References: ", componentRefs)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
		"Bdir,p(approved) Bdir,p(certified) Bdir((production [operations] [AND] handling operations)) "+//Bdir,p1(accredited) Bdir1(certifying agents) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Static output
	SetDynamicOutput(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
		"Bdir,p(recognized) Bdir,p1(accredited) Bdir1(certifying agents) Bdir(other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Static output
	SetDynamicOutput(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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
		"Bdir,p(recognized) Bdir,p1(accredited) Bdir1(certifying agents) Bdir(other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir,p1(responsible) Bdir1(Program Manager) and Bdir,p2(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// No shared elements
	INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

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

	// Take separator for Google Sheets output
	separator := headerRowSeparator

	statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650", separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "")
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

// Test annotation mapping based on private nodes

// ensure ordering of column headers

// introduce for statement combinations

// introduce feature for other components

// test with invalid statement and empty input nodes, unbalanced parentheses, missing ID
