package basic

import (
	"testing"
)

func Test_Plain(t *testing.T) {
	testParserGood(t, "plain", AcceptPlain, RejectPlain, "bar", "bar")
	//testParserBad() plain takes anything
	// testParserNot()  plain takes anything
}
