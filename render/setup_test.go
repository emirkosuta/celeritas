package render

import (
	"os"
	"testing"

	"github.com/CloudyKit/jet/v6"
)

var view = jet.NewSet(
	jet.NewOSFileSystemLoader("./testdata/views"),
	jet.InDevelopmentMode(),
)

var testRenderer = Render{
	Renderer: "",
	RootPath: "",
	JetViews: view,
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
