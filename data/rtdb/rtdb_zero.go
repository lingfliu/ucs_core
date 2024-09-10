package rtdb

/*
int add(int a, int b) {
	return a + b;
}
*/
import "C"

func Add(a int, b int) int {
	return C.add(a, b)
}
