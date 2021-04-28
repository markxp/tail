// +build windows

package powershell

import (
	"fmt"
	"strings"
	"testing"
)

func TestPSFlagsCommand(t *testing.T) {
	tt := []struct {
		name  string
		input PSFlags
	}{
		{"default: an empty profile settings, non-interactive, UTF-8 terminal", PSFlags{}},
		{"use profile", PSFlags{UseProfile: true}},
		{"interative terminal", PSFlags{Interactive: true}},
		{"use local encoding", PSFlags{LocalEncoding: true}},

		{"use local encoding + interactive terminal", PSFlags{UseProfile: true, Interactive: true}},

		{"user's default: use profile, interactive terminal, use local encoding", PSFlags{UseProfile: true, Interactive: true, LocalEncoding: true}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			cmd := tc.input.Command("echo", "hello")
			txt := strings.Join(cmd.Args[:len(cmd.Args)-2], " ")

			if (strings.Index(txt, flagNoProfile) < 0) != tc.input.UseProfile {
				var msg string
				if tc.input.UseProfile {
					msg = fmt.Sprintf("expect no %q.", flagNoProfile)
				} else {
					msg = fmt.Sprintf("expect %q.", flagNoProfile)
				}
				t.Fatalf("%s full command: %q", msg, cmd.String())
			}

			if (strings.Index(txt, flagNoInteractive) < 0) != tc.input.Interactive {
				var msg string
				if tc.input.Interactive {
					msg = fmt.Sprintf("expect no %q.", flagNoInteractive)
				} else {
					msg = fmt.Sprintf("expect %q.", flagNoInteractive)
				}
				t.Fatalf("%s full command: %q", msg, cmd.String())
			}

			if (strings.Index(txt, cmdUTF8Endcoding) < 0) != tc.input.LocalEncoding {
				var msg string
				if tc.input.LocalEncoding {
					msg = fmt.Sprintf("expect no %q.", "utf8_encoding_setting")
				} else {
					msg = fmt.Sprintf("expect %q.", "utf8_encoding_setting")
				}
				t.Fatalf("%s full command: %q", msg, cmd.String())
			}

		})

	}
}
