package utils

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_IntToString(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []string{"1", "2", "3", "4", "5"}

	result := Map(input, func(i int) string {
		return strconv.Itoa(i)
	})

	assert.Equal(t, expected, result)
	assert.Len(t, result, len(input))
}

func TestMap_StringToInt(t *testing.T) {
	input := []string{"10", "20", "30"}
	expected := []int{10, 20, 30}

	result := Map(input, func(s string) int {
		i, _ := strconv.Atoi(s)
		return i
	})

	assert.Equal(t, expected, result)
	assert.Len(t, result, len(input))
}

func TestMap_StringToUpper(t *testing.T) {
	input := []string{"hello", "world", "test"}
	expected := []string{"HELLO", "WORLD", "TEST"}

	result := Map(input, func(s string) string {
		return strings.ToUpper(s)
	})

	assert.Equal(t, expected, result)
}

func TestMap_EmptySlice(t *testing.T) {
	input := []int{}

	result := Map(input, func(i int) string {
		return strconv.Itoa(i)
	})

	assert.Empty(t, result)
	assert.NotNil(t, result) // Should return empty slice, not nil
	assert.Len(t, result, 0)
}

func TestMap_SingleElement(t *testing.T) {
	input := []int{42}
	expected := []string{"42"}

	result := Map(input, func(i int) string {
		return strconv.Itoa(i)
	})

	assert.Equal(t, expected, result)
	assert.Len(t, result, 1)
}

func TestMap_ComplexTransformation(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	type PersonSummary struct {
		Info    string
		IsAdult bool
	}

	input := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 17},
		{Name: "Charlie", Age: 30},
	}

	expected := []PersonSummary{
		{Info: "Alice (25)", IsAdult: true},
		{Info: "Bob (17)", IsAdult: false},
		{Info: "Charlie (30)", IsAdult: true},
	}

	result := Map(input, func(p Person) PersonSummary {
		return PersonSummary{
			Info:    p.Name + " (" + strconv.Itoa(p.Age) + ")",
			IsAdult: p.Age >= 18,
		}
	})

	assert.Equal(t, expected, result)
	assert.Len(t, result, len(input))
}

func TestMap_BooleanTransformation(t *testing.T) {
	input := []int{0, 1, 2, 0, 5}
	expected := []bool{false, true, true, false, true}

	result := Map(input, func(i int) bool {
		return i != 0
	})

	assert.Equal(t, expected, result)
}

func TestMap_SliceLength(t *testing.T) {
	input := []string{"a", "ab", "abc", "abcd"}
	expected := []int{1, 2, 3, 4}

	result := Map(input, func(s string) int {
		return len(s)
	})

	assert.Equal(t, expected, result)
}

func TestMap_NilFunction(t *testing.T) {
	input := []int{1, 2, 3}
	var nilFunc func(int) string

	// This should panic since fn is nil
	assert.Panics(t, func() {
		Map(input, nilFunc)
	})
}

func TestMap_FloatToInt(t *testing.T) {
	input := []float64{1.1, 2.7, 3.9, 4.2}
	expected := []int{1, 2, 3, 4}

	result := Map(input, func(f float64) int {
		return int(f)
	})

	assert.Equal(t, expected, result)
}

func TestReduce_Sum(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := 15

	result := Reduce(input, func(acc, val int) int {
		return acc + val
	}, 0)

	assert.Equal(t, expected, result)
}

func TestReduce_Product(t *testing.T) {
	input := []int{2, 3, 4}
	expected := 24

	result := Reduce(input, func(acc, val int) int {
		return acc * val
	}, 1)

	assert.Equal(t, expected, result)
}

func TestReduce_StringConcatenation(t *testing.T) {
	input := []string{"Hello", " ", "World", "!"}
	expected := "Hello World!"

	result := Reduce(input, func(acc, val string) string {
		return acc + val
	}, "")

	assert.Equal(t, expected, result)
}

func TestReduce_EmptySlice(t *testing.T) {
	input := []int{}
	initial := 42

	result := Reduce(input, func(acc, val int) int {
		return acc + val
	}, initial)

	assert.Equal(t, initial, result)
}

func TestReduce_SingleElement(t *testing.T) {
	input := []int{10}
	initial := 5
	expected := 15

	result := Reduce(input, func(acc, val int) int {
		return acc + val
	}, initial)

	assert.Equal(t, expected, result)
}

func TestReduce_Max(t *testing.T) {
	input := []int{3, 7, 2, 9, 1, 8}
	expected := 9

	result := Reduce(input, func(acc, val int) int {
		if val > acc {
			return val
		}
		return acc
	}, input[0])

	assert.Equal(t, expected, result)
}

func TestReduce_Min(t *testing.T) {
	input := []int{3, 7, 2, 9, 1, 8}
	expected := 1

	result := Reduce(input, func(acc, val int) int {
		if val < acc {
			return val
		}
		return acc
	}, input[0])

	assert.Equal(t, expected, result)
}

func TestReduce_CountingCondition(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := 5 // Count of even numbers

	result := Reduce(input, func(acc, val int) int {
		if val%2 == 0 {
			return acc + 1
		}
		return acc
	}, 0)

	assert.Equal(t, expected, result)
}

func TestReduce_ComplexAccumulation(t *testing.T) {
	type Item struct {
		Name  string
		Price float64
	}

	type Summary struct {
		TotalPrice float64
		ItemCount  int
		Names      []string
	}

	input := []Item{
		{Name: "Apple", Price: 1.50},
		{Name: "Banana", Price: 0.75},
		{Name: "Orange", Price: 2.00},
	}

	expected := Summary{
		TotalPrice: 4.25,
		ItemCount:  3,
		Names:      []string{"Apple", "Banana", "Orange"},
	}

	result := Reduce(input, func(acc Summary, item Item) Summary {
		return Summary{
			TotalPrice: acc.TotalPrice + item.Price,
			ItemCount:  acc.ItemCount + 1,
			Names:      append(acc.Names, item.Name),
		}
	}, Summary{})

	assert.Equal(t, expected.TotalPrice, result.TotalPrice)
	assert.Equal(t, expected.ItemCount, result.ItemCount)
	assert.Equal(t, expected.Names, result.Names)
}

func TestReduce_DifferentTypes(t *testing.T) {
	input := []string{"1", "2", "3", "4"}
	expected := 10 // Sum of integers converted from strings

	result := Reduce(input, func(acc int, val string) int {
		num, _ := strconv.Atoi(val)
		return acc + num
	}, 0)

	assert.Equal(t, expected, result)
}

func TestReduce_NilFunction(t *testing.T) {
	input := []int{1, 2, 3}
	var nilFunc func(int, int) int

	// This should panic since fn is nil
	assert.Panics(t, func() {
		Reduce(input, nilFunc, 0)
	})
}

func TestReduce_BooleanLogic(t *testing.T) {
	input := []bool{true, false, true, true}

	// Test AND operation
	resultAnd := Reduce(input, func(acc, val bool) bool {
		return acc && val
	}, true)
	assert.False(t, resultAnd)

	// Test OR operation
	resultOr := Reduce(input, func(acc, val bool) bool {
		return acc || val
	}, false)
	assert.True(t, resultOr)
}

func TestMap_LargeSlice(t *testing.T) {
	// Test with a larger slice to ensure performance is reasonable
	size := 1000
	input := make([]int, size)
	for i := 0; i < size; i++ {
		input[i] = i
	}

	result := Map(input, func(i int) int {
		return i * 2
	})

	assert.Len(t, result, size)
	assert.Equal(t, 0, result[0])
	assert.Equal(t, 2, result[1])
	assert.Equal(t, 1998, result[999])
}

func TestReduce_LargeSlice(t *testing.T) {
	// Test with a larger slice to ensure performance is reasonable
	size := 1000
	input := make([]int, size)
	for i := 0; i < size; i++ {
		input[i] = 1
	}

	result := Reduce(input, func(acc, val int) int {
		return acc + val
	}, 0)

	assert.Equal(t, size, result)
}

func TestMap_Chaining(t *testing.T) {
	// Test chaining Map operations
	input := []int{1, 2, 3, 4, 5}

	// First map: multiply by 2
	doubled := Map(input, func(i int) int {
		return i * 2
	})

	// Second map: convert to string
	strings := Map(doubled, func(i int) string {
		return strconv.Itoa(i)
	})

	expected := []string{"2", "4", "6", "8", "10"}
	assert.Equal(t, expected, strings)
}

func TestMapAndReduce_Integration(t *testing.T) {
	// Test Map and Reduce working together
	input := []int{1, 2, 3, 4, 5}

	// Map: square each number
	squared := Map(input, func(i int) int {
		return i * i
	})

	// Reduce: sum all squared numbers
	sum := Reduce(squared, func(acc, val int) int {
		return acc + val
	}, 0)

	expected := 55 // 1 + 4 + 9 + 16 + 25
	assert.Equal(t, expected, sum)
}

func TestMap_ZeroValues(t *testing.T) {
	input := []int{0, 0, 0}

	result := Map(input, func(i int) string {
		return strconv.Itoa(i)
	})

	expected := []string{"0", "0", "0"}
	assert.Equal(t, expected, result)
}

func TestReduce_ZeroValues(t *testing.T) {
	input := []int{0, 0, 0}

	result := Reduce(input, func(acc, val int) int {
		return acc + val
	}, 5)

	assert.Equal(t, 5, result) // Should return initial value
}

func TestMap_PointerTypes(t *testing.T) {
	input := []*int{intPtr(1), intPtr(2), intPtr(3)}

	result := Map(input, func(p *int) int {
		if p != nil {
			return *p * 2
		}
		return 0
	})

	expected := []int{2, 4, 6}
	assert.Equal(t, expected, result)
}

func TestReduce_PointerTypes(t *testing.T) {
	input := []*int{intPtr(1), intPtr(2), intPtr(3)}

	result := Reduce(input, func(acc int, p *int) int {
		if p != nil {
			return acc + *p
		}
		return acc
	}, 0)

	assert.Equal(t, 6, result)
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}
