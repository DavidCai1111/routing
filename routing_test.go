package routing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefine(t *testing.T) {
	assert := assert.New(t)

	t.Run("should create a root node", func(t *testing.T) {
		root := New()

		assert.Nil(root.parent)
		assert.NotNil(root.children)

		root.Define("/")

		assert.Equal(1, len(root.children))

		assert.Equal(root, root.children[option{str: ""}].parent)
	})

	t.Run("should panic when url does not start with /", func(t *testing.T) {
		root := New()

		assert.Panics(func() {
			root.Define("test")
		})
	})

	t.Run("should create the first level children", func(t *testing.T) {
		root := New()

		root.Define("/test1")
		root.Define("/test2")

		assert.Equal(2, len(root.children))

		assert.Equal(root, root.children[option{str: "test1"}].parent)
		assert.Equal(root, root.children[option{str: "test2"}].parent)
	})

	t.Run("should create the two level children", func(t *testing.T) {
		root := New()

		root.Define("/test1/test2")
		root.Define("/test1/test3")
		root.Define("/test1/:_id(test3|test4)")

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

		root.Define("/(" + regStr + ")")
		root.Define("/:_id(" + regStr + ")")

		assert.Equal(2, len(root.children))

		assert.Equal(root, root.children[option{reg: regStr}].parent)
		assert.Equal(root, root.children[option{reg: regStr, name: "_id"}].parent)
	})

	t.Run("should create a separated strings url", func(t *testing.T) {
		root := New()

		root.Define("/test1|test2")
		root.Define("/:_id(test1|test2)")

		assert.Equal(4, len(root.children))

		assert.Equal(root, root.children[option{str: "test1"}].parent)
		assert.Equal(root, root.children[option{str: "test2"}].parent)
		assert.Equal(root, root.children[option{str: "test1", name: "_id"}].parent)
		assert.Equal(root, root.children[option{str: "test2", name: "_id"}].parent)
	})
}
