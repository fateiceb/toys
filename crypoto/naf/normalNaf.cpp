#include <iostream>
#include <math.h>
using namespace std;

int main(){
    int i = 0 ,n = 0,p =1,q=1;;
    int k[32] = {0};
    cout << "please enter a num"<<endl;
    cin >> n;
    while(n >= 1){
        // ���n������
        if (n %2 != 0){
            k[i] = 2 -(n%4);
            n = n - k[i];
        }else{      //ż��
            k[i] = 0;
        }
        n = n /2;
        i = i +1;
    }

    for(int j = i-1;j>=0;j--){      // ���naf code
        cout << " "<<k[j];
    }
      cout <<endl;
  for(int j=i-1;j>0;j--)     //NAF�������㷨��ʵ��
  {
   q=2*q;
   if(k[j-1]>0)
    q=q+p;
   else if(k[j-1]<0)
    q=q-p;
  }
 cout<<"The scalar number is: ";
 cout<<q<<endl;
    cout<<endl;
}