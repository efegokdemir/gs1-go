package gs1

import "fmt"

// ValidateAssociations checks the Application Identifier pairing and
// exclusion rules supported by this package.
func (b Barcode) ValidateAssociations() error {
	seen := make(map[string]bool, len(b.Elements))
	for _, element := range b.Elements {
		seen[element.AI] = true
	}
	if err := validatePairRules(seen); err != nil {
		return err
	}
	if err := validateExclusiveRules(seen); err != nil {
		return err
	}
	return validateWeightRules(seen)
}

func validatePairRules(seen map[string]bool) error {
	if seen["01"] && seen["02"] {
		return associationError("AI (01) excludes AI (02)")
	}
	if seen["02"] && !seen["37"] {
		return associationError("AI (02) requires AI (37)")
	}
	if seen["37"] {
		if seen["01"] {
			return associationError("AI (01) excludes AI (37)")
		}
		if !seen["00"] || (!seen["02"] && !seen["8026"]) {
			return associationError("AI (37) requires AI (00) with AI (02) or AI (8026)")
		}
	}
	if seen["21"] && !seen["01"] && !seen["03"] && !seen["8006"] {
		return associationError("AI (21) requires AI (01), AI (03), or AI (8006)")
	}
	return nil
}

func validateExclusiveRules(seen map[string]bool) error {
	if seen["8017"] && seen["8018"] {
		return associationError("AI (8017) excludes AI (8018)")
	}
	return nil
}

func validateWeightRules(seen map[string]bool) error {
	prefixes := make(map[string]struct{})
	for ai := range seen {
		if len(ai) == 4 && ai[0] == '3' {
			prefixes[ai[:3]] = struct{}{}
		}
	}
	for prefix := range prefixes {
		variants := 0
		for ai := range seen {
			if len(ai) == 4 && ai[:3] == prefix {
				variants++
			}
		}
		if variants > 1 {
			return associationError("AI (%s) allows only one decimal variant", prefix)
		}
		if variants > 0 {
			requires := []string{"01", "02"}
			if prefix == "330" || prefix == "340" {
				requires = []string{"00", "01"}
			}
			if !seen[requires[0]] && !seen[requires[1]] {
				return associationError("AI (%s) requires AI (%s) or AI (%s)", prefix, requires[0], requires[1])
			}
		}
	}
	return nil
}

func associationError(message string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidAssociation, fmt.Sprintf(message, args...))
}
