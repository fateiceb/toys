package thread;

import java.util.LinkedList;
import java.util.Queue;

/**
 * WaitAndNotifyQueue
 */
public class WaitAndNotifyQueue {
    Queue<String> queue = new LinkedList<>();
    public synchronized void addTask(String s) {
        this.queue.add(s);
        this.notify();
    }
    public synchronized String getTask() {
        while(queue.isEmpty()){
            try {
                this.wait();
            } catch (InterruptedException e) {
                // TODO Auto-generated catch block
                e.printStackTrace();
            }
        }
        return queue.remove();
    }
}