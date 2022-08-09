package thread;

import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

public class FutureDemo {
    public static void main(String[] args) {
        ExecutorService executor = Executors.newFixedThreadPool(2);
        Callable<String> task = new TaskFuture();
        Future<String> future = executor.submit(task);
        String result = "";
        try {
            result = future.get();
        } catch (InterruptedException | ExecutionException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
        System.out.println(result);
    }
}
class TaskFuture implements Callable<String> {
    public String call() throws Exception {
        return "2";
    }
}
