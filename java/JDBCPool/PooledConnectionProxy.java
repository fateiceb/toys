package JDBCPool;

import java.sql.Connection;
import java.sql.SQLException;
import java.util.Queue;

public class PooledConnectionProxy extends AbstractConnectionProxy {

    Connection target;

    Queue<PooledConnectionProxy> idQueue; 

    public PooledConnectionProxy(Queue<PooledConnectionProxy> idQueue, Connection target){
        this.idQueue = idQueue;
        this.target = target;
    }

    public void close() throws SQLException{
        System.out.println("Fake close and released to idle queue for future reuse: " + target);
        idQueue.offer(this);
    }

    @Override
    protected Connection getRealConnection() {
        return target;
    }
    
}
