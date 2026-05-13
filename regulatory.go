package gs1

import "fmt"

// Regulator represents a pharmaceutical regulatory body that mandates
// specific GS1 Application Identifiers for drug traceability.
type Regulator struct {
	Name        string   // regulatory body name, e.g., "ANVISA"
	Country     string   // ISO 3166-1 alpha-2 country code
	RequiredAIs []string // AI codes that must be present in compliant barcodes
}

// LATAM pharmaceutical regulators and their required AIs.
var (
	// ANVISA — Brazil (RDC 157/2017, SNCM). Requires GTIN + expiry + lot + serial.
	ANVISA = Regulator{Name: "ANVISA", Country: "BR", RequiredAIs: []string{"01", "17", "10", "21"}}

	// ANMAT — Argentina (Disposición 3683/2011, SNT). Requires GTIN + expiry + lot + serial.
	ANMAT = Regulator{Name: "ANMAT", Country: "AR", RequiredAIs: []string{"01", "17", "10", "21"}}

	// SNFA — Chile (ISP/MINSAL). Requires GTIN + expiry + lot.
	SNFA = Regulator{Name: "SNFA", Country: "CL", RequiredAIs: []string{"01", "17", "10"}}

	// COFEPRIS — Mexico. Requires GTIN + expiry + lot.
	COFEPRIS = Regulator{Name: "COFEPRIS", Country: "MX", RequiredAIs: []string{"01", "17", "10"}}
)

// Validate checks that the barcode contains all Application Identifiers
// required by this regulator. Returns nil if compliant, or an error listing
// the first missing AI.
func (r Regulator) Validate(b Barcode) error {
	for _, ai := range r.RequiredAIs {
		if _, ok := b.Get(ai); !ok {
			return fmt.Errorf("%w: %s requires AI (%s)", ErrMissingRequiredAI, r.Name, ai)
		}
	}
	return nil
}

// ValidateANVISA checks compliance with Brazil's ANVISA traceability requirements.
func (b Barcode) ValidateANVISA() error { return ANVISA.Validate(b) }

// ValidateANMAT checks compliance with Argentina's ANMAT traceability requirements.
func (b Barcode) ValidateANMAT() error { return ANMAT.Validate(b) }

// ValidateSNFA checks compliance with Chile's SNFA traceability requirements.
func (b Barcode) ValidateSNFA() error { return SNFA.Validate(b) }

// ValidateCOFEPRIS checks compliance with Mexico's COFEPRIS traceability requirements.
func (b Barcode) ValidateCOFEPRIS() error { return COFEPRIS.Validate(b) }
