using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;
using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
using System.Text.Json.Serialization; // Add this at the top

namespace PersonApiService.Models
{
    public class Person
    {
        [JsonPropertyName("id")]
        [BsonElement("id")]
        public string? Id { get; set; } // Optional for POST; set automatically if null

        [Key]
        [DatabaseGenerated(DatabaseGeneratedOption.Identity)]
        [JsonIgnore]
        [BsonIgnore]
        public int? SqliteId { get; set; } // Used only for SQLite, not exposed in API

        [Required]
        [JsonPropertyName("first_name")]
        [BsonElement("first_name")]
        public string FirstName { get; set; } = string.Empty;
        [Required]
        [JsonPropertyName("last_name")]
        [BsonElement("last_name")]
        public string LastName { get; set; } = string.Empty;
        [Required]
        [JsonPropertyName("age")]
        [BsonElement("age")]
        public int Age { get; set; }
        [Required]
        [JsonPropertyName("address")]
        [BsonElement("address")]
        public Address Address { get; set; } = new Address();
    }
}

