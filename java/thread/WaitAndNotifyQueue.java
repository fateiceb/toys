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
    }
    public synchronized String getTask() {
        while(queue.isEmpty()){
            this.wait();
        }
        return queue.remove();
    }
}