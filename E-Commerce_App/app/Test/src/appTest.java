import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.URL;
import java.net.URLConnection;
import java.net.URLEncoder;


//cartID = 2020


public class appTest {
    static String charset = "UTF-8";
    static String cartID = "2020";
    static String itemID = "111";

    public static void testCheckout() throws Exception{
        String query = String.format("cartID=%s", URLEncoder.encode(cartID, charset));
        URLConnection connection = new URL("http://localhost:1339/app/checkout").openConnection();
        connection.setDoOutput(true); // Triggers POST.
        connection.setRequestProperty("Accept-Charset", charset);
        connection.setRequestProperty("Content-Type", "application/x-www-form-urlencoded;charset=" + charset);
        try (OutputStream output = connection.getOutputStream()) {
            output.write(query.getBytes(charset));
        }

        BufferedReader rd = new BufferedReader(
                new InputStreamReader(connection.getInputStream()));

        String response = rd.readLine();
        if(response == null) {
            response = "Error reading response from catalog:get\n";
        }
        System.out.println("Response: " + response);
    }



    public static void testGet() throws Exception {
        String query = String.format("itemID=%s", URLEncoder.encode(itemID, charset));
        URLConnection connection = new URL("http://localhost:1339/app/browse" + "?itemID=111").openConnection();
        connection.setRequestProperty("Accept-Charset", charset);

        BufferedReader rd = new BufferedReader(
                new InputStreamReader(connection.getInputStream()));

        String response = rd.readLine();
        if(response == null) {
            response = "Error reading response from catalog:get\n";
        }
        System.out.println("Response: " + response);


    }


    public static void main(String args[]) throws Exception{
        testCheckout();
    }



}
