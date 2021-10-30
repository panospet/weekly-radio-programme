package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStartEnd(t *testing.T) {
	input := "15:00-21:00"
	start, end, err := GetStartEnd(input)
	assert.Nil(t, err)
	assert.Equal(t, "15:00", start.Format("15:04"))
	assert.Equal(t, "21:00", end.Format("15:04"))
}

func TestHasConflict(t *testing.T) {
	input1 := "15:00-21:00"
	input2 := "22:00-00:00"
	has, err := HasConflict(input1, input2)
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestHasConflict2(t *testing.T) {
	input1 := "15:00-21:00"
	input2 := "20:00-00:00"
	has, err := HasConflict(input1, input2)
	assert.Nil(t, err)
	assert.True(t, has)
}


func TestHasConflict3(t *testing.T) {
	input1 := "15:00-21:00"
	input2 := "21:00-00:00"
	has, err := HasConflict(input1, input2)
	assert.Nil(t, err)
	assert.False(t, has)
}
func TestHasConflict4(t *testing.T) {
	input1 := "21:00-22:00"
	input2 := "21:59-00:00"
	has, err := HasConflict(input1, input2)
	assert.Nil(t, err)
	assert.True(t, has)
}

func TestHasConflict5(t *testing.T) {
	input1 := "00:00-02:00"
	input2 := "22:00-00:00"
	has, err := HasConflict(input1, input2)
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestHasConflict6(t *testing.T) {
	input1 := "00:00-02:00"
	input2 := "01:00-20:00"
	has, err := HasConflict(input1, input2)
	assert.Nil(t, err)
	assert.True(t, has)
}
