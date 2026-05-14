/** A single AI-value pair from a GS1 barcode. */
export interface GS1Element {
  ai: string;
  value: string;
}

/** Result of parsing a GS1 barcode string. */
export interface GS1ParseResult {
  raw: string;
  elements: GS1Element[];
  gtin: string;
  lot: string;
  serial: string;
}

/** Error object returned when parsing fails. */
export interface GS1ParseError {
  error: string;
}

declare global {
  namespace gs1 {
    /**
     * Parse a GS1 barcode string (GS1-128, DataMatrix, bracket notation).
     * Returns a parsed result object, or an object with an `error` field.
     */
    function parse(input: string): GS1ParseResult | GS1ParseError;

    /**
     * Validate a GTIN check digit.
     * Accepts GTIN-8, GTIN-12, GTIN-13, and GTIN-14.
     */
    function validateGTIN(gtin: string): boolean;

    /**
     * Validate barcode against a LATAM regulator's requirements.
     * @param input - Raw barcode string
     * @param regulator - One of: "anvisa", "anmat", "snfa", "cofepris"
     * @returns null if compliant, or an error string
     */
    function validateRegulatory(
      input: string,
      regulator: "anvisa" | "anmat" | "snfa" | "cofepris"
    ): string | null;
  }
}
