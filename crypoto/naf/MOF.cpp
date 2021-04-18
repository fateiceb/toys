#include <iostream> 
#include <math.h>
using namespace std;
int main(){
	int z[40],m[40],n[40],v[40];
	int num=0,k,a=0,c=1,j,q,p=1;
	cout<<"Please enter a number: ";
	cin>>k;
	do{
		z[num++]=k%2;
		k=k/2;
	}while(k!=0);
	cout<<"The Binary Code is: ";               //二进制编码形式 
	for(int i=num-1;i>=0;i--)
	   cout<<z[i]<<" ";
	cout<<endl;
	for(int i=num-1;i>=0;i--)
	   m[a++]=z[i];
	n[0]=0;
	m[a]=0;
	cout<<"The MOF Code is: ";                 
	for(int i=0;i<a;i++)
	   n[c++]=m[i];
	for(int i=0;i<c;i++){
	  v[i]=m[i]-n[i];                       
	  cout<<v[i]<<" ";
    }
	  cout<<endl;                          //MOF(彼此相反型编码) 形式，相邻非零位符号相反 
	for(int i=0;i<c;i++){
		if (v[i]>=0)
		  v[i]=v[i];
		else 
		  v[i]=1;	
	}
	for(int i=c;i<32;i++) 
	  v[i]=2;
	if(v[1]==0){
	    q=p;
	    j=2;
    }
    else
    {
      	q=2*p;
	    j=2;
	}
	for (int j=2;j<c-1;j++){
	  if (v[j]==v[j+1]){
	  	q=2*q;
	  } 
	  else{
	  	if(v[j]!=v[j+1] && v[j]==1 && v[j+2]==1){
	  		q=4*q;
	  		q=q-p;
	  		j=j+1;
		  }
	  	else if (v[j]!=v[j+1] && v[j]==1 && v[j+2]==0){
	  		q=2*q;
	  		q=q-p;
	  		q=2*q;
	  		j=j+1;
		  }
		else if (v[j]!=v[j+1] && v[j]==0 && v[j+2]==1){
	  		q=2*q;
	  		q=q+p;
	  		q=2*q;
	  		j=j+1;
		  }
		else if (v[j]!=v[j+1] && v[j]==0 && v[j+2]==0)
		{
	  		q=2*q;
	  		q=2*q;
	  		q=q+p;
	  		j=j+1;
		  }	
	  }		
	}
	if( v[c-1]==0) 
      q=2*q;
    else
    {
    	q=2*q;
    	q=q+p;
	}
	cout<<"The scalar number is: ";
	cout<<q<<endl;    
	return 0;	
}
