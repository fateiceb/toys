package thread_pool;

import java.util.concurrent.TimeUnit;

public class Task implements Runnable {
    private String name;
    public Task(String name) {
        this.name = name;
    }
    public String getName(){
        return this.name;
    }
    public void run() {
        try {
            Long duration = (long)(Math.random() * 10);
            System.out.println("Excutiong:" + name);
            TimeUnit.SECONDS.sleep(duration);
        } catch (InterruptedException e) {
            //TODO: handle exception
            e.printStackTrace();
        }
    }
}
