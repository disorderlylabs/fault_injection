
public class microservice {

    //code for microservice
    static class service {
        protected int portNum;
        protected String addr;
        protected String baseURL;


        public service(String a, int p) {
            addr = a;
            portNum = p;
            baseURL = "http://" + addr + ":" + portNum + "/";
        }

    }


    static class Catalog extends service{
        String get;
        String add;
        String update;
        String delete;

        public Catalog() {
            super("localhost", 8000);  //catalog listening on port 8000

            get = baseURL + "/catalog/get";
            add = baseURL + "/catalog/add";
            update = baseURL + "/catalog/update";
            delete = baseURL + "/catalog/delete";
        }

        public String get() { return get; }
        public String add() { return add; }
        public String update() { return update; }
        public String delete() { return delete; }
    }





    static class Cart extends service {
        //endpoints for the different functions
        String create;
        String add;
        String delete;
        String items;

        public Cart() {
            super("localhost", 1339);

            create = baseURL + "/cart/create";
            add = baseURL + "/cart/add";
            delete = baseURL + "/cart/delete";
            items = baseURL + "/cart/items";
        }

        public String create() {return create; }
        public String add() { return add; }
        public String delete() { return delete; }
        public String items() { return items; }
    }



    static class OrderManagement extends service {
        String create;
        String shipping;
        String payment;
        String summary;

        public OrderManagement() {
            super("localhost", 8002);

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
