package gs1

// dataType distinguishes numeric-only from alphanumeric AI data fields.
type dataType int

const (
	dataNumeric      dataType = iota // N: digits 0-9 only
	dataAlphanumeric                 // X: any character
)

// aiSpec describes the format of a single GS1 Application Identifier.
type aiSpec struct {
	AI       string   // AI code, e.g. "01", "10", "3100"
	Name     string   // human-readable name
	FixedLen int      // data length for fixed-length AIs; 0 means variable
	MaxLen   int      // max data length (equals FixedLen for fixed-length AIs)
	DataType dataType // numeric or alphanumeric
}

// aiTable maps AI codes to their specifications. The parser looks up 2-digit
// keys first, then 3-digit, then 4-digit.
var aiTable = map[string]aiSpec{
	// SSCC
	"00": {AI: "00", Name: "SSCC", FixedLen: 18, MaxLen: 18, DataType: dataNumeric},

	// GTIN
	"01": {AI: "01", Name: "GTIN", FixedLen: 14, MaxLen: 14, DataType: dataNumeric},
	"02": {AI: "02", Name: "Content GTIN", FixedLen: 14, MaxLen: 14, DataType: dataNumeric},

	// Batch / Lot
	"10": {AI: "10", Name: "Batch/Lot", FixedLen: 0, MaxLen: 20, DataType: dataAlphanumeric},

	// Dates (YYMMDD)
	"11": {AI: "11", Name: "Production Date", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"13": {AI: "13", Name: "Packaging Date", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"15": {AI: "15", Name: "Best Before Date", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"17": {AI: "17", Name: "Expiration Date", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},

	// Serial number
	"21": {AI: "21", Name: "Serial Number", FixedLen: 0, MaxLen: 20, DataType: dataAlphanumeric},

	// Counts
	"30": {AI: "30", Name: "Count", FixedLen: 0, MaxLen: 8, DataType: dataNumeric},
	"37": {AI: "37", Name: "Count of Trade Items", FixedLen: 0, MaxLen: 8, DataType: dataNumeric},

	// Additional product identification
	"240": {AI: "240", Name: "Additional Product ID", FixedLen: 0, MaxLen: 30, DataType: dataAlphanumeric},
	"241": {AI: "241", Name: "Customer Part Number", FixedLen: 0, MaxLen: 30, DataType: dataAlphanumeric},

	// Net weight, kg (310n — 4-digit AIs where n is decimal position)
	"3100": {AI: "3100", Name: "Net Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3101": {AI: "3101", Name: "Net Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3102": {AI: "3102", Name: "Net Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3103": {AI: "3103", Name: "Net Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3104": {AI: "3104", Name: "Net Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3105": {AI: "3105", Name: "Net Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},

	// Net weight, lb (320n)
	"3200": {AI: "3200", Name: "Net Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3201": {AI: "3201", Name: "Net Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3202": {AI: "3202", Name: "Net Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3203": {AI: "3203", Name: "Net Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3204": {AI: "3204", Name: "Net Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3205": {AI: "3205", Name: "Net Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},

	// Gross weight, kg (330n)
	"3300": {AI: "3300", Name: "Gross Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3301": {AI: "3301", Name: "Gross Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3302": {AI: "3302", Name: "Gross Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3303": {AI: "3303", Name: "Gross Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3304": {AI: "3304", Name: "Gross Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3305": {AI: "3305", Name: "Gross Weight kg", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},

	// Gross weight, lb (340n)
	"3400": {AI: "3400", Name: "Gross Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3401": {AI: "3401", Name: "Gross Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3402": {AI: "3402", Name: "Gross Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3403": {AI: "3403", Name: "Gross Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3404": {AI: "3404", Name: "Gross Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},
	"3405": {AI: "3405", Name: "Gross Weight lb", FixedLen: 6, MaxLen: 6, DataType: dataNumeric},

	// GLN — Global Location Number
	"414": {AI: "414", Name: "GLN", FixedLen: 13, MaxLen: 13, DataType: dataNumeric},

	// GSIN — Global Shipment Identification Number
	"402": {AI: "402", Name: "GSIN", FixedLen: 17, MaxLen: 17, DataType: dataNumeric},

	// Custom / internal use (90-99)
	"90": {AI: "90", Name: "Internal", FixedLen: 0, MaxLen: 30, DataType: dataAlphanumeric},
	"91": {AI: "91", Name: "Internal", FixedLen: 0, MaxLen: 90, DataType: dataAlphanumeric},
	"92": {AI: "92", Name: "Internal", FixedLen: 0, MaxLen: 90, DataType: dataAlphanumeric},
	"93": {AI: "93", Name: "Internal", FixedLen: 0, MaxLen: 90, DataType: dataAlphanumeric},
	"94": {AI: "94", Name: "Internal", FixedLen: 0, MaxLen: 90, DataType: dataAlphanumeric},
	"95": {AI: "95", Name: "Internal", FixedLen: 0, MaxLen: 90, DataType: dataAlphanumeric},
	"96": {AI: "96", Name: "Internal", FixedLen: 0, MaxLen: 90, DataType: dataAlphanumeric},
	"97": {AI: "97", Name: "Internal", FixedLen: 0, MaxLen: 90, DataType: dataAlphanumeric},
	"98": {AI: "98", Name: "Internal", FixedLen: 0, MaxLen: 90, DataType: dataAlphanumeric},
	"99": {AI: "99", Name: "Internal", FixedLen: 0, MaxLen: 90, DataType: dataAlphanumeric},

	// National Healthcare Reimbursement Number (NHRN)
	"710": {AI: "710", Name: "NHRN Germany", FixedLen: 0, MaxLen: 20, DataType: dataAlphanumeric},
	"711": {AI: "711", Name: "NHRN France", FixedLen: 0, MaxLen: 20, DataType: dataAlphanumeric},
	"712": {AI: "712", Name: "NHRN Spain", FixedLen: 0, MaxLen: 20, DataType: dataAlphanumeric},
	"713": {AI: "713", Name: "NHRN Brazil", FixedLen: 0, MaxLen: 20, DataType: dataAlphanumeric},
	"714": {AI: "714", Name: "NHRN Portugal", FixedLen: 0, MaxLen: 20, DataType: dataAlphanumeric},
}

// lookupAI finds the AI spec at the given position in the input string.
// It tries 2-digit, then 3-digit, then 4-digit keys. Returns the spec
// and the number of characters consumed for the AI code, or false if
// no match was found.
func lookupAI(input string, pos int) (aiSpec, int, bool) {
	remaining := len(input) - pos

	if remaining >= 2 {
		if spec, ok := aiTable[input[pos:pos+2]]; ok {
			return spec, 2, true
		}
	}
	if remaining >= 3 {
		if spec, ok := aiTable[input[pos:pos+3]]; ok {
			return spec, 3, true
		}
	}
	if remaining >= 4 {
		if spec, ok := aiTable[input[pos:pos+4]]; ok {
			return spec, 4, true
		}
	}

	return aiSpec{}, 0, false
}
