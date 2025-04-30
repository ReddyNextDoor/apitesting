package com.example.personapi.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Profile;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;

@Configuration
@Profile("sqlite")
@EnableJpaRepositories(basePackages = "com.example.personapi.repository.jpa")
public class JpaConfig {}
