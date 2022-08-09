package thread;

import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public class ReentrantLockDemo {
    private final Lock lock = new ReentrantLock();
    private int count;
    public void add(int n) {
        lock.lock();
        try {
            count += n;
        }finally{
            lock.unlock();
        }
    }
    public void addTry(int n) {
        try {
            if (lock.tryLock(1,TimeUnit.SECONDS)) {
                try {

                }finally {
                    lock.unlock();
                }
            }
        } catch (InterruptedException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
    }
}
