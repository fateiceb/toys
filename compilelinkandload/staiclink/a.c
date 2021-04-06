//gcc -c a.c b.c 得到a.o,b.o
/*
链接 step1：空间与地址分配
段合并
收集符号定义及应用，放入全局符号表
step2：符号解析与重定位
符号解析，重定位、调整代码中的地址。重定位和链接过程的核心。
*/
//ld a.o b.o -e main -o ab 
//但是ld出错了，还是看看书上的吧。
// -e main表示main函数是程序入口
//-o ab 输出文件名为ab，默认是a.out
//objdump -h a.o   VMA表示虚拟地址，LMA表示加载地址，正常情况是相同的
extern int shared;
// void swap(int *a,int *b);
int main()
{
    int a = 100;
    swap( &a, &shared);
}