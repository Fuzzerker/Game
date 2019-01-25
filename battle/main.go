package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var validAttrs = []string{
	"strength", "defense", "speed", "accuracy", "vitality", "resistance", "willpower",
}

func main() {
	rand.Seed(time.Now().UnixNano())
	Play()
}

type Unit struct {
	Name              string
	Type              string
	BaseAttributes    Attributes
	combatAttrMods    Attributes
	currentAttributes Attributes
	BaseStats         BaseStats
	Hp                int
	Fatigue           int
	IsHuman           bool
	Attacks           []Attack
	Team              int
	AiLevel           int
}

func (u Unit) Crunch() Unit {
	u.currentAttributes = Attributes{
		Strength:   u.BaseAttributes.Strength + u.combatAttrMods.Strength,
		Speed:      u.BaseAttributes.Speed + u.combatAttrMods.Speed,
		Defense:    u.BaseAttributes.Defense + u.combatAttrMods.Defense,
		Accuracy:   u.BaseAttributes.Accuracy + u.combatAttrMods.Accuracy,
		Vitality:   u.BaseAttributes.Vitality + u.combatAttrMods.Vitality,
		Willpower:  u.BaseAttributes.Willpower + u.combatAttrMods.Willpower,
		Resistance: u.BaseAttributes.Resistance + u.combatAttrMods.Resistance,
	}

	u.currentAttributes.Speed -= u.Fatigue
	u.currentAttributes.Defense -= int(u.Fatigue / 2)
	u.currentAttributes.Accuracy -= int(u.Fatigue / 3)
	u.currentAttributes.Strength -= int(u.Fatigue / 4)
	u.currentAttributes.Willpower -= u.Fatigue
	u.currentAttributes.Resistance -= (u.Fatigue - u.currentAttributes.Willpower)

	u.BaseStats.MaxHp = u.currentAttributes.Vitality + int(u.currentAttributes.Defense/2)

	return u
}

func (u Unit) ModStrength(i int) Unit {
	u.combatAttrMods.Strength += i
	return u
}

func (u Unit) ModDefense(i int) Unit {
	u.combatAttrMods.Defense += i
	return u
}
func (u Unit) ModAccuracy(i int) Unit {
	u.combatAttrMods.Accuracy += i
	return u
}
func (u Unit) ModVitality(i int) Unit {
	u.combatAttrMods.Vitality += i
	return u
}
func (u Unit) ModSpeed(i int) Unit {
	u.combatAttrMods.Speed += i
	return u
}
func (u Unit) ModWillpower(i int) Unit {
	u.combatAttrMods.Willpower += i
	return u
}
func (u Unit) ModResistance(i int) Unit {
	u.combatAttrMods.Resistance += i
	return u
}
func (u Unit) Speed() int {
	return u.Crunch().currentAttributes.Speed
}

func (u Unit) Defense() int {
	return u.Crunch().currentAttributes.Defense
}

func (u Unit) Strength() int {
	return u.Crunch().currentAttributes.Strength
}

func (u Unit) Accuracy() int {
	return u.Crunch().currentAttributes.Accuracy
}

func (u Unit) Vitality() int {
	return u.Crunch().currentAttributes.Vitality
}

func (u Unit) Willpower() int {
	return u.Crunch().currentAttributes.Willpower
}

func (u Unit) Resistance() int {
	return u.Crunch().currentAttributes.Resistance
}

type BaseStats struct {
	MaxHp int
}

type Attributes struct {
	Strength   int
	Defense    int
	Speed      int
	Accuracy   int
	Vitality   int
	Willpower  int
	Resistance int
}

type EffectType int

const (
	PHYS EffectType = iota
	MAG
	SELF
)

type Attack struct {
	Name        string
	FatigueCost int
	PowerMod    int
	Accuracy    int
	Targets     int
	Stat        string
	EffType     EffectType
	Team        int
}

type AttackMod struct {
	PowerMod int
	AccMod   int
}

type AttackResult struct {
	Damange int
	Attr    string
}

func Play() {
	player := createPlayer()

	for fight := 1; ; fight++ {
		fmt.Println("starting fight", fight)
		enemies := genEnemies(fight, 1)
		combatants := enemies
		combatants = append(combatants, player)

		for round := 0; ; round++ {
			fmt.Println("round", round)
			combatants = playRound(combatants)
			for _, c := range combatants {
				if c.Name == player.Name {
					player = c
					break
				}
			}

			if player.Hp <= 0 {
				fmt.Println("you ded, try again")
				os.Exit(0)
			}

			if checkOver(combatants) {
				read("round over")
				break
			}
		}

	}

}

func playRound(combatants []Unit) []Unit {
	orderedUnitNames := getPlayOrder(combatants)

	for _, name := range orderedUnitNames {
		if checkOver(combatants) {
			return combatants
		}
		for _, c := range combatants {
			if c.Name == name {
				printUnit(c)
				combatants = Turn(c, combatants)
				break
			}

		}
	}

	return combatants

}

func checkOver(cbts []Unit) bool {
	teams := make(map[int]bool)

	for _, cbt := range cbts {
		if cbt.Hp > 0 {
			teams[cbt.Team] = true
		}
	}

	if len(teams) < 2 {
		return true
	}

	return false
}

type order struct {
	Name  string
	Speed int
}

func getPlayOrder(units []Unit) []string {

	var orders []order
	for _, u := range units {
		orders = append(orders, order{u.Name, u.Speed()})

	}

	return handleOrders(orders)

}

func handleOrders(orders []order) []string {

	ordered := []string{}

	for {
		if allAccounted(orders, ordered) {
			break
		}

		next := getNext(orders)
		ordered = append(ordered, next)
		for i, order := range orders {
			if order.Name == next {
				order.Speed -= 5
				orders[i] = order
			}

		}

	}
	return ordered

}

func allAccounted(orders []order, list []string) bool {

	for _, o := range orders {
		if !contains(o.Name, list) {
			return false
		}
	}
	return true
}

func getNext(orders []order) string {
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Speed > orders[j].Speed
	})

	return orders[0].Name

}

func printUnit(unit Unit) {
	fmt.Printf("Name: %v. Hp: %v/%v. Lvl: %v. Team: %v. Fat: %v. CurAttrs[Str: %v, Def: %v, Spd: %v, Acc: %v, Vit: %v]\r\n", unit.Name, unit.Hp, unit.BaseStats.MaxHp,
		unit.AiLevel, unit.Team,
		unit.Fatigue, unit.Strength(), unit.Defense(), unit.Speed(), unit.Accuracy(), unit.Vitality())
}

func printAttack(atk Attack) {
	fmt.Println("")
	fmt.Printf("Chosen Attack: Name: %v. Stat: %v. Pow: %v. Acc: %v. NumTargets: %v. FatCost: %v.\r\n", atk.Name, atk.Stat, atk.PowerMod, atk.Accuracy, atk.Targets, atk.FatigueCost)
	fmt.Println("")
}

func createPlayer() Unit {
	name := read("name")
	strong := readAttr("good")
	weak := readAttr("bad")

	attrs := Attributes{
		5, 5, 5, 5, 5, 5, 5,
	}

	for _, attr := range validAttrs {
		if attr == strong {
			switch attr {
			case "strength":
				attrs.Strength += 5
			case "defense":
				attrs.Defense += 5
			case "speed":
				attrs.Speed += 5
			case "accuracy":
				attrs.Accuracy += 5
			case "vitality":
				attrs.Vitality += 5

			}

		}
		if attr == weak {
			switch attr {
			case "strength":
				attrs.Strength -= 5
			case "defense":
				attrs.Defense -= 5
			case "speed":
				attrs.Speed -= 5
			case "accuracy":
				attrs.Accuracy -= 5
			case "vitality":
				attrs.Vitality -= 5

			}
		}
	}

	player := CreateUnit(name, attrs)
	player.Attacks = getAttacks()
	player.IsHuman = true
	return player

}

func CreateUnit(name string, attrs Attributes) Unit {

	unit := Unit{
		Name:           name,
		BaseAttributes: attrs,
		Attacks:        []Attack{getAttacks()[0]},
	}

	unit = unit.Crunch()
	unit.Hp = unit.BaseStats.MaxHp

	return unit

}

func readAttr(msg string) string {
	raw := read(msg)

	if !contains(raw, validAttrs) {
		fmt.Println("invalid attr")
		return readAttr(msg)
	}
	return raw
}

func contains(target string, strs []string) bool {
	for _, s := range strs {
		if s == target {
			return true
		}
	}

	return false
}

var abs = 0

func genEnemies(i int, team int) []Unit {

	units := []Unit{}

	count := 1
	if i > 5 {
		count++
	}
	if i > 10 {
		count++
	}
	if i > 15 {
		count++
	}

	for x := 0; x < count; x++ {
		units = append(units, genUnit(x, i, team))

	}

	return units

}

var atk = 0

func genUnit(x int, i int, team int) Unit {
	potentialAtks := getAttacks()

	attrs := Attributes{offset(i), offset(i), offset(i), offset(i), offset(i), offset(i), offset(i)}

	unit := CreateUnit(fmt.Sprintf("goblin%v-%v-%v", i, x, abs), attrs)
	unit.Team = team
	unit.AiLevel = randomInt(0, 10)

	if i%3 == 0 && atk < len(potentialAtks) {
		unit.Attacks = append(unit.Attacks, potentialAtks[atk])
		atk++

	}

	abs++
	return unit

}

func offset(i int) int {
	return i + (int(i/2) * randomInt(-1, 1))
}

func Turn(active Unit, units []Unit) []Unit {
	if active.Hp <= 0 {
		fmt.Println(active.Name, "is dead")
		return units
	}
	atk := PickAttack(active)
	printAttack(atk)
	var targets []Unit
	if atk.EffType == SELF {
		targets = []Unit{active}
	} else {
		targets = PickTargets(atk, active.AiLevel, active.Team, active.IsHuman, units)
	}

	hitMap := make(map[string][]AttackResult)
	for _, target := range targets {
		hitMap[target.Name] = append(hitMap[target.Name], resolveAttack(atk, target))
	}

	for i, unit := range units {
		for name, results := range hitMap {
			if name == unit.Name {
				for _, result := range results {
					switch result.Attr {
					case "strength":
						unit = unit.ModStrength(result.Damange)
					case "defense":
						unit = unit.ModDefense(result.Damange)
					case "speed":
						unit = unit.ModSpeed(result.Damange)
					case "accuracy":
						unit = unit.ModAccuracy(result.Damange)
					case "fatigue":
						unit.Fatigue += result.Damange
					case "hp":
						unit.Hp -= result.Damange
						if unit.Hp > unit.BaseStats.MaxHp {
							unit.Hp = unit.BaseStats.MaxHp
						}
					case "select":
						var attr string
						if active.IsHuman {
							attr = readAttr("pick stat")
						} else {
							attr = validAttrs[randomInt(0, len(validAttrs)-1)]
						}
						result.Attr = attr

						switch attr {
						case "strength":
							unit = unit.ModStrength(result.Damange)
						case "defense":
							unit = unit.ModDefense(result.Damange)
						case "speed":
							unit = unit.ModSpeed(result.Damange)
						case "accuracy":
							unit = unit.ModAccuracy(result.Damange)
						case "vitality":
							unit = unit.ModVitality(result.Damange)
						}
					}

					if result.Attr == "miss" {
						fmt.Printf("%v attack %v missed %v\r\n", active.Name, atk.Name, unit.Name)
					} else {
						fmt.Printf("%v attack %v hit %v and dealt %v to %v\r\n", active.Name, atk.Name, unit.Name, result.Damange, result.Attr)
						if unit.Hp <= 0 {
							fmt.Println(active.Name, "killed", unit.Name)
						}
					}
				}

			}
		}

		units[i] = unit
	}

	activeIndex := -1

	for i, unit := range units {
		if unit.Name == active.Name {
			active = unit
			activeIndex = i
			break
		}
	}

	active.Fatigue += atk.FatigueCost
	active = active.Crunch()
	units[activeIndex] = active
	return units

}

func resolveAttack(attack Attack, unit Unit) AttackResult {
	//TODO: support magic attacks

	if attack.Team == unit.Team {
		return AttackResult{
			Attr:    attack.Stat,
			Damange: attack.PowerMod,
		}
	}

	if attack.EffType == PHYS {
		if attack.Accuracy < unit.Speed() {
			return AttackResult{
				Attr: "miss",
			}
		}

		dmg := attack.PowerMod - unit.Defense()
		if dmg < 0 {
			dmg = 0
		}

		return AttackResult{
			Attr:    attack.Stat,
			Damange: dmg,
		}
	}

	if attack.Accuracy < unit.Willpower() {
		return AttackResult{
			Attr: "miss",
		}
	}

	dmg := attack.PowerMod - unit.Resistance()
	if dmg < 0 {
		dmg = 0
	}

	return AttackResult{
		Attr:    attack.Stat,
		Damange: dmg,
	}

}

func PickTargets(atk Attack, lvl int, team int, human bool, units []Unit) []Unit {
	if human {
		fmt.Println("")
		for _, unit := range units {
			printUnit(unit)
		}
		fmt.Println("")
		return pickPlayerTargets(atk.Targets, units)
	}

	return npcPickTargets(atk, lvl, team, units)

}

func PickAttack(unit Unit) Attack {
	var atk Attack
	if unit.IsHuman {
		atk = pickPlayerAttack(unit.Attacks)
	} else {
		atk = npcPickAttack(unit)
	}
	atk.Team = unit.Team
	//TODO: phys vs mag
	if atk.EffType == PHYS {
		atk.PowerMod += unit.Strength()
		atk.Accuracy += unit.Accuracy()
	} else {
		atk.PowerMod += unit.Willpower()
		atk.Accuracy += int((unit.Willpower() + unit.Resistance()) / 2)
	}

	return atk
}

func pickPlayerTargets(num int, units []Unit) []Unit {
	var targets []Unit
	for i := 0; i < num; i++ {
		fmt.Println("Target", i)
		targets = append(targets, selectPlayerTarget(units))
	}

	return targets

}

func npcPickAttack(unit Unit) Attack {
	if unit.AiLevel < 3 {
		return unit.Attacks[randomInt(0, len(unit.Attacks)-1)]
	}

	if unit.Fatigue > 4 {
		for _, atk := range unit.Attacks {
			if atk.Name == "rest" {
				return atk
			}
		}
	}
	var basic Attack
	for _, atk := range unit.Attacks {
		if atk.Name == "big" && unit.Strength() < unit.Accuracy() {
			return atk
		}
		if atk.Name == "small" && unit.Accuracy() < unit.Strength() {
			return atk
		}
		if atk.Name == "basic" {
			basic = atk
		}
	}
	return basic

}

func npcPickTargets(atk Attack, lvl int, team int, units []Unit) []Unit {
	targets := []Unit{}
	for i := 0; i < atk.Targets; i++ {
		if lvl < 2 {
			randUnit := randomInt(0, len(units)-1)
			targets = append(targets, units[randUnit])
		} else {
			targets = append(targets, findFirstValidTarget(atk.PowerMod, team, units))
		}

	}

	return targets

}

func findFirstValidTarget(power int, team int, units []Unit) Unit {
	for _, unit := range units {
		if power > 0 {
			if unit.Team != team {
				return unit
			}
		} else if unit.Team == team {
			return unit
		}

	}

	if len(units) < 1 {
		return Unit{}
	}
	return units[0]

}

func selectPlayerTarget(units []Unit) Unit {
	chosen := read("pick target")
	target := Unit{}
	for _, unit := range units {
		if unit.Name == chosen {
			target = unit
		}
	}
	if target.Name == "" {
		fmt.Println("invalid target")
		return selectPlayerTarget(units)
	}
	return target
}
func pickPlayerAttack(atks []Attack) Attack {
	atkStr := read("choose attack")
	var atk Attack
	for _, ak := range atks {
		if ak.Name == atkStr {
			atk = ak
		}
	}
	if atk.Name == "" {
		fmt.Println("Ivalid attack")
		return pickPlayerAttack(atks)
	}
	return atk

}

func randomInt(min, max int) int {
	max++
	return min + rand.Intn(max-min)
}

func read(msg string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(msg)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
	return text
}

func readAsInt(msg string) int {
	raw := read(msg)
	i, err := strconv.Atoi(raw)
	if err != nil {
		readAsInt("invalid number, try again")
	}
	return i
}
