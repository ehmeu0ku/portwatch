// Package portcategory classifies ports into broad functional categories
// such as database, web, messaging, and so on.
package portcategory

// Category is a broad classification for a port.
type Category string

const (
	Web       Category = "web"
	Database  Category = "database"
	Messaging Category = "messaging"
	Remote    Category = "remote"
	DNS       Category = "dns"
	Mail      Category = "mail"
	Unknown   Category = "unknown"
)

// Classifier maps ports to categories.
type Classifier struct {
	custom map[uint16]Category
}

var builtins = map[uint16]Category{
	21:    Remote,
	22:    Remote,
	23:    Remote,
	25:    Mail,
	53:    DNS,
	80:    Web,
	110:   Mail,
	143:   Mail,
	443:   Web,
	3306:  Database,
	5432:  Database,
	6379:  Database,
	27017: Database,
	4222:  Messaging,
	5672:  Messaging,
	9092:  Messaging,
}

// New returns a Classifier seeded with built-in mappings.
func New(custom map[uint16]Category) *Classifier {
	c := &Classifier{custom: make(map[uint16]Category)}
	for k, v := range custom {
		c.custom[k] = v
	}
	return c
}

// Classify returns the Category for the given port number.
func (c *Classifier) Classify(port uint16) Category {
	if cat, ok := c.custom[port]; ok {
		return cat
	}
	if cat, ok := builtins[port]; ok {
		return cat
	}
	return Unknown
}
