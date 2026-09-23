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
	if seen["37"] && !seen["02"] {
		return associationError("AI (37) requires AI (02)")
	}
	if seen["21"] && !seen["01"] {
		return associationError("AI (21) requires AI (01)")
	}
	return nil
}

func validateExclusiveRules(seen map[string]bool) error {
	if seen["8003"] && seen["8004"] {
		return associationError("AI (8003) excludes AI (8004)")
	}
	if seen["8017"] && seen["8018"] {
		return associationError("AI (8017) excludes AI (8018)")
	}
	if seen["253"] && len(seen) != 1 {
		return associationError("AI (253) must appear alone")
	}
	return nil
}

func validateWeightRules(seen map[string]bool) error {
	for _, prefix := range []string{"310", "320", "330", "340"} {
		variants := 0
		for ai := range seen {
			if len(ai) == 4 && ai[:3] == prefix {
				variants++
			}
		}
		if variants > 1 {
			return associationError("AI (%s) allows only one decimal variant", prefix)
		}
		if variants > 0 && !seen["01"] && !seen["02"] {
			return associationError("AI (%s) requires AI (01) or AI (02)", prefix)
		}
	}
	return nil
}

func associationError(message string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidAssociation, fmt.Sprintf(message, args...))
}
