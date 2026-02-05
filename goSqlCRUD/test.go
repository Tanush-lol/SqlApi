package main

import(
	"fmt"
)

func main(){
	arr:=[5]int {1,2,3,4,5}
	s:=arr[1:4]
	fmt.Println(s)
	s[0] = 32
	fmt.Println(s)
	fmt.Println(arr)

	naws := make([]int,3,5)
	fmt.Println(naws)
	fmt.Println(len(naws));
	fmt.Println("This the cap length",cap(naws));

	for i:=0; i<10;i++{
		naws = append(naws,i)
	}
	fmt.Println(naws)
	fmt.Println(len(naws));
	fmt.Println("This the cap length",cap(naws));

	//2d experimentation cuh i forgot how 2d arrays work in cpp 
	twoD := make([][]int,0)
	fmt.Println(twoD)
	
	for i:=0; i<10;i++{
//		twoD = append(twoD,i) not possible as we can only insert slices into 2d slices
		row:= make([]int,0)
		for j:=0;j<10;j++{
			row= append(row,j)
		}
		twoD = append(twoD,row)
	}

for _, row := range twoD {
    for _, val := range row {
        fmt.Printf("%2d ", val)
    }
    fmt.Println()
}

}








