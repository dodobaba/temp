package myutil

import (
	"errors"
	"regexp"
	"shoppingzone/mylib/mylog"
)

//Validata : validata
type Validata struct{}

//ValiString :
func (v *Validata) ValiString(s string, rxp string, min int, max int) (bool, error) {
	if len(s) >= min && len(s) <= max {
		vali := regexp.MustCompile(rxp)
		if vali.MatchString(s) {
			return true, nil
		}
		mylog.Tf("[Error]", "Validata", "ValidataString", "Fail for field's style. '%s'", s)
		return false, errors.New("Fail for field's style")
	}
	mylog.Tf("[Error]", "Validata", "ValidataString", "Field's length is fail. '%s'", s)
	return false, errors.New("Field's length is fail")
}

//ValiStrings :
func (v *Validata) ValiStrings(arg ...ValiStringType) <-chan error {
	out := make(chan error, 10)
	go func() {
		for _, s := range arg {
			_, err := v.ValiString(s.S, s.Rxp, s.Min, s.Max)
			if err != nil {
				out <- err
				break
			}
			out <- nil
		}
		close(out)
	}()
	return out
}

//ValiStringType :
type ValiStringType struct {
	S   string
	Rxp string
	Min int
	Max int
}
