package basic

import (
	"testing"
)

func Test_Bcrypt(t *testing.T) {
	testParserGood(t, "bcrypt", nil, RejectBcrypt, "$2y$05$NqR6p/K60C40W08.LDCTNeLpSC.gVVkadaLXysS5Y.nGNPltVacSi", "bar")
	testParserBad(t, "bcrypt", nil, RejectBcrypt, "$2y$0")
	testParserNot(t, "bcrypt", nil, RejectBcrypt, "plaintext")
}
