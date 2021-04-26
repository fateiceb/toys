#include<iostream>
#include <math.h>
using namespace std;
int main(){
    int n = 0,i = 0;
    int w = 0;
    int k[50] = {0};
    cout << "please enter a num "<<endl;
    cin >> n;
    cout << "please enter a window num"<<endl;
    cin >> w;
     int mid = pow(2,w);
     cout <<"mid"<<mid<<endl;
    while(n >= 1){
        //若k是奇数
        if (n % 2 != 0){
            k[i] = n % mid;
            n = n - k[i];
            n = n / mid;
            for(int j = 1;j<=w-1;j++){
                k[i+j] = 0;
            }
            i = i + w;
        }else{
            k[i] = 0;
            n = n/2;
            i = i + 1;
        }
      
    }
    for (int j = i-1; j >=0;j--){
        cout <<" "<<k[j];
    }
    cout<<endl;

}