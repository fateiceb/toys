#include <iostream>
using namespace std;
int main(){
    cout <<"多位生成"<<endl;
    int n = 0,p =1,q=1;
    int k[32] = {0};
    int i = 0;
    cout << "please input number:"<<endl;
    cin >> n;
    while(n > 0){
        int t = n &0x3;     //取最后两位，有可能的值0,1,2,3
        switch (t) {
            case 0: 
                k[i] = 0;       //k[i],k[i+1]一次生成两位
                k[i+1] = 0;     
                i=i+2;          //i后移两位
                n = n >> 2; //  n/4
                break;
            case 1: 
                k[i] = 1;
                k[i+1] = 0;
                i=i+2;
                n = n >> 2; //  n/4
                break;
            case 2: 
                k[i] = 0;
                i=i+1;
                n = n >> 1;      // n/2,n除以2
                break;
            case 3: 
                k[i] = -1;
                k[i+1] = 0;
                i=i+2;
                n = n+1;
                n = n >> 2;   //    n/4
                break;
        }
    }
    cout<<"The NAF Code is: ";
  for(int j=i-2;j>=0;j--)
    cout<<k[j]<<" ";
  cout<<endl;
  
  for(int j=i-2;j>0;j--)     //NAF标量乘算法的实现 
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