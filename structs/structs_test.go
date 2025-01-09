package structs

import (
	"math"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	manaOffset = 0
	manaLength = 10

	healthOffset = 10
	healthLength = 10

	respectOffset = 20
	respectLength = 4

	strengthOffset = 24
	strengthLength = 4

	experienceOffset = 28
	experienceLength = 4

	levelOffset = 32
	levelLength = 4

	houseOffset = 36
	houseLength = 1

	gunOffset = 37
	gunLength = 1

	familyOffset = 38
	familyLength = 1

	typeOffset = 39
	typeLength = 2
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, manaOffset, manaLength, mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, healthOffset, healthLength, health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, respectOffset, respectLength, respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, strengthOffset, strengthLength, strength)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, experienceOffset, experienceLength, experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, levelOffset, levelLength, level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, houseOffset, houseLength, 1)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, gunOffset, gunLength, 1)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, familyOffset, familyLength, 1)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, typeOffset, typeLength, personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	//  Строка имени, хранит копию имени пользователя (до 42 символов)
	name [42]byte

	//  Служебный буфер, куда побитово упакованы поля:
	//    offset=0..9 для Mana (10 бит)
	//    offset=10..19 для Health (10 бит)
	//    offset=20..23 для Respect (4 бита)
	//    offset=24..27 для Strength (4 бита)
	//    offset=28..31 для Experience (4 бита)
	//    offset=32..35 для Level (4 бита)
	//    offset=36 для House (1 бит)
	//    offset=37 для Gun (1 бит)
	//    offset=38 для Family (1 бит)
	//    offset=39..40 для Type (2 бита)
	data [6]byte

	// Поля для координат и золота
	x, y, z, gold int32
}

func setBits(p *GamePerson, offset, length, value int) {
	var v uint64
	for i := 0; i < 6; i++ {
		v |= uint64(p.data[i]) << (8 * i)
	}

	mask := (uint64(1) << length) - 1
	v &^= mask << offset

	v |= (uint64(value) & mask) << offset
	for i := 0; i < 6; i++ {
		p.data[i] = byte(v >> (8 * i))
	}
}

func getBits(p *GamePerson, offset, length int) int {
	var v uint64
	for i := 0; i < 6; i++ {
		v |= uint64(p.data[i]) << (8 * i)
	}

	mask := (uint64(1) << length) - 1
	return int((v >> offset) & mask)
}

func NewGamePerson(options ...Option) GamePerson {
	var person GamePerson
	for _, opt := range options {
		opt(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	return strings.TrimRight(string(p.name[:]), "\x00")
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return getBits(p, manaOffset, manaLength)
}

func (p *GamePerson) Health() int {
	return getBits(p, healthOffset, healthLength)
}

func (p *GamePerson) Respect() int {
	return getBits(p, respectOffset, respectLength)
}

func (p *GamePerson) Strength() int {
	return getBits(p, strengthOffset, strengthLength)
}

func (p *GamePerson) Experience() int {
	return getBits(p, experienceOffset, experienceLength)
}

func (p *GamePerson) Level() int {
	return getBits(p, levelOffset, levelLength)
}

func (p *GamePerson) HasHouse() bool {
	return getBits(p, houseOffset, houseLength) == 1
}

func (p *GamePerson) HasGun() bool {
	return getBits(p, gunOffset, gunLength) == 1
}

func (p *GamePerson) HasFamilty() bool {
	return getBits(p, familyOffset, familyLength) == 1
}

func (p *GamePerson) Type() int {
	return getBits(p, typeOffset, typeLength)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
