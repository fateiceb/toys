#include <iostream>
#include <math.h>
using namespace std;

int main(){
    int i = 0 ,n = 0;
    int k[32] = {0};
    cin >> n;
    while(n >= 1){
        if (n %2 != 0){
            k[i] = 2 -(n%4);
            n = n - k[i];
        }else{
            k[i] = 0;
        }
        n = n /2;
        i = i +1;
    }
    for(int j = i-1;j>=0;j--){
        cout << " "<<k[j];
    }
    cout<<endl;
}