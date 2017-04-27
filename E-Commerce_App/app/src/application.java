import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

import java.io.*;
import java.net.*;
import java.lang.Object;


public class application {




    public static void main(String[] args) throws Exception{
        int socket = 1339;



        System.out.println("Application listening on socket: " + socket);
        //create an HttpServer accepting connections on localhost
        HttpServer appServer = HttpServer.create(new InetSocketAddress(socket), 0);


        appServer.createContext("/app/test", new RequestHandler.TestHandler());
        appServer.createContext("/app/browse", new RequestHandler.BrowseHandler());
        appServer.createContext("/app/createCart", new RequestHandler.CartCreate());
        appServer.createContext("/app/addToCart", new RequestHandler.CartAdd());
        appServer.createContext("/app/deleteFromCart", new RequestHandler.CartDelete());
        appServer.createContext("/app/checkout", new RequestHandler.Checkout());
        appServer.createContext("/app/orderStatus", new RequestHandler.OrderStatus());

        appServer.start();
    }



}
