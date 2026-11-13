package ddl

import "github.com/lib/pq"

// q quotes an identifier. Always quoting keeps reserved words ("order",
// "user"), mixed case and names like "uuid-ossp" valid.
var q = pq.QuoteIdentifier
