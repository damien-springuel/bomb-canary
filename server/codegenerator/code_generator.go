package codegenerator

import "sync"

const maxCodeLength = 6

type randomCodeGenerator struct {
	nextRandomRune func() rune
	generatedCodes map[string]struct{}
	mut            *sync.Mutex
}

func New(nextRandomRune func() rune) randomCodeGenerator {
	return randomCodeGenerator{
		nextRandomRune: nextRandomRune,
		generatedCodes: make(map[string]struct{}),
		mut:            &sync.Mutex{},
	}
}

func (r randomCodeGenerator) checkAndSetCodeSuccessfully(code string) bool {
	r.mut.Lock()
	defer r.mut.Unlock()

	_, exists := r.generatedCodes[code]
	if exists {
		return false
	} else {
		r.generatedCodes[code] = struct{}{}
		return true
	}
}

func (r randomCodeGenerator) GenerateCode() string {
	code := ""
	for {
		code += string(r.nextRandomRune())
		if len(code) == maxCodeLength {
			if r.checkAndSetCodeSuccessfully(code) {
				break
			} else {
				code = ""
			}
		}
	}
	return code
}
