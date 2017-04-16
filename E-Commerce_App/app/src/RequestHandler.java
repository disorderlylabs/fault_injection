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

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.util.ArrayList;
import java.util.List;


public class RequestHandler {

    static microservice.Catalog catalog = new microservice.Catalog();
    static microservice.Cart cart = new microservice.Cart();
    static microservice.OrderManagement orderManagement = new microservice.OrderManagement();
    static HttpClient client = HttpClientBuilder.create().build();


    static class BrowseHandler implements HttpHandler {
        public void handle(HttpExchange t) throws IOException {
            String response = "this should not be displayed";
            OutputStream os = t.getResponseBody();
            StringBuffer req = new StringBuffer();


            if(t.getRequestMethod().equalsIgnoreCase("GET")) {
                Headers headers = t.getRequestHeaders();

                String itemID = headers.getFirst("itemID");
                if(itemID == null) {
                    response = "Could not parse itemID\n";
                    t.sendResponseHeaders(400, response.length());
                }else{
                    //we have parsed the itemID, now ask the catalog for the information
                    //for now let's just retrieve the title and price from catalog
                    req.append(catalog.get());  //append the URL of the get request
                    req.append("?itemID=" + itemID);

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
            System.out.println("Checking out");
            String response = "Checking out\n";
            t.sendResponseHeaders(200, response.length());
            OutputStream os = t.getResponseBody();
            os.write(response.getBytes());
            os.close();
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
        }
    }




}
