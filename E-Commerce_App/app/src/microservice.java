
public class microservice {

    //code for microservice
    static class service {
        protected int portNum;
        protected String addr;
        protected String baseURL;


        public service(String a, int p) {
            addr = a;
            portNum = p;
            baseURL = "http://" + addr + ":" + portNum + "";
        }

    }


    static class Catalog extends service{
        String get;
        String batchGet;
        String add;
        String items;
        String update;
        String delete;

        public Catalog() {
            super("localhost", 8008);  //catalog listening on port 8000

            get = baseURL + "/catalog/get";
            add = baseURL + "/catalog/add";
            items = baseURL + "/catalog/items";
            update = baseURL + "/catalog/update";
            delete = baseURL + "/catalog/delete";
            batchGet = baseURL + "/catalog/batchGet";
        }

        public String get() { return get; }
        public String add() { return add; }
        public String items() { return items; }
        public String update() { return update; }
        public String delete() { return delete; }
        public String batchGet() { return batchGet; }
    }





    static class Cart extends service {
        //endpoints for the different functions
        String create;
        String addItem;
        String deleteItem;
        String deleteCart;
        String items;

        public Cart() {
            super("localhost", 8008);

            create = baseURL + "/cart/create";
            addItem = baseURL + "/cart/addItem";
            deleteItem = baseURL + "/cart/deleteItem";
            deleteCart = baseURL + "/cart/deleteCart";
            items = baseURL + "/cart/items";
        }

        public String create() {return create; }
        public String addItem() { return addItem; }
        public String deleteItem() { return deleteItem; }
        public String deleteCart() { return deleteCart; }
        public String items() { return items; }
    }



    static class OrderManagement extends service {
        String create;
        String shipping;
        String payment;
        String summary;

        public OrderManagement() {
            super("localhost", 8008);

            create = baseURL + "/orders/create";
            shipping = baseURL + "/orders/shipping";
            payment = baseURL + "/orders/payment";
            summary = baseURL + "/orders/summary";
        }

        public String create() { return create; }
        public String shipping() { return shipping; }
        public String payment() { return payment; }
        public String summary() { return summary; }
    }

}
