package com.example.personapi.repository;

import com.example.personapi.repository.jpa.JpaPersonRepositoryCustomImpl;
import com.example.personapi.repository.mongo.MongoPersonRepositoryCustomImpl;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.ApplicationContext;
import org.springframework.core.env.Environment;
import org.springframework.stereotype.Component;

@Component
public class PersonRepositoryFactory {
    private final ApplicationContext context;
    private final Environment env;

    @Autowired
    public PersonRepositoryFactory(ApplicationContext context, Environment env) {
        this.context = context;
        this.env = env;
    }

    public PersonRepository getPersonRepository() {
        String activeProfile = env.getProperty("spring.profiles.active", "mongo");
        if ("sqlite".equals(activeProfile)) {
            return context.getBean(JpaPersonRepositoryCustomImpl.class);
        } else {
            return context.getBean(MongoPersonRepositoryCustomImpl.class);
        }
    }
}
