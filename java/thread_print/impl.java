package thread_print;

public interface impl {
    default void hello(){
        System.out.printf("hello, ");
    }

    void who(String name);
}
