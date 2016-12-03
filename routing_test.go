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

func TestDefine(t *testing.T) {
	assert := assert.New(t)

	t.Run("should create a root node", func(t *testing.T) {
		root := New()

		assert.Nil(root.parent)
		assert.NotNil(root.children)

		root.Define("/", nil)

		assert.Equal(1, len(root.children))

		assert.Equal(root, root.children[option{str: ""}].parent)
	})

	t.Run("should panic when url does not start with /", func(t *testing.T) {
		root := New()

		assert.Panics(func() {
			root.Define("test", nil)
		})
	})

	t.Run("should create the first level children", func(t *testing.T) {
		root := New()

		root.Define("/test1", nil)
		root.Define("/test2", nil)

		assert.Equal(2, len(root.children))

		assert.Equal(root, root.children[option{str: "test1"}].parent)
		assert.Equal(root, root.children[option{str: "test2"}].parent)
	})

	t.Run("should create the two level children", func(t *testing.T) {
		root := New()

		root.Define("/test1/test2", nil)
		root.Define("/test1/test3", nil)
		root.Define("/test1/:_id(test3|test4)", nil)

		assert.Equal(1, len(root.children))

		assert.Equal(root, root.children[option{str: "test1"}].parent)

		r := root.children[option{str: "test1"}]

		assert.Equal(4, len(r.children))

		assert.Equal(r, r.children[option{str: "test2"}].parent)
		assert.Equal(r, r.children[option{str: "test3"}].parent)
		assert.Equal(r, r.children[option{str: "test3", name: "_id"}].parent)
		assert.Equal(r, r.children[option{str: "test4", name: "_id"}].parent)
	})

	t.Run("should create a reg url", func(t *testing.T) {
		root := New()
		regStr := `^[\w\.-]+$`

		root.Define("/("+regStr+")", nil)
		root.Define("/:_id("+regStr+")", nil)

		assert.Equal(2, len(root.children))

		assert.Equal(root, root.children[option{reg: regStr}].parent)
		assert.Equal(root, root.children[option{reg: regStr, name: "_id"}].parent)
	})

	t.Run("should create a separated strings url", func(t *testing.T) {
		root := New()

		root.Define("/test1|test2", nil)
		root.Define("/:_id(test1|test2)", nil)

		assert.Equal(4, len(root.children))

		assert.Equal(root, root.children[option{str: "test1"}].parent)
		assert.Equal(root, root.children[option{str: "test2"}].parent)
		assert.Equal(root, root.children[option{str: "test1", name: "_id"}].parent)
		assert.Equal(root, root.children[option{str: "test2", name: "_id"}].parent)
	})
}

func TestMatch(t *testing.T) {
	assert := assert.New(t)

	t.Run("should match /", func(t *testing.T) {
		root := New()

		root.Define("/", 1)
		c, p, ok := root.Match("/")

		assert.True(ok)
		assert.Zero(len(p))
		assert.Equal(1, c.(int))
	})

	t.Run("should match one level url", func(t *testing.T) {
		root := New()

		root.Define("/fav.icon", 1)
		c, p, ok := root.Match("/fav.icon")

		assert.True(ok)
		assert.Zero(len(p))
		assert.Equal(1, c.(int))
	})

	t.Run("should match named params", func(t *testing.T) {
		root := New()

		root.Define("/:_id", 1)
		c, p, ok := root.Match("/123")

		assert.True(ok)
		assert.Equal(1, len(p))
		assert.Equal("123", p["_id"])
		assert.Equal(1, c.(int))
	})

	t.Run("should match named params for regex", func(t *testing.T) {
		root := New()

		root.Define("/:_id(\\w{3,30})", 1)
		c, p, ok := root.Match("/haha")

		assert.True(ok)
		assert.Equal(1, len(p))
		assert.Equal("haha", p["_id"])
		assert.Equal(1, c.(int))
	})

	t.Run("should match string first", func(t *testing.T) {
		root := New()

		root.Define("/haha", 1)
		root.Define("/(\\w{3,30})", 2)
		c, p, ok := root.Match("/haha")

		assert.True(ok)
		assert.Equal(0, len(p))
		assert.Equal(1, c.(int))
	})

	t.Run("should match encoded url", func(t *testing.T) {
		root := New()

		root.Define("/:_id(@\\w+)", 1)
		c, p, ok := root.Match("/%40haha")

		assert.True(ok)
		assert.Equal(1, len(p))
		assert.Equal("@haha", p["_id"])
		assert.Equal(1, c.(int))
	})
}
