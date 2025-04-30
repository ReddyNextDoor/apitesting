package com.example.personapi;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication(
    exclude = {
        org.springframework.boot.autoconfigure.data.jpa.JpaRepositoriesAutoConfiguration.class,
        org.springframework.boot.autoconfigure.data.mongo.MongoRepositoriesAutoConfiguration.class
    }
)
public class PersonApiServiceApplication {
    public static void main(String[] args) {
        SpringApplication.run(PersonApiServiceApplication.class, args);
    }
}
