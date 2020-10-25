package codegenerator

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_GenerateCode(t *testing.T) {
	i := 0
	randomCharFunc := func() rune {
		defer func() {
			i += 1
		}()
		return rune("ABCDEF000000000"[i])
	}
	generator := New(randomCharFunc)
	code := generator.GenerateCode()

	g := NewWithT(t)
	g.Expect(code).To(Equal("ABCDEF"))
}

func Test_GenerateCode_CantGenerateTwiceTheSameCode(t *testing.T) {
	i := 0
	randomCharFunc := func() rune {
		defer func() {
			i += 1
		}()
		return rune("ABCDEFABCDEFABCDEE000000000"[i])
	}
	generator := New(randomCharFunc)
	_ = generator.GenerateCode()
	code := generator.GenerateCode()

	g := NewWithT(t)
	g.Expect(code).To(Equal("ABCDEE"))
}
