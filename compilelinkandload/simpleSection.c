//gcc -c simpleSection.c
//objdump -h simpleSection.o
//gcc bintuils 可以用于windwos和linux linux下的ELF文件 windows下的pe文件 readelf是linux下的工具解析elf文件
//size 可以查看ELF文件的代码段、数据段和BSS段的长度
//objdump -s -d simpleSection.o -s 16进制表示 -d 指令反汇编
//readelf -h simpleSection.o 详细查看ELF文件头
//elf定义 /usr/include/elf.h
//nm simpleSection.c 查看符号
/*
无论是可执行文件、目标文件或库，实际上都是一样基于段的文件或这种文件的集合。程序的源代码经过编译以后，按照代码和数据分别存放
到相应的断种，编译器还会将一些辅助性的信息，诸如符号、重定位信息等也按照表的方式存放到目标文件中，而通常情况下，一个表就是一个段
*/
int printf(const char * format, ...);
int global_init_var = 84;
int global_uninit_var;
void func1 (int i)
{
    printf("%d\n",i);
}
int main(void)
{
    static int static_var = 85;
    static int static_var2;
    int a = 1;
    int b;
    func1(static_var + static_var2 +a +b);
    return a;
}