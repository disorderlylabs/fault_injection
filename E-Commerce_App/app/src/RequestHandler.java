import com.sun.net.httpserver.Headers;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;

import java.io.*;
import java.net.*;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;


public class RequestHandler {
    static final int SUCCESS = 200;
    static final int BAD_REQUEST = 400;
    static final int METHOD_NOT_ALLOWED = 405;
    static final int INTERNAL_ERR = 500;

    static microservice.Catalog catalog = new microservice.Catalog();
    static microservice.Cart cart = new microservice.Cart();
    static microservice.OrderManagement orderManagement = new microservice.OrderManagement();
    //static HttpClient client = HttpClientBuilder.create().build();
    static Map<String, Object> parameters = new HashMap<String, Object>();
    static String charset = "UTF-8";

    public static void parseQuery(String query, Map<String,
            Object> parameters) throws UnsupportedEncodingException {

        if (query != null) {
            String pairs[] = query.split("[&]");
            for (String pair : pairs) {
                String param[] = pair.split("[=]");
                String key = null;
                String value = null;
                if (param.length > 0) {
                    key = URLDecoder.decode(param[0],
                            System.getProperty("file.encoding"));
                }

                if (param.length > 1) {
                    value = URLDecoder.decode(param[1],
                            System.getProperty("file.encoding"));
                }

                if (parameters.containsKey(key)) {
                    Object obj = parameters.get(key);
                    if (obj instanceof List<?>) {
                        List<String> values = (List<String>) obj;
                        values.add(value);

                    } else if (obj instanceof String) {
                        List<String> values = new ArrayList<String>();
                        values.add((String) obj);
                        values.add(value);
                        parameters.put(key, values);
                    }
                } else {
                    parameters.put(key, value);
                }
            }
        }
    }


    public static void writeResponse(HttpExchange httpExchange, String response, int code) throws IOException {
        httpExchange.sendResponseHeaders(code, response.length());
        OutputStream os = httpExchange.getResponseBody();
        os.write(response.getBytes());
        os.close();
    }

    static void _setPostParameters(HttpURLConnection httpConnection, int queryLength) throws IOException {
        httpConnection.setDoOutput(true);
        httpConnection.setInstanceFollowRedirects(false);
        httpConnection.setRequestMethod("POST");
        httpConnection.setRequestProperty("charset", "utf-8");
        httpConnection.setRequestProperty(
                "Content-Type", "application/x-www-form-urlencoded");
        httpConnection.setRequestProperty("Content-Length", Integer.toString(queryLength));
        httpConnection.setUseCaches(false);
    }


    static class BrowseHandler implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            String response = "this should not be displayed";
            OutputStream os = t.getResponseBody();
            int returnCode = INTERNAL_ERR;

            if (t.getRequestMethod().equalsIgnoreCase("GET")) {
                URL url = new URL(catalog.items());
                HttpURLConnection httpConnection= (HttpURLConnection)url.openConnection();
                httpConnection.setRequestProperty( "charset", "utf-8");

                int responseCode = httpConnection.getResponseCode();
                if(responseCode != SUCCESS) {
                    returnCode = responseCode;
                    response = "Error reading response from catalog:items\n";
                }else {
                    BufferedReader rd = new BufferedReader( new InputStreamReader(httpConnection.getInputStream()));
                    response = rd.readLine();

                    if(response == null) {
                        response = "Error reading response from catalog:get\n";
                        returnCode = BAD_REQUEST;
                    }else{
                        returnCode = SUCCESS;
                    }

                }
            }else{
                response = "Only GET requests\n";
                returnCode = METHOD_NOT_ALLOWED;
            }

            t.sendResponseHeaders(returnCode, response.length());
            os.write(response.getBytes());
            os.close();
        }
    }//end BrowseHandler



    static class CartCreate implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            String response = "this should not be displayed";
            OutputStream os = t.getResponseBody();
            StringBuffer req = new StringBuffer();
            int returnCode = INTERNAL_ERR;

            if (t.getRequestMethod().equalsIgnoreCase("POST")) {
                Headers headers = t.getRequestHeaders();

                String userID = headers.getFirst("userID");
                if (userID == null) {
                    response = "Could not parse userID\n";
                    returnCode = BAD_REQUEST;
                } else {
                    //create the post request
                    String request = cart.create();
                    URL url = new URL( request );
                    HttpURLConnection httpConnection= (HttpURLConnection) url.openConnection();
                    String query = String.format("userID=%s",userID);
                    _setPostParameters(httpConnection, query.length());

                    int responseCode = httpConnection.getResponseCode();
                    System.out.println("Response Code : " + responseCode);

                    if(responseCode != SUCCESS) {
                        response = "Error from cart.create()\n";
                        returnCode = responseCode;
                    }else {
                        BufferedReader br = new BufferedReader(new InputStreamReader(httpConnection.getInputStream()));
                        response = br.readLine();
                        if (response == null) {
                            response = "no response from cart.create()";
                            returnCode = BAD_REQUEST;
                        }else{
                            returnCode = SUCCESS;
                        }
                    }
                }
            } else {
                response = "Only POST requests\n";
                returnCode = METHOD_NOT_ALLOWED;
            }

            t.sendResponseHeaders(405, response.length());
            os.write(response.getBytes());
            os.close();
        }
    }//end cartCreate


    static class CartAdd implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            String response = "this should not be displayed";
            OutputStream os = t.getResponseBody();
            StringBuffer req = new StringBuffer();
            int returnCode = INTERNAL_ERR;

            if (t.getRequestMethod().equalsIgnoreCase("POST")) {
                Headers headers = t.getRequestHeaders();

                String itemID = headers.getFirst("itemID");
                if (itemID == null) {
                    response = "Could not parse itemID\n";
                    t.sendResponseHeaders(BAD_REQUEST, response.length());
                    os.write(response.getBytes());
                    os.close();
                    return;
                }

                String cartID = headers.getFirst("cartID");
                if (cartID == null) {
                    response = "Could not parse cartID\n";
                    returnCode = BAD_REQUEST;
                } else {
                    String request = cart.addItem();
                    URL url = new URL( request );
                    HttpURLConnection httpConnection= (HttpURLConnection) url.openConnection();
                    String query = String.format("itemID=%s&cartID=%s",itemID, cartID);
                    _setPostParameters(httpConnection, query.length());

                    int responseCode = httpConnection.getResponseCode();
                    System.out.println("Response Code : " + responseCode);

                    if(responseCode != SUCCESS) {
                        response = "Error from cart.add()\n";
                        returnCode = responseCode;
                    }else {
                        BufferedReader br = new BufferedReader(new InputStreamReader(httpConnection.getInputStream()));
                        response = br.readLine();
                        if (response == null) {
                            response = "no response from cart.add()";
                            returnCode = BAD_REQUEST;
                        }else{
                            returnCode = SUCCESS;
                        }
                    }
                }
            } else {
                response = "Only POST requests\n";
                returnCode = METHOD_NOT_ALLOWED;
            }

            t.sendResponseHeaders(returnCode, response.length());
            os.write(response.getBytes());
            os.close();
        }
    }//end cartAdd


    static class CartDelete implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            String response = "this should not be displayed";
            OutputStream os = t.getResponseBody();
            StringBuffer req = new StringBuffer();
            int returnCode = INTERNAL_ERR;

            if (t.getRequestMethod().equalsIgnoreCase("POST")) {
                Headers headers = t.getRequestHeaders();

                String itemID = headers.getFirst("itemID");
                if (itemID == null) {
                    response = "Could not parse itemID\n";
                    t.sendResponseHeaders(BAD_REQUEST, response.length());
                    os.write(response.getBytes());
                    os.close();
                    return;
                }

                String cartID = headers.getFirst("cartID");
                if (cartID == null) {
                    response = "Could not parse cartID\n";
                    returnCode = BAD_REQUEST;
                } else {
                    String request = cart.deleteItem();
                    URL url = new URL( request );
                    HttpURLConnection httpConnection= (HttpURLConnection) url.openConnection();
                    String query = String.format("itemID=%s&cartID=%s",itemID, cartID);
                    _setPostParameters(httpConnection, query.length());

                    int responseCode = httpConnection.getResponseCode();
                    System.out.println("Response Code : " + responseCode);

                    if(responseCode != SUCCESS) {
                        response = "Error from call to cart.deleteItem()\n";
                        returnCode = responseCode;
                    }else {
                        BufferedReader br = new BufferedReader(new InputStreamReader(httpConnection.getInputStream()));
                        response = br.readLine();
                        if (response == null) {
                            response = "no response from cart.deleteItem()";
                            returnCode = BAD_REQUEST;
                        }else{
                            returnCode = SUCCESS;
                        }
                    }
                }
            } else {
                response = "Only POST requests\n";
                returnCode = METHOD_NOT_ALLOWED;
            }

            t.sendResponseHeaders(returnCode, response.length());
            os.write(response.getBytes());
            os.close();
        }
    }


    static class Checkout implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            System.out.println("in checkout");
            String response = "this should not be displayed";
            OutputStream os = t.getResponseBody();
            StringBuffer req = new StringBuffer();
            BufferedReader br;
            URLConnection connection;
            int responseCode = 200;
            parameters.clear();


            if (t.getRequestMethod().equalsIgnoreCase("POST")) {
                InputStreamReader isr = new InputStreamReader(t.getRequestBody(), "utf-8");
                br = new BufferedReader(isr);
                String query = br.readLine();
                parseQuery(query, parameters);


                String cartID = parameters.get("cartID").toString();
                String userID = parameters.get("userID").toString();

                if (cartID == null) {
                    response = "Could not parse cartID\n";
                    //t.sendResponseHeaders(400, response.length());
                    responseCode = 400;
                    writeResponse(t, response, responseCode);
                } else if (userID == null) {
                    response = "Could not parse userID\n";
                    responseCode = 400;
                    writeResponse(t, response, responseCode);
                } else {
                    //1)get items from the cart
                    //connection.set
                    connection = new URL(cart.items() + "?cartID=" + cartID).openConnection();
                    connection.setRequestProperty("Accept-Charset", "UTF-8");

                    br = new BufferedReader(new InputStreamReader(connection.getInputStream()));
                    response = br.readLine();
                    if (response == null) {
                        response = "Error reading response from cart:get\n";
                        //t.sendResponseHeaders(400, response.length());
                        writeResponse(t, response, 400);
                    }
                    System.out.println("cart items: " + response);
                    String itemIDs = response;

                    //2)now that we have the list of items in the cart, do a batch get from catalog
                    connection = new URL(catalog.batchGet() + "?items=" + itemIDs).openConnection();
                    connection.setRequestProperty("Accept-Charset", "UTF-8");
                    br = new BufferedReader(new InputStreamReader(connection.getInputStream()));
                    response = br.readLine();
                    if (response == null) {
                        response = "Error reading response from catalog:batchget\n";
                        //t.sendResponseHeaders(400, response.length());
                        responseCode = 400;
                        writeResponse(t, response, responseCode);
                    }

                    String items = response;
                    System.out.println("items: " + items);

                    //3) create orderID, passing items in the cart
                    //System.out.println("ORDERMANAGEMENT: " + orderManagement.create());
                    String request = orderManagement.create();
                    URL url = new URL(request);
                    HttpURLConnection httpConnection = (HttpURLConnection) url.openConnection();
                    query = String.format("userID=%s&items=%s", userID, items);
                    _setPostParameters(httpConnection, query.length());


                    try (DataOutputStream wr = new DataOutputStream(httpConnection.getOutputStream())) {
                        wr.write(query.getBytes());
                    }

                    responseCode = httpConnection.getResponseCode();
                    System.out.println("Response Code : " + responseCode);

                    br = new BufferedReader(new InputStreamReader(httpConnection.getInputStream()));
                    response = br.readLine();
                    if (response == null) {
                        response = "Error reading response from orders:create\n";
                        //t.sendResponseHeaders(400, response.length());
                        responseCode = 400;
                        writeResponse(t, response, responseCode);
                    }
                    String orderID = response;
                    System.out.println("OrderID: " + orderID);


                    //4)delete the cart
                    request = cart.deleteCart();
                    url = new URL(request);
                    httpConnection = (HttpURLConnection) url.openConnection();
                    query = String.format("cartID=%s", cartID);
                    _setPostParameters(httpConnection, query.length());


                    try (DataOutputStream wr = new DataOutputStream(httpConnection.getOutputStream())) {
                        wr.write(query.getBytes());
                    }

                    responseCode = httpConnection.getResponseCode();
                    System.out.println("Response Code : " + responseCode);


                    if (responseCode != 200) {
                        response = "Error reading response from cart:delete\n";
                    } else {
                        response = orderID;
                    }

                }
            } else {
                response = "Only POST requests\n";
                //t.sendResponseHeaders(405, response.length());
                writeResponse(t, response, 405);
            }
            System.out.println("END, response code: " + responseCode);
            t.sendResponseHeaders(responseCode, response.length());
            os.write(response.getBytes());
            os.close();
        }

    }

    static class TestHandler implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            System.out.println("Checking Order status");
            String response = "Checking Order status\n";
            BufferedReader br;

            String urlParameters = "param1=a&param2=b&param3=c";
            byte[] postData = urlParameters.getBytes(charset);
            int postDataLength = postData.length;
            String request = "http://localhost:8008/orders/create";
            URL url = new URL(request);
            HttpURLConnection conn = (HttpURLConnection) url.openConnection();
            conn.setDoOutput(true);
            conn.setInstanceFollowRedirects(false);
            conn.setRequestMethod("POST");
            conn.setRequestProperty("Content-Type", "application/x-www-form-urlencoded");
            conn.setRequestProperty("charset", "utf-8");
            conn.setRequestProperty("Content-Length", Integer.toString(postDataLength));
            conn.setUseCaches(false);
            DataOutputStream wr = new DataOutputStream(conn.getOutputStream());
            wr.writeBytes(urlParameters);
            wr.flush();
            wr.close();

            System.out.println("here");
            int responseCode = conn.getResponseCode();
            System.out.println("Response Code : " + responseCode);

            br = new BufferedReader(new InputStreamReader(conn.getInputStream()));
            response = br.readLine();
            System.out.println("--here");
            br.close();
            System.out.println("response: " + response);


            t.sendResponseHeaders(200, response.length());
            OutputStream os = t.getResponseBody();
            os.write(response.getBytes());
            os.close();


            //given orderid, call /orders/summary
        }
    }


    static class OrderStatus implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            System.out.println("Checking Order status");
            String response = "This should not be printed\n";
            int returnCode = INTERNAL_ERR;

            if (t.getRequestMethod().equalsIgnoreCase("POST")) {
                String query = t.getRequestURI().getQuery();
                parseQuery(query, parameters);
                String orderID = parameters.get("orderID").toString();

                if (orderID == null) {
                    response = "Could not parse itemID\n";
                    returnCode = BAD_REQUEST;
                } else {
                    //we have parsed the itemID, now ask the catalog for the information
                    //for now let's just retrieve the title and price from catalog
                    URL url = new URL(orderManagement.summary());
                    HttpURLConnection httpConnection= (HttpURLConnection)url.openConnection();
                    httpConnection.setRequestProperty( "charset", "utf-8");

                    int responseCode = httpConnection.getResponseCode();
                    if(responseCode != SUCCESS) {
                        returnCode = responseCode;
                        response = "Error reading response from catalog:items\n";
                    }else {
                        BufferedReader rd = new BufferedReader( new InputStreamReader(httpConnection.getInputStream()));
                        response = rd.readLine();

                        if(response == null) {
                            response = "Error reading response from orders:summary\n";
                            returnCode = BAD_REQUEST;
                        }else{
                            returnCode = SUCCESS;
                        }

                    }
                }
            } else {
                response = "Only POST requests\n";
                returnCode = METHOD_NOT_ALLOWED;
            }

            t.sendResponseHeaders(returnCode, response.length());
            OutputStream os = t.getResponseBody();
            os.write(response.getBytes());
            os.close();
        }
    }//end OrderStatus

}//end RequestHandler.java
