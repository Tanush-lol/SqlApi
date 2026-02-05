#include<iostream>

int main(){
  int a,b=0;
  std::cout<<"enter the number of columns in the array"<<'\n';
  std::cin>>a;
  std::cout<<"enter the number of rows"<<'\n';
  std::cin>>b;

  int arr[a][b];
  
  for(int i=0;i<a;i++){
    for(int j=0;j<b;j++){
      arr[i][j] = i + j;
      std::cout<<arr[i][j]<<" ";
    }
    std::cout<<'\n';
  }
}
//
//but i still don't understand stack memory is stiff and like stack memory arrays we can't change their size in runtime, right so in this code the compiler doesn't know the size of array initially rather makes the array during runtime right
