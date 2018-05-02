package parsers

import (
	"unicode/utf8"
	"strings"
	"bytes"
)

// Parses the supplied file content against the list of licences and
// returns the matching licences with the shortname and the percentage of match.
// If fast lookup methods fail it will try deep matching which is slower.
func guessLicenseFast(content []byte, deepguess bool, licenses []License) []LicenseMatch {
	var matchingLicenses []LicenseMatch

	for _, license := range keywordGuessLicenseFast(content, licenses) {
		matchingLicense := License{}

		for _, l := range licenses {
			if l.LicenseId == license.LicenseId {
				matchingLicense = l
				break
			}
		}

		trimto := utf8.RuneCountInString(matchingLicense.LicenseText)

		if trimto > len(content) {
			trimto = len(content)
		}

		//contentConcordance := vectorspace.BuildConcordance(string(runecontent[:trimto]))
		//relation := vectorspace.Relation(matchingLicense.Concordance, contentConcordance)
		//
		//if relation >= confidence {
		//	matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: relation})
		//}
	}

	//if len(matchingLicenses) == 0 && deepguess == true {
	//	for _, license := range licenses {
	//		trimto := utf8.RuneCountInString(license.LicenseText)
	//
	//		if trimto > len(content) {
	//			trimto = len(content)
	//		}
	//
	//		contentConcordance := vectorspace.BuildConcordance(string(content[:trimto]))
	//		relation := vectorspace.Relation(license.Concordance, contentConcordance)
	//
	//		if relation >= confidence {
	//			matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: relation})
	//		}
	//	}
	//}

	//sort.Slice(matchingLicenses, func(i, j int) bool {
	//	return matchingLicenses[i].Percentage > matchingLicenses[j].Percentage
	//})

	//// Special cases such as MIT and JSON here
	//if len(matchingLicenses) > 2 && ((matchingLicenses[0].LicenseId == "JSON" && matchingLicenses[1].LicenseId == "MIT") ||
	//	(matchingLicenses[0].LicenseId == "MIT" && matchingLicenses[1].LicenseId == "JSON")) {
	//	if strings.Contains(strings.ToLower(content), "not evil") {
	//		// Its JSON
	//		matchingLicenses = []LicenseMatch{}
	//		matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: "JSON", Percentage: 1})
	//	} else {
	//		// Its MIT
	//		matchingLicenses = []LicenseMatch{}
	//		matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: "MIT", Percentage: 1})
	//	}
	//}

	return matchingLicenses
}

// Fast method of checking if supplied content contains a licence using
// matching keyword ngrams to find if the licence is a match or not
// returns the maching licences with shortname and the percentage of match.
func keywordGuessLicenseFast(content []byte, licenses []License) []LicenseMatch {
	content = cleanTextFast(content)

	matchingLicenses := []LicenseMatch{}

	for _, license := range licenses {
		keywordmatch := 0
		contains := false

		for _, keyword := range license.Keywords {
			contains = bytes.Contains(content, []byte(strings.ToLower(keyword)))

			if contains == true {
				keywordmatch++
			}
		}

		if keywordmatch > 0 {
			percentage := (float64(keywordmatch) / float64(len(license.Keywords))) * 100
			matchingLicenses = append(matchingLicenses, LicenseMatch{LicenseId: license.LicenseId, Percentage: percentage})
		}
	}

	return matchingLicenses
}

func cleanTextFast(content []byte) []byte {
	content = bytes.ToLower(content)

	tmp := alphaNumericRegex.ReplaceAllString(string(content), " ")
	tmp = multipleSpacesRegex.ReplaceAllString(tmp, " ")

	return []byte(tmp)
}