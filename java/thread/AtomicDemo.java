package thread;

import java.util.concurrent.atomic.AtomicInteger;

/**
 * AtomicDemo
 */
public class AtomicDemo {
    /*
        使用CAS实现原子类
    */
    public int incrementAndGet(AtomicInteger var){
        int prev,next;
        do {
            prev = var.get();
            next = prev +1 ;
        }while(!var.compareAndSet(prev, next));
        return next;
    }
}