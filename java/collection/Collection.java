package collection;

import java.util.ArrayList;
import java.util.Hashtable;
import java.util.Map;
import java.util.TreeMap;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.AbstractQueuedSynchronizer;

/**
 * Collection
 */
public class Collection {
    public static void main(String[] args) {
        Map<String,String> m = new Hashtable<String,String>();
        m.put("a", "a");
        System.out.println(m);
        ConcurrentHashMap cm = new ConcurrentHashMap<String,String>();
        Map<String,String> m2 = new TreeMap<String,String>();
        m2.put(key, value)
        
    }
}