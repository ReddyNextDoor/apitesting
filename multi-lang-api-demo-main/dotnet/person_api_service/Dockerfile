# Use .NET 8 SDK for build
FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src
COPY . .
RUN dotnet restore
RUN dotnet publish -c Release -o /app

# Runtime image
FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS runtime
WORKDIR /app
COPY --from=build /app .
RUN mkdir -p /app/db
EXPOSE 8000
ENTRYPOINT ["dotnet", "PersonApiService.dll"]
