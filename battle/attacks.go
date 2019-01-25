package main

const _buff = "buff"
const _bane = "bane"
const _phys = "phys"
const _self = "self"

func getAttacks() []Attack {
	return []Attack{
		Attack{
			Name:        "basic",
			Accuracy:    0,
			Power:       0,
			Stat:        "hp",
			Targets:     1,
			FatigueCost: 1,
			Type:        _phys,
		}, Attack{
			Name:        "rest",
			Accuracy:    100,
			Power:       -2,
			Stat:        "fatigue",
			Targets:     1,
			FatigueCost: 1,
			Type:        _self,
		}, Attack{
			Name:        "big",
			Accuracy:    -1,
			Power:       3,
			Stat:        "hp",
			Targets:     1,
			FatigueCost: 2,
			Type:        _phys,
		}, Attack{
			Name:        "small",
			Accuracy:    1,
			Power:       -1,
			Stat:        "hp",
			Targets:     1,
			FatigueCost: 1,
			Type:        _phys,
		}, Attack{
			Name:        "burden",
			Accuracy:    1,
			Power:       2,
			Stat:        "fatigue",
			Targets:     1,
			FatigueCost: 1,
			Type:        _bane,
		}, Attack{
			Name:        "heal 1",
			Accuracy:    1,
			Power:       -3,
			Stat:        "hp",
			Targets:     1,
			FatigueCost: 1,
			Type:        _buff,
		}, Attack{
			Name:        "drain",
			Accuracy:    1,
			Power:       -2,
			Stat:        "select",
			Targets:     1,
			FatigueCost: 1,
			Type:        _bane,
		}, Attack{
			Name:        "buff",
			Accuracy:    1,
			Power:       1,
			Stat:        "select",
			Targets:     1,
			FatigueCost: 1,
			Type:        _buff,
		},
	}
}
