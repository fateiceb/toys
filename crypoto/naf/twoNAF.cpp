#include <iostream>
using namespace std;
int main(){
    cout <<"多位生成"<<endl;
    int n = 0,p =1,q=1;
    int k[32] = {0};
    int i = 0;
    cout << "please input n:"<<endl;
    cin >> n;
    while(n > 0){
        int t = n &0x3;
        switch (t) {
            case 0: 
                k[i] = 0;
                k[i+1] = 0;
                i=i+2;
                n = n >> 2;
                break;
            case 1: 
                k[i] = 1;
                k[i+1] = 0;
                i=i+2;
                n = n >> 2;
                break;
            case 2: 
                k[i] = 0;
                i=i+1;
                n = n >> 1;
                break;
            case 3: 
                k[i] = -1;
                k[i+1] = 0;
                i=i+2;
                n = n+1;
                n = n >> 2;
                break;
        }
    }
    cout<<"The NAF Code is: ";
  for(int j=i-2;j>=0;j--)
    cout<<k[j]<<" ";
  cout<<endl;
  
  for(int j=i-2;j>0;j--)     
  {
   q=2*q;
   if(k[j-1]>0)
    q=q+p;
   else if(k[j-1]<0)
    q=q-p;
  }
 cout<<"The scalar number is: ";
 cout<<q<<endl;
  return 1;
}