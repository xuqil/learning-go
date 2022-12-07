package animal_test

import (
	animal2 "animal"
	"animal/dog"
	"testing"
)

func TestEat(t *testing.T) {
	testCases := []struct {
		name   string
		animal animal2.Animal
	}{
		{
			name:   "dog",
			animal: dog.Dog{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.animal.Eat()
			tc.animal.Sleep()
		})
	}
}
