package thread_print;

import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public class P100 {
   Thread t1 = new Thread(new Counter(0));
   Thread t2 = new Thread(new Counter(1));
   Thread t3 = new Thread(new Counter(2));
   t1.start();
   t2.start();
   t3.start();
}
class Counter implements Runnable {
    private static final ReentrantLock LOCK = new ReentrantLock();
    private static final Condition c = LOCK.newCondition(); 
    private static int cnt = 0;
    private Integer threadNumber;
    public Counter(Integer threadNumber){
        this.threadNumber = threadNumber;
    }
    public void run(){
        while(true) {
            try{
                LOCK.lock();
                while(cnt % 3 != this.threadNumber) {
                    if (cnt >= 101){
                        break;
                    }
                    try {
                        c.await();
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                        //TODO: handle exception
                    }
                }
            if (cnt >= 101) {
                break;
            }
            System.out.println("thread-"+this.threadNumber+":"+cnt);
            cnt++;
            c.signalAll();
            }finally {
                LOCK.unlock();
            }
        }
   }
}

