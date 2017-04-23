import com.sun.net.httpserver.Headers;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import org.apache.http.HttpResponse;
import org.apache.http.NameValuePair;
import org.apache.http.client.HttpClient;
import org.apache.http.client.entity.UrlEncodedFormEntity;
import org.apache.http.client.methods.HttpDelete;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.message.BasicNameValuePair;

import java.io.*;
import java.net.*;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;


public class RequestHandler {

    static microservice.Catalog catalog = new microservice.Catalog();
    static microservice.Cart cart = new microservice.Cart();
    static microservice.OrderManagement orderManagement = new microservice.OrderManagement();
    static HttpClient client = HttpClientBuilder.create().build();
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


    static class BrowseHandler implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            String response = "this should not be displayed";
            OutputStream os = t.getResponseBody();
            StringBuffer req = new StringBuffer();
            parameters.clear();


            if(t.getRequestMethod().equalsIgnoreCase("GET")) {
                //get the URI and parse the parameters
                String query = t.getRequestURI().getQuery();
                parseQuery(query, parameters);
                String itemID = parameters.get("itemID").toString();


                if(itemID == null) {
                    response = "Could not parse itemID\n";
                    t.sendResponseHeaders(400, response.length());
                }else{
                    //we have parsed the itemID, now ask the catalog for the information
                    //for now let's just retrieve the title and price from catalog
                    req.append(catalog.items());  //append the URL of the get request

                    HttpGet request = new HttpGet(req.toString());
                    HttpResponse res = client.execute(request);
                    System.out.println("Response Code : "
                            + res.getStatusLine().getStatusCode());

                    BufferedReader rd = new BufferedReader(
                            new InputStreamReader(res.getEntity().getContent()));

                    response = rd.readLine();
                    if(response == null) {
                        response = "Error reading response from catalog:get\n";
                        t.sendResponseHeaders(400, response.length());
                    }

                    //response should be in the form "title:price"
                    //just return it for now.
                }

            }else{
                response = "Only GET requests\n";
                t.sendResponseHeaders(405, response.length());
            }


            os.write(response.getBytes());
            os.close();
        }
    }

    static class CartCreate implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            String response = "this should not be displayed";
            OutputStream os = t.getResponseBody();
            StringBuffer req = new StringBuffer();


            if(t.getRequestMethod().equalsIgnoreCase("POST")) {
                Headers headers = t.getRequestHeaders();

                String userID = headers.getFirst("userID");
                if(userID == null) {
                    response = "Could not parse userID\n";
                    t.sendResponseHeaders(400, response.length());
                }else{
                    //create the post request
                    String url = cart.create();
                    HttpPost postReq = new HttpPost(url);

                    //add parameters for post request
                    List<NameValuePair> params = new ArrayList<NameValuePair>(1);
                    params.add(new BasicNameValuePair("userID", userID));
                    postReq.setEntity(new UrlEncodedFormEntity(params, "UTF-8"));

                    HttpResponse res = client.execute(postReq);
                    System.out.println("Response Code : "
                            + res.getStatusLine().getStatusCode());

                    BufferedReader rd = new BufferedReader(
                            new InputStreamReader(res.getEntity().getContent()));

                    response = rd.readLine();
                    if(response == null) {
                        response = "Error reading response from catalog:get\n";
                        t.sendResponseHeaders(400, response.length());
                    }
                }
            }else{
                response = "Only POST requests\n";
                t.sendResponseHeaders(405, response.length());
            }

            os.write(response.getBytes());
            os.close();
        }
    }



    static class CartAdd implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            String response = "this should not be displayed";
            OutputStream os = t.getResponseBody();
            StringBuffer req = new StringBuffer();


            if(t.getRequestMethod().equalsIgnoreCase("POST")) {
                Headers headers = t.getRequestHeaders();

                String itemID = headers.getFirst("itemID");
                if(itemID == null) {
                    response = "Could not parse itemID\n";
                    t.sendResponseHeaders(400, response.length());
                    os.write(response.getBytes());
                    os.close();
                    return;
                }

                String cartID = headers.getFirst("cartID");
                if(cartID == null) {
                    response = "Could not parse cartID\n";
                    t.sendResponseHeaders(400, response.length());
                }else{
                    String url = cart.add();
                    HttpPost postReq = new HttpPost(url);


                    List<NameValuePair> params = new ArrayList<NameValuePair>(2);
                    params.add(new BasicNameValuePair("itemID", itemID));
                    params.add(new BasicNameValuePair("cartID", cartID));
                    postReq.setEntity(new UrlEncodedFormEntity(params, "UTF-8"));

                    HttpResponse res = client.execute(postReq);
                    System.out.println("Response Code : "
                            + res.getStatusLine().getStatusCode());

                    BufferedReader rd = new BufferedReader(
                            new InputStreamReader(res.getEntity().getContent()));

                    response = rd.readLine();
                    if(response == null) {
                        response = "Error reading response from catalog:get\n";
                        t.sendResponseHeaders(400, response.length());
                    }

                }
            }else{
                response = "Only POST requests\n";
                t.sendResponseHeaders(405, response.length());
            }

            os.write(response.getBytes());
            os.close();

        }
    }



    static class CartDelete implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            String response = "this should not be displayed";
            OutputStream os = t.getResponseBody();
            StringBuffer req = new StringBuffer();


            if(t.getRequestMethod().equalsIgnoreCase("DELETE")) {
                Headers headers = t.getRequestHeaders();

                String itemID = headers.getFirst("itemID");
                if(itemID == null) {
                    response = "Could not parse itemID\n";
                    t.sendResponseHeaders(400, response.length());
                    os.write(response.getBytes());
                    os.close();
                    return;
                }

                String cartID = headers.getFirst("cartID");
                if(cartID == null) {
                    response = "Could not parse cartID\n";
                    t.sendResponseHeaders(400, response.length());
                }else{
                    String url = cart.delete();
                    HttpDelete deleteReq = new HttpDelete(url);

                    deleteReq.addHeader("itemID", itemID);
                    deleteReq.addHeader("cartID", cartID);

                    HttpResponse res = client.execute(deleteReq);
                    System.out.println("Response Code : "
                            + res.getStatusLine().getStatusCode());

                    BufferedReader rd = new BufferedReader(
                            new InputStreamReader(res.getEntity().getContent()));

                    response = rd.readLine();
                    if(response == null) {
                        response = "Error reading response from catalog:get\n";
                        t.sendResponseHeaders(400, response.length());
                    }

                }
            }else{
                response = "Only DELETE requests\n";
                t.sendResponseHeaders(405, response.length());
            }

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


            if(t.getRequestMethod().equalsIgnoreCase("POST")) {
                InputStreamReader isr = new InputStreamReader(t.getRequestBody(), "utf-8");
                br = new BufferedReader(isr);
                String query = br.readLine();
                parseQuery(query, parameters);


                String cartID = parameters.get("cartID").toString();
                String userID = parameters.get("userID").toString();

                if(cartID == null) {
                    response = "Could not parse cartID\n";
                    //t.sendResponseHeaders(400, response.length());
                    responseCode = 400;
                    writeResponse(t, response, responseCode);
                }else if(userID == null) {
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
                    if(response == null) {
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
                    if(response == null) {
                        response = "Error reading response from catalog:batchget\n";
                        //t.sendResponseHeaders(400, response.length());
                        responseCode = 400;
                        writeResponse(t, response, responseCode);
                    }

                    String items = response;
                    System.out.println("items: " + items);

                    //3) create orderID, passing items in the cart
                    HttpURLConnection httpConnection = (HttpURLConnection) new URL(orderManagement.create()).openConnection();
                    httpConnection.setRequestMethod("POST");

                    query = String.format("userID=%s&items=%s",
                            URLEncoder.encode(userID, charset),
                            URLEncoder.encode(items, charset));
//                    connection = new URL(orderManagement.create()).openConnection();
                    httpConnection.setDoOutput(true);
                    httpConnection.setRequestProperty("Accept-Charset", "UTF-8");
                    httpConnection.setRequestProperty(
                            "Content-Type", "application/x-www-form-urlencoded" );
                    try (OutputStream output = httpConnection.getOutputStream()) {
                        output.write(query.getBytes(charset));
                    }

                    br = new BufferedReader(new InputStreamReader(httpConnection.getInputStream()));
                    response = br.readLine();
                    if(response == null) {
                        response = "Error reading response from orders:create\n";
                        //t.sendResponseHeaders(400, response.length());
                        responseCode = 400;
                        writeResponse(t, response, responseCode);
                    }


                    //4)delete the cart
                    query = String.format("cartID=%s", URLEncoder.encode(cartID, charset));
                    connection = new URL(cart.delete()).openConnection();
                    connection.setRequestProperty("Content-Type", "application/x-www-form-urlencoded;charset=" + charset);
                    try (OutputStream output = connection.getOutputStream()) {
                        output.write(query.getBytes(charset));
                    }
                    br = new BufferedReader(new InputStreamReader(connection.getInputStream()));
                    response = br.readLine();
                    if(response == null) {
                        response = "Error reading response from cart:delete\n";
                        //t.sendResponseHeaders(400, response.length());
                        responseCode = 400;
                        writeResponse(t, response, responseCode);

                    }else{
                        System.out.println("response: " + response);
                        writeResponse(t, response, 200);
                    }
                }
            }else{
                response = "Only POST requests\n";
                //t.sendResponseHeaders(405, response.length());
                writeResponse(t, response, 405);
            }
            t.sendResponseHeaders(responseCode, response.length());
        }

    }




    static class OrderStatus implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            System.out.println("Checking Order status");
            String response = "Checking Order status\n";
            t.sendResponseHeaders(200, response.length());
            OutputStream os = t.getResponseBody();
            os.write(response.getBytes());
            os.close();

            //given orderid, call /orders/summary
        }
    }




}
