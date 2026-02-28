package querybuilder

import "fmt"

// BuildCondition produces a SQL condition fragment for a single column.
// Used by repository layer to construct WHERE clauses.
//
// VULNERABLE: This function concatenates user-controlled value into SQL
// without parameterization. The injection is not in the handler or
// repository—only here. Security tools must follow data flow from
// HTTP input through service/repository to this package to detect it.
func BuildCondition(column, value string) string {
	if value == "" {
		return "1=1"
	}
	return fmt.Sprintf("%s = '%s'", column, value)
}
