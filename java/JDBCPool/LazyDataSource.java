package JDBCPool;

import java.io.PrintWriter;
import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.SQLFeatureNotSupportedException;
import java.util.logging.Logger;

import javax.management.RuntimeErrorException;
import javax.sql.DataSource;

public class LazyDataSource implements DataSource {

    private String url;
    private String username;
    private String password;

    public LazyDataSource(String url, String username, String password){
        this.url = url;
        this.password = password;
        this.username = username;
    }

    public Connection getConnection(String username, String password) throws SQLException{
        return new LazyConnectionProxy(() -> {
           try {
               Connection connection = DriverManager.getConnection(url, username, password);
               System.out.println("open connection: " + connection);
               return connection;
           }catch (SQLException e){
                throw new RuntimeException(e);
           } 
        });
    }

    @Override
    public Logger getParentLogger() throws SQLFeatureNotSupportedException {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public <T> T unwrap(Class<T> iface) throws SQLException {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public boolean isWrapperFor(Class<?> iface) throws SQLException {
        // TODO Auto-generated method stub
        return false;
    }

    @Override
    public Connection getConnection() throws SQLException {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public PrintWriter getLogWriter() throws SQLException {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public void setLogWriter(PrintWriter out) throws SQLException {
        // TODO Auto-generated method stub
        
    }

    @Override
    public void setLoginTimeout(int seconds) throws SQLException {
        // TODO Auto-generated method stub
        
    }

    @Override
    public int getLoginTimeout() throws SQLException {
        // TODO Auto-generated method stub
        return 0;
    }
    
}
