package files

import "testing"

func TestRoot(t *testing.T) {
	Root = "/test/"
}

func TestGetTemplate(t *testing.T) {
	Root = "../files/"
	GetTemplate("test", "home")
}
