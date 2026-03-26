package backend

// Formula represents a mathematical formula with a name, description, and the formula template.
type Formula struct {
	Name        string
	Description string
	Template    string
}

// CommonFormulas is a list of commonly used mathematical formulas.
var CommonFormulas = []Formula{
	{
		Name:        "Circle Area",
		Description: "Area of a circle given radius r",
		Template:    "pi * r^2",
	},
	{
		Name:        "Circle Circumference",
		Description: "Circumference of a circle given radius r",
		Template:    "2 * pi * r",
	},
	{
		Name:        "Quadratic Formula",
		Description: "Roots of ax^2 + bx + c = 0",
		Template:    "(-b + sqrt(b^2 - 4*a*c)) / (2*a)",
	},
	{
		Name:        "Pythagorean Theorem",
		Description: "Hypotenuse c of a right triangle given legs a and b",
		Template:    "sqrt(a^2 + b^2)",
	},
	{
		Name:        "Sphere Volume",
		Description: "Volume of a sphere given radius r",
		Template:    "(4/3) * pi * r^3",
	},
	{
		Name:        "Triangle Area",
		Description: "Area of a triangle given base b and height h",
		Template:    "0.5 * b * h",
	},
	{
		Name:        "Cylinder Volume",
		Description: "Volume of a cylinder given radius r and height h",
		Template:    "pi * r^2 * h",
	},
	{
		Name:        "Logarithm (Base 10)",
		Description: "Logarithm of x with base 10",
		Template:    "log10(x)",
	},
	{
		Name:        "Natural Logarithm",
		Description: "Logarithm of x with base e",
		Template:    "ln(x)",
	},
	{
		Name:        "Compound Interest",
		Description: "Amount A after time t with principal P, rate r, and n times compounded",
		Template:    "P * (1 + r/n)^(n*t)",
	},
}
