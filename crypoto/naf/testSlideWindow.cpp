#include<iostream>
#include<math.h>
using namespace std;

int count_zero(int t)
{
    int counter = 0;
    int mid = 0;
    while(t > 0) {
        mid = t%10;
        t = t /10;
        if(mid == 0){
            counter++;
        }
    }
    return  counter;
}

int main(){
    int n,w,t,mid,i = 0;
    int k[50] = {0};
    cout<<"please enter num"<<endl;
    cin>>n;
    cout<<"please enter w"<< endl;
    cin>>w;
    mid = pow(2,w) - 1; //2^w - 1
    while( n >= 1){
        //n为奇数
        if( n %2 != 0){
            t = n & mid; //1 k & 2^w - 1 
            // cout << t; //debug
            //2
            if (t > pow(2,w-1)){
                k[i] = t - pow(2,w);
                n = n - k[i];
            }
            else{
                k[i] = t;
               
            }
            //debug
            //3
            for (int j = 1; j<=w-1;j++) {
                k[i+j] = 0;
            }
            //4
            n = n >> w;
            //5
            i = i+2;
        }else {
           int count = count_zero(t);
        //    if(count == 0)
        //    {
        //       count = 1;
        //    }
           cout <<"counter:"<<count<<endl;
           for(int f = 0 ;f <count;f++){
               k[i+f] = 0;
           }
           i = i + count;
           n = n >> count;
        }


    }
     for(int j = 49;j>=0;j--) {
        cout<<" "<<k[i]<<" ";
     }

    // int count = count_zero(n);
    // cout<<"counter"<<count<<endl; 
    // mid = pow(2,w) -1; //2^w - 1
    // t = k & mid;  // k & 2^w - 1 
    // cout << w<<endl;
    // cout <<k<<endl;
    // cout << pow(k,w-1);
   
    
    return 1;
}