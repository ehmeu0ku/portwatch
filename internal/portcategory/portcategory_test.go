package portcategory_test

import (
	"testing"

	"github.com/yourusername/portwatch/internal/portcategory"
)

func TestClassifyWebPort(t *testing.T) {
	cl := portcategory.New(nil)
	if got := cl.Classify(80); got != portcategory.Web {
		t.Fatalf("expected web, got %s", got)
	}
}

func TestClassifyHTTPS(t *testing.T) {
	cl := portcategory.New(nil)
	if got := cl.Classify(443); got != portcategory.Web {
		t.Fatalf("expected web, got %s", got)
	}
}

func TestClassifyDatabasePort(t *testing.T) {
	cl := portcategory.New(nil)
	for _, p := range []uint16{3306, 5432, 6379, 27017} {
		if got := cl.Classify(p); got != portcategory.Database {
			t.Fatalf("port %d: expected database, got %s", p, got)
		}
	}
}

func TestClassifyMessagingPort(t *testing.T) {
	cl := portcategory.New(nil)
	if got := cl.Classify(9092); got != portcategory.Messaging {
		t.Fatalf("expected messaging, got %s", got)
	}
}

func TestClassifyUnknownPort(t *testing.T) {
	cl := portcategory.New(nil)
	if got := cl.Classify(9999); got != portcategory.Unknown {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestCustomMappingOverridesBuiltin(t *testing.T) {
	cl := portcategory.New(map[uint16]portcategory.Category{80: portcategory.Database})
	if got := cl.Classify(80); got != portcategory.Database {
		t.Fatalf("expected database override, got %s", got)
	}
}

func TestCustomMappingAddsNewPort(t *testing.T) {
	cl := portcategory.New(map[uint16]portcategory.Category{8888: portcategory.Web})
	if got := cl.Classify(8888); got != portcategory.Web {
		t.Fatalf("expected web, got %s", got)
	}
}

func TestBuiltinUnaffectedByOtherCustom(t *testing.T) {
	cl := portcategory.New(map[uint16]portcategory.Category{8080: portcategory.Web})
	if got := cl.Classify(22); got != portcategory.Remote {
		t.Fatalf("expected remote, got %s", got)
	}
}
