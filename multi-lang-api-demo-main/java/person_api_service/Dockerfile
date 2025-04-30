# ---- Build Stage ----
FROM maven:3.9.9-amazoncorretto-21 AS build
WORKDIR /app
COPY . .
RUN mvn clean package -DskipTests

# ---- Run Stage ----
FROM amazoncorretto:21
WORKDIR /app
RUN mkdir -p /app/db
COPY --from=build /app/target/person-api-service-1.0.0.jar ./person-api-service.jar
EXPOSE 8000
ENTRYPOINT ["java", "-jar", "person-api-service.jar"]
