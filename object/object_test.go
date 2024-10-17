package object

import "testing"

func TestStringHashKey(t *testing.T) {
	firInput := &String{Value: "Hello World"}
	secInput := &String{Value: "Hello World"}

	thrdInput := &String{Value: "This is another test"}
	frthInput := &String{Value: "This is another test"}

	if firInput.HashKey() != secInput.HashKey() {
		t.Errorf("Strings with the same content have different Hash Keys!")
	}

	if thrdInput.HashKey() != frthInput.HashKey() {
		t.Errorf("Strings with the same content have different Hash Keys!")
	}

	if firInput.HashKey() == thrdInput.HashKey() {
		t.Errorf("Strings with different content have the same Hash Keys!")
	}

}

func TestIntegerHashKey(t *testing.T) {
	firInput := &Integer{Value: 1}
	secInput := &Integer{Value: 1}

	thrdInput := &Integer{Value: 3}
	frthInput := &Integer{Value: 3}

	if firInput.HashKey() != secInput.HashKey() {
		t.Errorf("Integers with the same values have different Hash Keys!")
	}

	if thrdInput.HashKey() != frthInput.HashKey() {
		t.Errorf("Integers with the same values have different Hash Keys!")
	}

	if firInput.HashKey() == thrdInput.HashKey() {
		t.Errorf("Integers with different values have the same Hash Keys!")
	}
}

func TestBooleanHashKey(t *testing.T) {
	firInput := &Boolean{Value: true}
	secInput := &Boolean{Value: true}

	thrdInput := &Boolean{Value: false}
	frthInput := &Boolean{Value: false}

	if firInput.HashKey() != secInput.HashKey() {
		t.Errorf("Booleans with the same values have different Hash Keys!")
	}

	if thrdInput.HashKey() != frthInput.HashKey() {
		t.Errorf("Booleans with the same values have different Hash Keys!")
	}

	if firInput.HashKey() == thrdInput.HashKey() {
		t.Errorf("Booleans with different values have the same Hash Keys!")
	}
}
