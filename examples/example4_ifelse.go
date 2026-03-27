// Example 4: If-else conditional
var age int = 20;
if age >= 18 {
	fmt.Println("Adult");
} else {
	fmt.Println("Minor");
}

var score int = 75;
if score >= 90 {
	fmt.Println("A");
} else {
	if score >= 80 {
		fmt.Println("B");
	} else {
		fmt.Println("C");
	}
}
