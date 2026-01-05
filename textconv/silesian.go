package textconv

import (
	"slices"
	"strings"
)

// Orthography ...
type Orthography int

// Orthography values
const (
	Lysohorsky Orthography = iota
	Wieczorek
	Slabikorz
)

// OrthographyFromString ...
func OrthographyFromString(s string) (Orthography, bool) {
	switch s {
	case "lysohorsky":
		return Lysohorsky, true
	case "wieczorek":
		return Wieczorek, true
	case "slabikorz":
		return Slabikorz, true
	}
	return Lysohorsky, false
}

var (
	pairs1 = []string{
		"mjo", "mio",
		"Mjo", "Mio",
		"mja", "mia",
		"Mja", "Mia",
		"mje", "mie",
		"Mje", "Mie",
		"mjé", "miy",
		"Mjé", "Miy",
		"bja", "bia",
		"Bja", "Bia",
		"bje", "bie",
		"Bje", "Bie",
		"pje", "pie",
		"Pje", "Pie",
		"bjo", "bio",
		"Bjo", "Bio",
		"bjé", "biy",
		"Bjé", "Biy",
		"źé", "ziy",
		"Źé", "Ziy",
		"ňõ", "niõ",
		"Ňõ", "Niŏ",
		"ňu", "niu",
		"Ňu", "Niu",
		"wje", "wie",
		"Wje", "Wie",
		"wě", "wie",
		"Wě", "Wie",
		"wja", "wia",
		"Wja", "Wia",
		"wjy", "wiy",
		"Wjy", "Wiy",
		"wjé", "wiy",
		"Wjé", "Wiy",
		"ňéj", "nij",
		"Ňéj", "Nij",
		"kjéj", "kij",
		"Kjéj", "Kij",
		"kjé", "ki",
		"Kjé", "Ki",
		"léj", "lij",
		"Léj", "Lij",
		"jéj", "jij",
		"Jéj", "Jij",
		"śe", "sie",
		"Śe", "Sie",
		"źe", "zie",
		"Źe", "Zie",
		"źé", "ziy",
		"Źé", "Ziy",
		"će", "cie",
		"Će", "Cie",
		"ćé", "ciy",
		"Ćé", "Ciy",
		"ća", "cia",
		"Ća", "Cia",
		"dźa", "dzia",
		"Dźa", "Dzia",
		"śo", "sio",
		"Śo", "Sio",
		"źo", "zio",
		"Źo", "Zio",
		"źi", "zi",
		"Źi", "Zi",
		"śi", "si",
		"Śi", "Si",
		"śa", "sia",
		"Śa", "Sia",
		"ći", "ci",
		"Ći", "Ci",
		"ňi", "ni",
		"Ňi", "Ni",
		"kě", "kie",
		"Kě", "Kie",
		"gě", "gie",
		"Gě", "Gie",
		"mě", "mie",
		"Mě", "Mie",
		"bě", "bie",
		"Bě", "Bie",
		"pě", "pie",
		"Pě", "Pie",
		"ně", "nie",
		"Ně", "Nie",
		"ňa", "nia",
		"Ňa", "Nia",
		"ňé", "niy",
		"Ňé", "Niy",
		"é", "y",
		"É", "Y",
		"š", "sz",
		"Š", "Sz",
		"č", "cz",
		"Č", "Cz",
		"ž", "ż",
		"Ž", "Ż",
		"ř", "rz",
		"Ř", "Rz",
		"ň", "ń",
		"Ň", "Ń",
		"ci", "cy",
		"Ci", "Cy",
		"si", "sy",
		"Si", "Sy",
	}
	// pairs2 = []string{
	// }
	replacerWieczorek = strings.NewReplacer(slices.Concat([]string{
		"źô", "ziô",
		"Źô", "Ziô",
		"ćô", "ciô",
		"Ćô", "ciô",
		"wjó", "wió",
		"Wjó", "Wió",
		"mjó", "mió",
		"Mjó", "Mió",
		"wjô", "wiô",
		"Wjô", "Wiô",
		"śó", "sió",
		"Śó", "sió",
		"ňó", "nió",
		"Ňó", "nió",
		"ćó", "ció",
		"ćó", "ció",
		"ňô", "niô",
		"Ňô", "niô",
	}, pairs1)...)
	replacerSlabikorz = strings.NewReplacer(slices.Concat([]string{
		"źô", "ziŏ",
		"Źô", "Ziŏ",
		"ćô", "ciŏ",
		"Ćô", "ciŏ",
		"mjó", "miō",
		"Mjó", "Miō",
		"wjó", "wiō",
		"Wjó", "Wiō",
		"wjô", "wiŏ",
		"Wjô", "Wiŏ",
		"śó", "siō",
		"Śó", "siō",
		"ňó", "niō",
		"Ňó", "niō",
		"ňô", "niŏ",
		"Ňô", "niŏ",
		"ćó", "ciō",
		"ćó", "ciō",
	},
		pairs1,
		[]string{
			"ò", "ô",
			"Ò", "Ô",
			"ô", "ŏ",
			"Ô", "Ŏ",
			"ó", "ō",
			"Ó", "Ō",
		})...)
	replacerSlabikorzOl = strings.NewReplacer(slices.Concat([]string{
		"ŏł", "oł",
	})...)
)

// ConvertEnclosed ...
func ConvertEnclosed(text string, ortho Orthography) string {
	var sb strings.Builder
	var block strings.Builder
	var inBlock bool
	for _, r := range text {
		if inBlock {
			if r == '}' {
				inBlock = false
				sb.WriteString(Convert(block.String(), ortho))
			} else {
				block.WriteRune(r)
			}
		} else {
			if r == '{' {
				inBlock = true
				block.Reset()
			} else {
				sb.WriteRune(r)
			}
		}
	}
	return sb.String()
}

// Convert ...
func Convert(text string, ortho Orthography) string {
	switch ortho {
	case Wieczorek:
		return replacerWieczorek.Replace(text)
	case Slabikorz:
		text = replacerSlabikorz.Replace(text)
		return replacerSlabikorzOl.Replace(text)
	default:
		return text
	}
}
