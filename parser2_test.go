package argh_test

import (
	"testing"

	"git.meatballhat.com/x/box-o-sand/argh"
	"github.com/davecgh/go-spew/spew"
)

func TestParser2(t *testing.T) {
	for _, tc := range []struct {
		name     string
		args     []string
		commands []string
	}{
		{
			name: "basic",
			args: []string{
				"pies", "-eat", "--wat", "hello",
			},
			commands: []string{
				"hello",
			},
		},
	} {
		t.Run(tc.name, func(ct *testing.T) {
			pt, err := argh.ParseArgs2(tc.args, tc.commands)
			if err != nil {
				ct.Logf("err=%+#v", err)
				return
			}

			spew.Dump(pt)
		})
	}
}
