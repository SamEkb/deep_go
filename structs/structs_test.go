package structs

import (
	"math"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
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
		setBits(person, 0, 10, mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, 10, 10, health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, 20, 4, respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, 24, 4, strength)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, 28, 4, experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, 32, 4, level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, 36, 1, 1)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, 37, 1, 1)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, 38, 1, 1)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		setBits(person, 39, 2, personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	name          [42]byte
	data          [6]byte
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
	return getBits(p, 0, 10)
}

func (p *GamePerson) Health() int {
	return getBits(p, 10, 10)
}

func (p *GamePerson) Respect() int {
	return getBits(p, 20, 4)
}

func (p *GamePerson) Strength() int {
	return getBits(p, 24, 4)
}

func (p *GamePerson) Experience() int {
	return getBits(p, 28, 4)
}

func (p *GamePerson) Level() int {
	return getBits(p, 32, 4)
}

func (p *GamePerson) HasHouse() bool {
	return getBits(p, 36, 1) == 1
}

func (p *GamePerson) HasGun() bool {
	return getBits(p, 37, 1) == 1
}

func (p *GamePerson) HasFamilty() bool {
	return getBits(p, 38, 1) == 1
}

func (p *GamePerson) Type() int {
	return getBits(p, 39, 2)
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
