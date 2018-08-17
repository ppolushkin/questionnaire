package main

func main() {
	a := App{}
	// You need to set your Username and Password here
	a.Initialize("quest", "quest", "localhost", "quest")

	a.Run(":8080")
}