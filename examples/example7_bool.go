// Example 7: Boolean operations and logical operators
var x int = 5;
var y int = 10;
var isGreater bool = x > y;
var isEqual bool = x == y;
var combined bool = isGreater || isEqual;
fmt.Println(combined);

if x < y && x > 0 {
	fmt.Println("x is between 0 and y");
}
