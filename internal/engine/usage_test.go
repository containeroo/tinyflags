package engine

import (
	"bytes"
	"testing"

	"github.com/containeroo/tinyflags/internal/core"
	"github.com/stretchr/testify/assert"
)

type dummyValue struct {
	def     string
	changed bool
}

func (d *dummyValue) Set(s string) error { d.changed = true; return nil }
func (d *dummyValue) Get() any           { return d.def }
func (d *dummyValue) Changed() bool      { return d.changed }
func (d *dummyValue) Default() string    { return d.def }
func (d *dummyValue) IsSlice()           {}

func TestBuildFlagDescriptionHideDefault(t *testing.T) {
	t.Parallel()

	t.Run("hidesDefault", func(t *testing.T) {
		t.Parallel()

		flag := &core.BaseFlag{
			Name:        "f",
			Usage:       "short desc",
			HideDefault: true,
			DisableEnv:  true,
			Value:       &dummyValue{def: "default"},
		}

		desc := buildFlagDescription(flag, false, "app")
		assert.Equal(t, "short desc", desc)
	})
}

func TestUsageNotesUseConfiguredWriter(t *testing.T) {
	t.Parallel()

	t.Run("staticUsageNote", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		var buf bytes.Buffer
		fs.RegisterFlag("name", &core.BaseFlag{
			Name:  "name",
			Usage: "name",
			Value: &dummyValue{},
		})
		fs.SetStaticUsageNote("static note")

		fs.PrintStaticDefaults(&buf, 2, 20, 80)

		assert.Contains(t, buf.String(), "static note")
	})

	t.Run("dynamicUsageNote", func(t *testing.T) {
		t.Parallel()

		fs := NewFlagSet("app", ContinueOnError)
		var buf bytes.Buffer
		fs.DynamicGroup("http").String("host", "", "host")
		fs.SetDynamicUsageNote("dynamic note")

		fs.PrintDynamicDefaults(&buf, 2, 20, 80)

		assert.Contains(t, buf.String(), "dynamic note")
	})
}
