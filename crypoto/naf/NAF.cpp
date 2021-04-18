#include <iostream> 
using namespace std;
int main()
{    
 int n,p=1,q=1,z;  
 int a[32]={0};    
 int i=0,j,k; 
 cout<<"Please enter a number: ";
 cin>>n;           
  while(n>0)               //NAF(二元非相邻形式) 非零元素不得相邻 
  {  
   if(n%2==0)
   {
    a[i++]=n%2;         
       n=n/2;   
   }
   else 
   {
    if((((n-1)/2)%2)==0)    //所得商为偶数 
    {
       a[i++]=1;         
        n=n/2;   
    }
    else                    //if((((n+1)/2)%2)==0)
    {
        a[i++]=-1;         
        n=n/2+1;   
    }
   }      
 }
 cout<<"The NAF Code is: ";         
  for(j=i-1;j>=0;j--)        
   cout<<a[j]<<" ";    
 cout<<endl; 
 
 for(k=i-1;k>0;k--)          //NAF标量乘算法的实现 
  {
   q=2*q;
   if(a[k-1]>0)
    q=q+p;
   else if(a[k-1]<0)
    q=q-p;
  }
 cout<<"The scalar number is: ";
 cout<<q<<endl;
  return 0;
}
