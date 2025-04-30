using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
using MongoDB.Bson.Serialization.Attributes;
using System.Text.Json.Serialization;

namespace PersonApiService.Models
{
    public class Address
    {
        [Required]
        [JsonPropertyName("address_line1")]
        [BsonElement("address_line1")]
        public string AddressLine1 { get; set; } = string.Empty;
        
        [JsonPropertyName("address_line2")]
        [BsonElement("address_line2")]
        public string? AddressLine2 { get; set; }
        
        [Required]
        [JsonPropertyName("city")]
        [BsonElement("city")]
        public string City { get; set; } = string.Empty;
        
        [Required]
        [JsonPropertyName("state")]
        [BsonElement("state")]
        public string State { get; set; } = string.Empty;
        
        [Required]
        [JsonPropertyName("zip")]
        [BsonElement("zip")]
        public string Zip { get; set; } = string.Empty;
    }
}

