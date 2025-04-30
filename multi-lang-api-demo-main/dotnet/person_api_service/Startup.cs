using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.EntityFrameworkCore;
using MongoDB.Driver;
using PersonApiService.Repositories;
using PersonApiService.Models;

namespace PersonApiService
{
    public class Startup
    {
        public IConfiguration Configuration { get; }
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public void ConfigureServices(IServiceCollection services)
        {
            var backend = Environment.GetEnvironmentVariable("PERSON_REPO_BACKEND")?.ToLower() ?? "sqlite";
            if (backend == "mongo")
            {
                var mongoUri = Environment.GetEnvironmentVariable("MONGO_URI");
                if (string.IsNullOrEmpty(mongoUri))
                    throw new InvalidOperationException("MONGO_URI environment variable must be set for MongoDB backend.");
                var mongoClient = new MongoClient(mongoUri);
                var database = mongoClient.GetDatabase("person_db");
                services.AddSingleton(database);
                services.AddScoped<IPersonRepository, MongoPersonRepository>();
            }
            else
            {
                var dbPath = Environment.GetEnvironmentVariable("DOTNET_RUNNING_IN_CONTAINER") == "true"
                    ? "/app/db/persons.db"
                    : "db/persons.db";
                services.AddDbContext<PersonContext>(options =>
                    options.UseSqlite($"Data Source={dbPath}"));
                services.AddScoped<IPersonRepository, SqlitePersonRepository>();
            }
            services.AddControllers();
            services.AddEndpointsApiExplorer();
            services.AddSwaggerGen(c =>
            {
                c.SwaggerDoc("v1", new Microsoft.OpenApi.Models.OpenApiInfo
                {
                    Title = "Person API Service",
                    Version = "v1",
                    Description = "ASP.NET Core REST API for Person management with SQLite and MongoDB"
                });
            });

            // Ensure SQLite schema is created
            if (backend != "mongo")
            {
                using (var scope = services.BuildServiceProvider().CreateScope())
                {
                    var db = scope.ServiceProvider.GetRequiredService<PersonContext>();
                    db.Database.EnsureCreated();
                }
            }
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            app.UseSwagger();
            app.UseSwaggerUI(c =>
            {
                c.SwaggerEndpoint("/swagger/v1/swagger.json", "Person API V1");
                c.RoutePrefix = "docs";
            });

            // Redirect root / to /docs
            app.Use(async (context, next) =>
            {
                if (context.Request.Path == "/")
                {
                    context.Response.Redirect("/docs");
                    return;
                }
                await next();
            });
            app.UseRouting();
            app.UseAuthorization();
            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
            });
        }
    }
}
