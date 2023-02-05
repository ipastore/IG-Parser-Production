package endpoints

import (
	"IG-Parser/core/exporter"
	"IG-Parser/core/parser"
	"IG-Parser/core/tree"
)

/*
This file contains the application endpoints that integrate the core parsing features, as well as file/output
handling. Both can be invoked with IG Script-encoded institutional statements to produce either tabular or
visual output for downstream processing, serving as endpoints for the use by specific applications, such as
web applications, console tools, etc.
*/

/*
Consumes statement as input and produces outfile.
Arguments include the IGScript-annotated statement, statement ID based on which substatements are generated,
the nature of the output type (see TabularOutputGeneratorConfig #OUTPUT_TYPE_CSV, #OUTPUT_TYPE_GOOGLE_SHEETS)
and a filename for the output. If the filename is empty, no output will be written.
If printHeaders is set, the output includes the header row.
Returns tabular output as string, and error (defaults to tree.PARSING_NO_ERROR).
*/
func ConvertIGScriptToTabularOutput(statement string, stmtId string, outputType string, filename string, printHeaders bool) (string, tree.ParsingError) {

	// Use separator specified by default
	separator := exporter.CellSeparator

	Println(" Step: Parse input statement")
	// Explicitly activate printing of shared elements
	//exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Parse IGScript statement into tree
	s, err := parser.ParseStatement(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	Println("Parsed statement:", s.String())

	// Run composite generation and return output and error. Will write file if filename != ""
	output, statementMap, statementHeader, statementHeaderNames, err :=
		exporter.GenerateTabularOutputFromParsedStatement(s, "", stmtId, filename, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, outputType, printHeaders)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	Println("  - Results:")
	Println("  - Header Symbols: ", statementHeader)
	Println("  - Header Names: ", statementHeaderNames)
	Println("  - Data: ", statementMap)

	Println("  - Output generation complete.")

	return output, err

}

/*
Consumes statement as input and produces outfile reflecting visual tree structure consumable by D3.
Arguments include the IGScript-annotated statement, statement ID (currently not used in visualization),
and a filename for the output. If the filename is empty, no output will be written.
Returns Visual tree structure as string, and error (defaults to tree.PARSING_NO_ERROR).
*/
func ConvertIGScriptToVisualTree(statement string, stmtId string, filename string) (string, tree.ParsingError) {

	Println(" Step: Parse input statement")
	// Explicitly activate printing of shared elements
	//exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Parse IGScript statement into tree
	s, err := parser.ParseStatement(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	// Prepare visual output
	Println(" Step: Generate visual output structure")
	output, err := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), exporter.IncludeAnnotations(),
		exporter.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	Println("  - Generated visual tree:", output)

	Println("  - Output generation complete.")

	if filename != "" {
		Println("  - Writing to file ...")

		err2 := exporter.WriteToFile(filename, output.String())
		if err2 != nil {
			Println("  - Problems when writing file "+filename+", Error:", err2)
		}

		Println("  - Writing completed.")
	}

	return output.String(), err
}
