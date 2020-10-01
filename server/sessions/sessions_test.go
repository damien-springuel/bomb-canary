package sessions

import (
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
	s.Create("code", "name")

	g := NewWithT(t)
	g.Expect(s.sessionById).To(Equal(map[string]playerSession{
		"myUuid": {partyCode: "code", name: "name"},
	}))
}
