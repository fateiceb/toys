import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;

import javax.sql.DataSource;

import thread_print.*;
import JDBCPool.*;
/**  
 * main
 */
public class main {
    public static void main(String[] args) throws SQLException {
        String jdbcUrl = "jdbc:mysql://localhost:3306/test";
        String jdbcUsername = "root";
        String jdbcPassword = "root";
        
        DataSource pooledDataSource = new PooledDataSource(jdbcUrl, jdbcUsername, jdbcPassword);
        try (Connection conn = pooledDataSource.getConnection()) {
        }
        try (Connection conn = pooledDataSource.getConnection()) {
            // 获取到的是同一个Connection
        }
        try (Connection conn = pooledDataSource.getConnection()) {
            // 获取到的是同一个Connection
        }

        DataSource lazyDataSource = new LazyDataSource(jdbcUrl, jdbcUsername, jdbcPassword);
        System.out.println("get lazy connection...");
        try (Connection conn1 = lazyDataSource.getConnection()) {
            // 并没有实际打开真正的Connection
        }
        System.out.println("get lazy connection...");
        try (Connection conn2 = lazyDataSource.getConnection()) {
            try (PreparedStatement ps = conn2.prepareStatement("SELECT * FROM stu")) { // 打开了真正的Connection
                try (ResultSet rs = ps.executeQuery()) {
                    while (rs.next()) {
                        System.out.println(rs.getString("name"));
                    }
                }
            }
        }
    }
}

class Student{
    public String name;
    Student(String name) {
        this.name = name;
    }
    public String toString() {
        return "{Person: " + name + "}";
    }
    // @Override
    // public int compareTo(Student stu) {
    //     return this.name.compareTo(stu.name);
    // }
}

class sayHello implements impl{

    @Override
    public void who(String name) {
        System.out.println(name);
    }
}