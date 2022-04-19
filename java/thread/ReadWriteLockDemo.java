package thread;

import java.util.Arrays;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;

public class ReadWriteLockDemo {
    private final Lock lock = new ReentrantLock();
    private int[] counts = new int[10];
    public void inc(int index) {
        lock.lock();
        try {
            counts[index] += 1;
        } finally {
            lock.unlock();
        }
    }
    public int[] get() {
        lock.lock();
        try {
            return Arrays.copyOf(counts, counts.length);
        }finally {
            lock.unlock();
        }
    }
}
/*
    readwriteLock
*/
private class Counter {
    private final ReadWriteLock rwlock = new ReentrantReadWriteLock();
    private final Lock rLock  = rwlock.readLock();
    private final Lock wlock = rwlock.writeLock();
    private int[] counts = new int[10];
    public void inc(int index) {
        wlock.lock();
        try {
            counts[index] += 1;
        }finally{
            wlock.unlock();
        }
    }
    public int[] get(){
        rLock.lock();
        try {
            return Arrays.copyOf(counts, counts.length);
        }
    }
}
