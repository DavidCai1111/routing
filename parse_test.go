package routing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	t.Run("should parse empty string", func(t *testing.T) {
		options := parse("")
		assert.Equal([]option{option{str: ""}}, options)
	})

	t.Run("should parse a simple string", func(t *testing.T) {
		options := parse("test")

		assert.Equal([]option{option{str: "test"}}, options)
	})

	t.Run("should parse a string including _ and -", func(t *testing.T) {
		options := parse("t-es_t")

		assert.Equal([]option{option{str: "t-es_t"}}, options)
	})

	t.Run("should parse a reg", func(t *testing.T) {
		regStr := `^[\w\.-]+$`
		options := parse("(" + regStr + ")")

		assert.Equal(1, len(options))
		assert.Equal("", options[0].name)
		assert.Equal(regStr, options[0].reg)
	})

	t.Run("should parse a named param", func(t *testing.T) {
		options := parse(":_id")

		assert.Equal([]option{option{name: "_id"}}, options)
	})

	t.Run("should parse a named param with strings", func(t *testing.T) {
		options := parse(":_id(test1|test2)")

		assert.Equal([]option{
			option{name: "_id", str: "test1"},
			option{name: "_id", str: "test2"},
		}, options)
	})

	t.Run("should parse a named param with regex", func(t *testing.T) {
		regStr := `^[\w\.-]+$`
		options := parse(":_id(" + regStr + ")")

		assert.Equal(1, len(options))
		assert.Equal("_id", options[0].name)
		assert.Equal(regStr, options[0].reg)
	})

	t.Run("should parse a named param with strings", func(t *testing.T) {
		options := parse("test1|test2")

		assert.Equal([]option{
			option{str: "test1"},
			option{str: "test2"},
		}, options)
	})

	t.Run("should panic on invalid strings", func(t *testing.T) {
		assert.Panics(func() {
			parse("test1|$$$")
		})
	})
}
