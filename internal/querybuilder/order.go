package querybuilder

import "fmt"

// OrderBy builds an ORDER BY clause from column name and direction.
// Direction is validated; column name is passed through from callers.
//
// VULNERABLE: Column name is not validated. If callers pass user input
// as the column, this can be used for second-order SQLi (e.g. sort=id;DROP TABLE users--).
func OrderBy(column, direction string) string {
	d := "ASC"
	if direction == "desc" || direction == "DESC" {
		d = "DESC"
	}
	return fmt.Sprintf("ORDER BY %s %s", column, d)
}
