package fetcher

import "strings"

func NormalizeString(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)

	replacer := strings.NewReplacer(
		"&", " and ",
		"+", " and ",
		"’", "'",
		"`", "'",
		"´", "'",
		"“", `"`,
		"”", `"`,
		"–", "-",
		"—", "-",
		"_", " ",
		".", " ",
		",", " ",
		":", " ",
		";", " ",
		"!", " ",
		"?", " ",
		"(", " ",
		")", " ",
		"[", " ",
		"]", " ",
		"{", " ",
		"}", " ",
	)
	s = replacer.Replace(s)

	// Collapse repeated whitespace.
	s = strings.Join(strings.Fields(s), " ")

	return s
}
