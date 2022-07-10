package argh

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkScannerPOSIXyScannerScan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		scanner := NewScanner(strings.NewReader(strings.Join([]string{
			"walrus",
			"-what",
			"--ball=awesome",
			"--elapsed",
			"carrot cake",
		}, string(nul))), nil)
		for {
			tok, _, _ := scanner.Scan()
			if tok == EOL {
				break
			}
		}
	}
}

func TestScannerPOSIXyScanner(t *testing.T) {
	for _, tc := range []struct {
		name              string
		argv              []string
		expectedTokens    []Token
		expectedLiterals  []string
		expectedPositions []Pos
	}{
		{
			name: "simple",
			argv: []string{"walrus", "-cake", "--corn-dog", "awkward"},
			expectedTokens: []Token{
				IDENT,
				ARG_DELIMITER,
				COMPOUND_SHORT_FLAG,
				ARG_DELIMITER,
				LONG_FLAG,
				ARG_DELIMITER,
				IDENT,
				EOL,
			},
			expectedLiterals: []string{
				"walrus", string(nul), "-cake", string(nul), "--corn-dog", string(nul), "awkward", "",
			},
			expectedPositions: []Pos{
				6, 7, 12, 13, 23, 24, 31, 32,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			scanner := NewScanner(strings.NewReader(strings.Join(tc.argv, string(nul))), nil)

			actualTokens := []Token{}
			actualLiterals := []string{}
			actualPositions := []Pos{}

			for {
				tok, lit, pos := scanner.Scan()

				actualTokens = append(actualTokens, tok)
				actualLiterals = append(actualLiterals, lit)
				actualPositions = append(actualPositions, pos)

				if tok == EOL {
					break
				}
			}

			r.Equal(tc.expectedTokens, actualTokens)
			r.Equal(tc.expectedLiterals, actualLiterals)
			r.Equal(tc.expectedPositions, actualPositions)
		})
	}
}
