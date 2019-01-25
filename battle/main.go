package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	player := Unit{}
	player.Name = read("name:")
	player.Stats.Power = readAsInt("power:")
	player.Stats.Speed = readAsInt("speed:")
	player.Stats.Accuracy = readAsInt("acc:")
	player.Stats.Defense = readAsInt("def:")

	player.MaxHp = player.Stats.Defense
	player.Hp = player.Stats.Defense

	player.Attacks = append(player.Attacks, Attack{
		Name: "basic",
		Type: "Dmg",
	}, Attack{
		Name:  "buffAll",
		Type:  "Buff",
		Stats: Stats{1, 1, 1, 1},
	})

	for i := 0; i < 100; i++ {
		enemy := genUnit(i)
		fmt.Println("you are fighting", enemy.Name)
		fmt.Println("your stats", player.Stats)
		fmt.Println("your attacks", player.Attacks)
		fmt.Println("enemey", enemy)
		player, enemy = battle(player, enemy)
		if player.Hp == 0 {
			fmt.Println("dead!")
			read("")
			os.Exit(0)
		}

		fmt.Println("killed", enemy.Name)
		player.MaxHp++
		player.Hp = player.MaxHp
		read("enter to continue")

	}
}

func battle(player Unit, enemy Unit) (Unit, Unit) {
	player, enemy = playerTurn(player, enemy)
	if player.Hp <= 0 || enemy.Hp <= 0 {
		return player, enemy
	}
	enemy, player = enemyTurn(enemy, player)
	if player.Hp <= 0 || enemy.Hp <= 0 {
		return player, enemy
	}

	return battle(player, enemy)
}

func playerTurn(player Unit, enemy Unit) (Unit, Unit) {
	atk := chooseAttack(player)

	return applyAttack(player, enemy, atk)
}

func enemyTurn(enemy Unit, player Unit) (Unit, Unit) {
	atk := Attack{
		Stats: enemy.Stats,
		Type:  "Dmg",
	}
	return applyAttack(enemy, player, atk)
}

func chooseAttack(unit Unit) Attack {
	atkString := read("atk")
	atak := Attack{}
	for _, atk := range unit.Attacks {
		if atk.Name == atkString {
			atak = atk
			break
		}

	}

	if atak.Name == "" {
		fmt.Println("atk not found")
		return chooseAttack(unit)
	}

	atak.Stats = Stats{
		Power:    atak.Stats.Power + unit.Stats.Power,
		Speed:    atak.Stats.Speed + unit.Stats.Speed,
		Defense:  atak.Stats.Defense + unit.Stats.Defense,
		Accuracy: atak.Stats.Accuracy + unit.Stats.Accuracy,
	}

	return atak
}

func genUnit(i int) Unit {

	return Unit{
		Name:  fmt.Sprintf("goblin%v", i),
		MaxHp: i,
		Hp:    i,
		Stats: Stats{
			i, i, i, i,
		},
	}
}

func read(msg string) string {
	fmt.Println(msg)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func readAsInt(msg string) int {
	str := read(msg)

	i, err := strconv.Atoi(str)
	if err != nil {
		return readAsInt("invalid int, try again")
	}

	return i

}

func applyAttack(attacker Unit, defender Unit, attack Attack) (Unit, Unit) {
	attackStats := attack.Stats
	if attack.Type == "Buff" {
		attacker.Stats.Accuracy += attackStats.Accuracy
		attacker.Stats.Defense += attackStats.Defense
		attacker.Stats.Power += attackStats.Power
		attacker.Stats.Speed += attackStats.Speed
	} else {
		if attackStats.Accuracy < attacker.Stats.Speed {
			fmt.Println(attacker.Name, "missed")
			return attacker, defender
		}

		if attack.Type == "Dmg" {
			dmg := attackStats.Power - defender.Stats.Defense
			if dmg < 0 {
				dmg = 0
			}
			fmt.Println(attacker.Name, "dealt", dmg)
			defender.Hp -= dmg
			return attacker, defender
		}

	}

	return attacker, defender
}

type Unit struct {
	Name    string
	Stats   Stats
	XP      int
	Level   int
	Hp      int
	MaxHp   int
	Attacks []Attack
}

type Attack struct {
	Name  string
	Type  string
	Stats Stats
}

type Stats struct {
	Speed    int
	Power    int
	Defense  int
	Accuracy int
}
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
	"strength", "defense", "speed", "accuracy", "vitality",
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
		Strength: u.BaseAttributes.Strength + u.combatAttrMods.Strength,
		Speed:    u.BaseAttributes.Speed + u.combatAttrMods.Speed,
		Defense:  u.BaseAttributes.Defense + u.combatAttrMods.Defense,
		Accuracy: u.BaseAttributes.Accuracy + u.combatAttrMods.Accuracy,
		Vitality: u.BaseAttributes.Vitality + u.combatAttrMods.Vitality,
	}

	u.currentAttributes.Speed -= u.Fatigue
	u.currentAttributes.Defense -= int(u.Fatigue / 2)
	u.currentAttributes.Accuracy -= int(u.Fatigue / 3)
	u.currentAttributes.Strength -= int(u.Fatigue / 4)

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
)

type Attack struct {
	Name        string
	FatigueCost int
	Power       int
	Accuracy    int
	Targets     int
	Stat        string
	EffType     EffectType
	Type        string
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
	fmt.Printf("Chosen Attack: Name: %v. Stat: %v. Pow: %v. Acc: %v. NumTargets: %v. FatCost: %v.\r\n", atk.Name, atk.Stat, atk.Power, atk.Accuracy, atk.Targets, atk.FatigueCost)
	fmt.Println("")
}

func createPlayer() Unit {
	name := read("name")
	strong := readAttr("good")
	weak := readAttr("bad")

	attrs := Attributes{
		5, 5, 5, 5, 5,
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

	attrs := Attributes{offset(i), offset(i), offset(i), offset(i), offset(i)}

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
	if atk.Type == _self {
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
	dmg := 0

	switch attack.Type {
	case _phys:
		if attack.Accuracy < unit.Speed() {
			return AttackResult{
				Attr: "miss",
			}
		}

		dmg = attack.Power - unit.Defense()
		if dmg < 0 {
			dmg = 0
		}
	case _buff, _self:
		dmg = attack.Power
	case _bane:
		if attack.Accuracy < unit.Speed() {
			return AttackResult{
				Attr: "miss",
			}
		}
		dmg = attack.Power - unit.Defense()
	}

	return AttackResult{
		Attr:    attack.Stat,
		Damange: dmg,
	}

}

func (u Unit) getModifiers() AttackMod {
	return AttackMod{
		u.Strength(),
		u.Accuracy(),
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
	if atk.EffType == Phys {
		mods := unit.getModifiers()
		atk.Power += mods.PowerMod
		atk.Accuracy += mods.AccMod
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
			targets = append(targets, findFirstValidTarget(atk.Type, team, units))
		}

	}

	return targets

}

func findFirstValidTarget(typ string, team int, units []Unit) Unit {
	for _, unit := range units {
		switch typ {
		case _phys, _bane:
			if unit.Team != team {
				return unit
			}
		case _buff:
			if unit.Team == team {
				return unit
			}
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
