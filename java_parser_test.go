package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJavaClass_givenValidClassDeclaration_returnsOkAndClass(t *testing.T) {
	ok, className, _ := GetJavaClass("public class  Foo   {")
	assert.True(t, ok)
	assert.Equal(t, "Foo", className)
}

func TestGetJavaClass_givenClassThatExtendsOtherClass_returnsOkClassAndParentClass(t *testing.T) {
	ok, className, parentClassNames := GetJavaClass("public class Foo extends  Bar   {")
	assert.True(t, ok)
	assert.Equal(t, "Foo", className)
	assert.Equal(t, 1, len(parentClassNames))
	assert.Contains(t, parentClassNames, "Bar")
}

func TestGetJavaClass_givenClassThatImplementsInterface_returnsOkClassAndInterface(t *testing.T) {
	ok, className, parentClassNames := GetJavaClass("public class Foo implements  Bar   {")
	assert.True(t, ok)
	assert.Equal(t, "Foo", className)
	assert.Equal(t, 1, len(parentClassNames))
	assert.Contains(t, parentClassNames, "Bar")
}

func TestGetJavaClass_givenClassThatExtendsOtherClassAndImplementsInterface_returnsOkClassAndParentAndInterfaces(t *testing.T) {
	ok, className, parentClassNames := GetJavaClass("public class Foo extends  Bar   implements Baz , Quux  {")
	assert.True(t, ok)
	assert.Equal(t, "Foo", className)
	assert.Equal(t, 3, len(parentClassNames))
	assert.Contains(t, parentClassNames, "Bar")
	assert.Contains(t, parentClassNames, "Baz")
	assert.Contains(t, parentClassNames, "Quux")
}

func TestGetJavaClass_givenClassThatImplementsInterfaceAndExtendsOtherClass_returnsOkClassAndParentAndInterfaces(t *testing.T) {
	ok, className, parentClassNames := GetJavaClass("public class Foo  implements Baz , Quux  extends  Bar   {")
	assert.True(t, ok)
	assert.Equal(t, "Foo", className)
	assert.Equal(t, 3, len(parentClassNames))
	assert.Contains(t, parentClassNames, "Bar")
	assert.Contains(t, parentClassNames, "Baz")
	assert.Contains(t, parentClassNames, "Quux")
}
