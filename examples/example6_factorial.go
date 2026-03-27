// Example 6: Complex example - factorial
var n int = 5;
var result int = 1;
var i int = 1;
for i <= n {
	result = result * i;
	i = i + 1;
}
fmt.Println("Factorial of 5 is:");
fmt.Println(result);
