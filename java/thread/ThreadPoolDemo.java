package thread;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class ThreadPoolDemo {
    public static void main(String[] args){
        //线程数固定的线程池
        ExecutorService es1 = Executors.newFixedThreadPool(2);
        //线程数根据任务动态调整的线程池
        ExecutorService es2 = Executors.newCachedThreadPool();
        //仅仅单个线程的线程池
        ExecutorService es3 = Executors.newSingleThreadExecutor();
        ScheduledExecutorService es4 = Executors.newScheduledThreadPool(4);

        for (int i = 0; i <6;i++ ){
            es1.submit(new Task(""+i));
            es2.submit(new Task(""+i));
            es3.submit(new Task(""+i));
            es4.scheduleAtFixedRate(new Task(""+i), 2, 3, TimeUnit.SECONDS);
        }
    }
}
class Task implements Runnable {
    private final String name;
    public Task(String name) {
        this.name = name;
    }
    @Override
    public void run() {
        // TODO Auto-generated method stub
        System.out.println("start task "+name);
        try {
            Thread.sleep(1000);
        } catch (InterruptedException e) {
        }
        System.out.println("end task"+name);
    }
    
}
