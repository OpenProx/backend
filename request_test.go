package backend

import "testing"

func TestRequestToken(t *testing.T) {
	token, err := GenerateRequestToken(2, 10)
	if err != nil {
		t.Error(err)
		return
	}

	p, id, _, err := DecodeRequestToken(token)
	if err != nil {
		t.Error(err)
		return
	}

	if p != 2 || id != 10 {
		t.Errorf("Wrong values after decode 2/10 -> %d/%d", p, id)
	}
}
