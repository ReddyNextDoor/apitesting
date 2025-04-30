package com.example.personapi.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.beans.factory.annotation.Value;
import java.util.Map;
import java.time.ZonedDateTime;
import java.util.LinkedHashMap;

@RestController
public class HealthController {
    @Value("${person.repo.backend}")
    private String repoBackend;

    @GetMapping("/health")
    public Map<String, Object> health() {
        Map<String, Object> status = new LinkedHashMap<>();
        status.put("service", "person-api-service");
        status.put("status", "ok");
        status.put("mode", repoBackend);
        status.put("time", ZonedDateTime.now().toString());
        status.put("database", "ok"); // Replace with real DB check if needed
        return status;
    }
}
