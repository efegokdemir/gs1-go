//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/galenzo17/health-interop/gs1"
)

type elementJSON struct {
	AI    string `json:"ai"`
	Value string `json:"value"`
}

type parseResultJSON struct {
	Raw      string        `json:"raw"`
	Elements []elementJSON `json:"elements"`
	GTIN     string        `json:"gtin"`
	Lot      string        `json:"lot"`
	Serial   string        `json:"serial"`
}

func parse(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return errorResult("parse requires 1 argument")
	}
	input := args[0].String()

	b, err := gs1.Parse(input)
	if err != nil {
		return errorResult(err.Error())
	}

	result := parseResultJSON{
		Raw:      b.Raw,
		Elements: make([]elementJSON, len(b.Elements)),
		GTIN:     b.GTIN(),
		Lot:      b.Lot(),
		Serial:   b.SerialNumber(),
	}
	for i, e := range b.Elements {
		result.Elements[i] = elementJSON{AI: e.AI, Value: e.Value}
	}

	data, err := json.Marshal(result)
	if err != nil {
		return errorResult(err.Error())
	}
	return js.Global().Get("JSON").Call("parse", string(data))
}

func validateGTIN(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return false
	}
	return gs1.ValidateGTIN(args[0].String()) == nil
}

func validateRegulatory(_ js.Value, args []js.Value) any {
	if len(args) < 2 {
		return "validateRegulatory requires 2 arguments: input, regulator"
	}
	input := args[0].String()
	regName := args[1].String()

	b, err := gs1.Parse(input)
	if err != nil {
		return err.Error()
	}

	regulators := map[string]gs1.Regulator{
		"anvisa":   gs1.ANVISA,
		"anmat":    gs1.ANMAT,
		"snfa":     gs1.SNFA,
		"cofepris": gs1.COFEPRIS,
	}

	reg, ok := regulators[regName]
	if !ok {
		return "unknown regulator: " + regName
	}

	if err := reg.Validate(b); err != nil {
		return err.Error()
	}
	return js.Null()
}

func errorResult(msg string) any {
	obj := js.Global().Get("Object").New()
	obj.Set("error", msg)
	return obj
}

func main() {
	ns := js.Global().Get("Object").New()
	ns.Set("parse", js.FuncOf(parse))
	ns.Set("validateGTIN", js.FuncOf(validateGTIN))
	ns.Set("validateRegulatory", js.FuncOf(validateRegulatory))
	js.Global().Set("gs1", ns)

	// Keep the Go runtime alive.
	<-make(chan struct{})
}
