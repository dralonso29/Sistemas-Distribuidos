package realstate

import (
	"errors"
	"sync"
)

//go run realstate.go
//go test
//go fmt realstate.go

type House struct {
	id           int
	calle        string
	numero       int
	piso         int
	letra        string
	codigo       int
	propietarios []DNI
}

type Owner struct {
	dni      DNI
	nombre   string
	apellido string
	casas    []int
}

type DNI struct {
	numero int
	letra  string
}

type RealState struct {
	listacasas        map[int]House
	listapropietarios map[DNI]Owner
	muxO              *sync.Mutex
	muxH              *sync.Mutex
}

func NewRealState() *RealState {
	muxO := &sync.Mutex{}
	muxH := &sync.Mutex{}
	ownermap := make(map[DNI]Owner)
	housemap := make(map[int]House)
	return &RealState{housemap, ownermap, muxO, muxH}
}

func (r *RealState) AddNewPair(o *Owner, h *House) error {
	_, own := r.listapropietarios[o.dni]
	_, ho := r.listacasas[h.id]
	if own || ho {
		return errors.New("Error AddNewPair")
	}
	r.muxH.Lock()
	r.muxO.Lock()
	o.casas = append(o.casas, h.id)
	h.propietarios = append(h.propietarios, o.dni)
	r.listapropietarios[o.dni] = *o
	r.listacasas[h.id] = *h
	r.muxO.Lock()
	r.muxH.Unlock()

	return nil
}

func (r *RealState) AddOwner(o *Owner) error {
	_, own := r.listapropietarios[o.dni]
	if own {
		return errors.New("Error AddOwner")
	}

	r.muxH.Lock()
	r.listapropietarios[o.dni] = *o
	r.muxO.Unlock()

	return nil
}

func (r *RealState) AddHouse(dni DNI, h *House) error {
	cas, ho := r.listacasas[h.id]
	prop, own := r.listapropietarios[dni]

	if !own || ho {
		return errors.New("Error AddHouse")
	}

	r.muxH.Lock()
	prop.casas = append(prop.casas, h.id)
	cas.propietarios = append(cas.propietarios, dni)
	r.muxH.Unlock()
	return nil
}

func (r *RealState) GetHouses(dni DNI) (listId []int, err error) {
	_, own := r.listapropietarios[dni]

	if !own {
		return listId, errors.New("Error GetHouses")
	}
	r.muxH.Lock()
	propietario := r.listapropietarios[dni]
	listId = propietario.casas
	r.muxH.Unlock()
	return listId, nil
}

func (r *RealState) GetOwners(id int) (listDni []DNI, err error) {
	_, ho := r.listacasas[id]

	if !ho {
		return listDni, errors.New("Error GetOwners")
	}
	r.muxO.Lock()
	casa := r.listacasas[id]
	listDni = casa.propietarios
	r.muxO.Unlock()
	return listDni, nil
}

func (r *RealState) GetHouse(id int) (h House, err error) {
	_, ho := r.listacasas[id]

	if !ho {
		return h, errors.New("Error GetHouse")
	}
	r.muxH.Lock()
	casa := r.listacasas[id]
	r.muxH.Unlock()
	return casa, nil
}

func (r *RealState) GetOwner(dni DNI) (o Owner, err error) {
	_, own := r.listapropietarios[dni]

	if !own {
		return o, errors.New("Error GetOwner")
	}
	r.muxO.Lock()
	propietario := r.listapropietarios[dni]
	r.muxO.Unlock()
	return propietario, nil
}

func (r *RealState) ChangeOwner(o *Owner) error {
	prop, own := r.listapropietarios[o.dni]

	if !own {
		return errors.New("Error ChangeOwner")
	}
	r.muxO.Lock()
	prop.nombre = "Carlos"
	r.listapropietarios[o.dni] = prop

	r.muxO.Unlock()
	return nil
}

func (r *RealState) ChangeHouse(h *House) error {
	cas, ho := r.listacasas[h.id]

	if !ho {
		return errors.New("Error ChangeHouse")
	}
	r.muxH.Lock()
	cas.calle = "Avenida"
	r.listacasas[h.id] = cas
	r.muxH.Unlock()
	return nil

}

func (r *RealState) DelOwner(dni DNI) error {
	prop, own := r.listapropietarios[dni]

	if !own {
		return errors.New("Error DelOwner")
	}
	r.muxO.Lock()
	for _, l := range prop.casas {
		cas := r.listacasas[l]
		for j, otracas := range cas.propietarios {
			if otracas.numero == dni.numero && otracas.letra == dni.letra {
				cas.propietarios = append(cas.propietarios[:j], cas.propietarios[j+1:]...)

			}

		}
		r.listacasas[l] = cas
	}
	r.muxO.Unlock()
	delete(r.listapropietarios, dni)
	return nil
}

func (r *RealState) DelHouse(id int) error {
	cas, ho := r.listacasas[id]

	if !ho {
		return errors.New("Error DelHouse")
	}

	r.muxH.Lock()
	for _, l := range cas.propietarios {
		prop := r.listapropietarios[l]
		r.listapropietarios[l] = prop
	}
	r.muxH.Unlock()
	delete(r.listacasas, id)
	return nil
}

