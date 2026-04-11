type Point struct {
	x int;
	y int;
}

type Person struct {
	name string;
	age  int;
}

var p Point;
p.x = 10;
p.y = 20;
fmt.Println(p.x, p.y);

var user Person;
user.name = "Alice";
user.age = 30;
fmt.Println(user.name, user.age);
