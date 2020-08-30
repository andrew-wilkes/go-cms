package page

import (
	"gocms/pkg/files"
	"testing"
)

func TestGetByRoute(t *testing.T) {
	files.Root = "../files/"
	GetByRoute("test", "test_page", false)
}
