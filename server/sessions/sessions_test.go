package sessions

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"
)

type testUUID struct {
	uuidToReturn string
}

func (t testUUID) Create() string {
	return t.uuidToReturn
}

func Test_Create(t *testing.T) {
	s := New(testUUID{uuidToReturn: "myUuid"})
	s.Create("name")

	g := NewWithT(t)
	g.Expect(s.nameBySessionId).To(Equal(map[string]string{
		"myUuid": "name",
	}))

	name, err := s.Get("myUuid")
	g.Expect(name).To(Equal("name"))
	g.Expect(err).To(BeNil())
}

func Test_Get_DoesntExist(t *testing.T) {
	s := New(testUUID{uuidToReturn: "myUuid"})
	s.Create("name")

	name, err := s.Get("randomUuid")

	g := NewWithT(t)
	g.Expect(name).To(Equal(""))
	g.Expect(err).To(Equal(errors.New("session doesn't exist")))
}
