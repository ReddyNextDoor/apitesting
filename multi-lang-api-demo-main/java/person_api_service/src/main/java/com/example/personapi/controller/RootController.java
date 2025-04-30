package com.example.personapi.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.stereotype.Controller;

@Controller
public class RootController {
    @GetMapping("/")
    public String redirectToSwaggerUi() {
        return "redirect:/swagger-ui/index.html";
    }
}
