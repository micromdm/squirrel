package basic

import (
	"testing"
)

func testParserGood(t *testing.T, name string, accept PasswdParser, reject PasswdParser, hashed string, passwd string) {
	if accept != nil {
		ep, err := accept(hashed)
		if err != nil {
			t.Errorf("%s accept (%s) failed to parse: %s", name, hashed, err.Error())
		} else if ep == nil {
			t.Errorf("%s accept (%s) failed to yield an EncodedPasswd", name, hashed)
		} else {
			if !ep.MatchesPassword(passwd) {
				t.Errorf("%s accept (%s) failed to match password (%s)", name, hashed, passwd)
			}
			if ep.MatchesPassword(passwd + "not") {
				t.Errorf("%s accept (%s) failed by matching password (%s)", name, hashed, passwd+"not")
			}
		}
	}

	if reject != nil {
		ep, err := reject(hashed)
		if ep != nil {
			t.Errorf("%s reject (%s) yielded an EncodedPasswd", name, hashed)
		} else if err == nil {
			t.Errorf("%s reject (%s) did not return an error", name, hashed)
		}
	}
}

func testParserBad(t *testing.T, name string, accept PasswdParser, reject PasswdParser, hashed string) {
	if accept != nil {
		ep, err := accept(hashed)
		if ep != nil {
			t.Errorf("%s accept (%s) yielded a EncodedPasswd", name, hashed)
		} else if err == nil {
			t.Errorf("%s accept (%s) did not return an error", name, hashed)
		}
	}
	if reject != nil {
		ep, err := reject(hashed)
		if ep != nil {
			t.Errorf("%s reject (%s) yielded a EncodedPasswd", name, hashed)
		} else if err == nil {
			t.Errorf("%s reject (%s) did not return an error", name, hashed)
		}
	}
}

func testParserNot(t *testing.T, name string, accept PasswdParser, reject PasswdParser, hashed string) {
	if accept != nil {
		ep, err := accept(hashed)
		if ep != nil {
			t.Errorf("%s accept (%s) yielded a EncodedPasswd", name, hashed)
		} else if err != nil {
			t.Errorf("%s accept (%s) errored instead of ignoring", name, hashed)
		}
	}
	if reject != nil {
		ep, err := reject(hashed)
		if ep != nil {
			t.Errorf("%s reject (%s) yielded a EncodedPasswd", name, hashed)
		} else if err != nil {
			t.Errorf("%s reject (%s) errored instead of ignoring", name, hashed)
		}
	}
}
