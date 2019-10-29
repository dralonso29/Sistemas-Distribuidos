package realstate

import (
	"testing"
)

func Test1(t *testing.T) {
	var (
		r *RealState = NewRealState()
	)
	dni := DNI{111, "a"}
	dni2 := DNI{222, "b"}
	o := Owner{dni, "aaa", "a1", nil}
	o2 := Owner{dni2, "bbb", "b2", nil}
	o3 := Owner{dni, "ccc", "c3", nil}
	h := House{11, "calle", 2323, 23, "letra", 2323, nil}
	h2 := House{22, "calle2", 4444, 24, "letraB", 2323, nil}
	h3 := House{33, "calle3", 4444, 25, "letraC", 2323, nil}

	go func() {
		var er = r.AddNewPair(&o, &h)
		if er != nil {
			t.Error(er)
		}
	}()

	go func() {
		var er = r.AddNewPair(&o2, &h)
		if er != nil {
			t.Error(er)
		}
	}()

	go func() {
		var er = r.AddNewPair(&o3, &h2)
		if er != nil {
			t.Error(er)
		}
	}()

	go func() {
		var er = r.AddOwner(&o2)
		if er != nil {
			t.Error(er)
		}
	}()

	go func() {
		var er = r.AddHouse(dni2, &h2)
		if er != nil {
			t.Error(er)
		}
	}()

	go func() {
		var er = r.DelOwner(dni)
		if er != nil {
			t.Error(er)
		}
	}()

	go func() {
		var er = r.DelOwner(dni2)
		if er != nil {
			t.Error(er)
		}
	}()

	go func() {
		var er = r.DelHouse(22222)
		if er != nil {
			t.Error(er)
		}
	}()

	go func() {
		var er = r.ChangeHouse(&h3)
		if er != nil {
			t.Error(er)
		}
	}()

}
