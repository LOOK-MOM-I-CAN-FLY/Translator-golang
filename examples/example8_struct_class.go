struct Point {
	x int;
	y int;
};

class Person {
	var name string;
	var age int;
};

var p Point;
p.x = 10;
p.y = 20;
fmt.Println(p.x, p.y);

var user Person;
user.name = "Alice";
user.age = 30;
fmt.Println(user.name, user.age);

