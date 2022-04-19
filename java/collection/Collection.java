package collection;

import java.util.ArrayList;
import java.util.Hashtable;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Collection
 */
public class Collection {
    public static void main(String[] args) {
        Map<String,String> m = new Hashtable<String,String>();
        m.put("a", "a");
        System.out.println(m);
        ConcurrentHashMap cm = new ConcurrentHashMap<String,String>();
    }
}