package main

func getAttacks() []Attack {
	return []Attack{
		Attack{
			Name:        "basic",
			Stat:        "hp",
			Targets:     1,
			FatigueCost: 1,
			EffType:     PHYS,
		}, Attack{
			Name:        "rest",
			Accuracy:    100,
			PowerMod:    -2,
			Stat:        "fatigue",
			Targets:     1,
			FatigueCost: 1,
			EffType:     SELF,
		}, Attack{
			Name:        "big",
			Accuracy:    -1,
			PowerMod:    3,
			Stat:        "hp",
			Targets:     1,
			FatigueCost: 2,
			EffType:     PHYS,
		}, Attack{
			Name:        "small",
			Accuracy:    1,
			PowerMod:    -1,
			Stat:        "hp",
			Targets:     1,
			FatigueCost: 1,
			EffType:     PHYS,
		}, Attack{
			Name:        "burden",
			Accuracy:    1,
			PowerMod:    0,
			Stat:        "fatigue",
			Targets:     1,
			FatigueCost: 1,
			EffType:     MAG,
		}, Attack{
			Name:        "heal 1",
			Accuracy:    1,
			PowerMod:    -10,
			Stat:        "hp",
			Targets:     1,
			FatigueCost: 1,
			EffType:     MAG,
		}, Attack{
			Name:        "drain",
			Accuracy:    1,
			PowerMod:    3,
			Stat:        "select",
			Targets:     1,
			FatigueCost: 1,
			EffType:     MAG,
		}, Attack{
			Name:        "buff",
			Accuracy:    1,
			PowerMod:    -1,
			Stat:        "select",
			Targets:     1,
			FatigueCost: 1,
			EffType:     MAG,
		},
	}
}
